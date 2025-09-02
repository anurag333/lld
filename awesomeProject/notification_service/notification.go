package notification_service

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"
)

// --- Domain types ---

type Channel string

const (
	ChannelEmail Channel = "EMAIL"
	ChannelSMS   Channel = "SMS"
)

type NotificationStatus string

const (
	StatusDraft     NotificationStatus = "DRAFT"
	StatusQueued    NotificationStatus = "QUEUED"
	StatusSending   NotificationStatus = "SENDING"
	StatusSent      NotificationStatus = "SENT"
	StatusPartially NotificationStatus = "PARTIALLY_SENT"
	StatusFailed    NotificationStatus = "FAILED"
	StatusCancelled NotificationStatus = "CANCELLED"
)

type AttemptStatus string

const (
	AttemptPending    AttemptStatus = "PENDING"
	AttemptDelivering AttemptStatus = "DELIVERING"
	AttemptDelivered  AttemptStatus = "DELIVERED"
	AttemptFailed     AttemptStatus = "FAILED"
	AttemptCancelled  AttemptStatus = "CANCELLED"
)

// --- App ---

type App struct {
	store    *InMemoryStore
	outbox   chan OutboxEvent // dont use
	shutdown chan struct{}
}

func NewApp() *App {
	return &App{
		store:    NewInMemoryStore(),
		outbox:   make(chan OutboxEvent, 1000),
		shutdown: make(chan struct{}),
	}
}

// create notification
func (a *App) CreateNotification(noti *Notification) (id string, status string, err error) {
	if noti.ID == "" {
		noti.ID = strconv.FormatInt(time.Now().UnixNano(), 36)
	}
	if len(noti.Channels) == 0 {
		noti.Channels = []Channel{ChannelEmail}
	}
	noti.Status = StatusQueued
	a.store.CreateNotification(noti)
	for _, rc := range noti.Recipients {
		for _, ch := range noti.Channels {
			att := &DeliveryAttempt{
				NotificationID: noti.ID,
				RecipientID:    rc.ID,
				Channel:        ch,
				Status:         AttemptPending,
				AttemptNo:      0,
			}
			a.store.AddAttempt(att)
		}
	}
	select {
	case a.outbox <- OutboxEvent{NotificationID: noti.ID, CreatedAt: time.Now()}:
	default:
	}
	return noti.ID, string(noti.Status), nil
}

// list notifications
func (a *App) ListNotifications() []*Notification {
	return a.store.ListNotifications()
}

// patch notification status
//func PatchNotificationStatus(store *InMemoryStore, id string, status string) error {
//	return store.UpdateNotificationStatus(id, NotificationStatus(status))
//}

// send notification
func (a *App) SendNotification(id string) error {
	n, ok := a.store.GetNotification(id)
	if !ok {
		return fmt.Errorf("not found")
	}
	if n.Status == StatusDraft {
		n.Status = StatusQueued
	}
	for _, rc := range n.Recipients {
		for _, ch := range n.Channels {
			att := &DeliveryAttempt{NotificationID: n.ID, RecipientID: rc.ID, Channel: ch, Status: AttemptPending}
			a.store.AddAttempt(att)
		}
	}
	select {
	case a.outbox <- OutboxEvent{NotificationID: n.ID, CreatedAt: time.Now()}:
	default:
	}
	return nil
}



// worker loop: poll pending attempts and send
func (a *App) SenderWorker(ctx context.Context, id int) {
	log.Printf("senderWorker %d started", id)
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			log.Printf("senderWorker %d stopping", id)
			return
		case <-ticker.C:
			atts := a.store.NextPendingAttempts(10)
			for _, att := range atts {
				// quick state flip
				att.Status = AttemptDelivering
				a.store.UpdateAttempt(att)
				// load notification + recipient
				n, ok := a.store.GetNotification(att.NotificationID)
				if !ok {
					att.Status = AttemptFailed
					a.store.UpdateAttempt(att)
					continue
				}
				r := findRecipient(n, att.RecipientID)
				subj, body := render(n, r, att.Channel)
				// choose address
				//addr := ""
				//if r.Addresses != nil {
				//	if v, ok := r.Addresses[string(att.Channel)]; ok {
				//		addr = v
				//	}
				//}
				//if addr == "" {
				//	att.Status = AttemptFailed
				//	a.store.UpdateAttempt(att)
				//	continue
				//}
				ctx2, cancel := context.WithTimeout(ctx, 5*time.Second)
				_, err := sendToProvider(ctx2, att.Channel, subj, body)
				cancel()
				if err != nil {
					att.AttemptNo++
					if att.AttemptNo >= 3 {
						att.Status = AttemptFailed
					} else {
						att.Status = AttemptPending
						att.NextAttemptAt = time.Now().Add(time.Duration(1<<att.AttemptNo) * time.Second)
					}
				} else {
					att.Status = AttemptDelivered
				}
				a.store.UpdateAttempt(att)
			}
		}
	}
}
