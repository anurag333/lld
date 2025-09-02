package notification_service_outbox

import (
	"errors"
	"sync"
	"time"
)

// --- In-memory repo ---

type InMemoryStore struct {
	mu            sync.RWMutex
	notifications map[string]*Notification
	attempts      map[int64]*DeliveryAttempt
	nextAttemptID int64 // auto-incrementing ID for attempts
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		notifications: make(map[string]*Notification),
		attempts:      make(map[int64]*DeliveryAttempt),
	}
}

func (s *InMemoryStore) CreateNotification(n *Notification) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.notifications[n.ID]; exists {
		return errors.New("id exists")
	}
	n.CreatedAt = time.Now().UTC()
	s.notifications[n.ID] = n
	return nil
}

func (s *InMemoryStore) ListNotifications() []*Notification {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*Notification, 0, len(s.notifications))
	for _, n := range s.notifications {
		list = append(list, n)
	}
	return list
}

func (s *InMemoryStore) GetNotification(id string) (*Notification, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	n, ok := s.notifications[id]
	return n, ok
}

func (s *InMemoryStore) UpdateNotificationStatus(id string, to NotificationStatus) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	n, ok := s.notifications[id]
	if !ok {
		return errors.New("not found")
	}
	// simple allowed transitions for demo
	if n.Status == StatusSent || n.Status == StatusCancelled {
		return errors.New("cannot transition")
	}
	n.Status = to
	return nil
}

func (s *InMemoryStore) AddAttempt(a *DeliveryAttempt) int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.nextAttemptID++
	a.ID = s.nextAttemptID
	a.UpdatedAt = time.Now().UTC()
	s.attempts[a.ID] = a
	return a.ID
}

func (s *InMemoryStore) NextPendingAttempts(limit int) []*DeliveryAttempt {
	s.mu.Lock()
	defer s.mu.Unlock()
	r := make([]*DeliveryAttempt, 0)
	now := time.Now().UTC()
	for _, a := range s.attempts {
		if a.Status == AttemptPending && (a.NextAttemptAt.IsZero() || a.NextAttemptAt.Before(now) || a.NextAttemptAt.Equal(now)) {
			r = append(r, a)
			if len(r) >= limit {
				break
			}
		}
	}
	return r
}

func (s *InMemoryStore) UpdateAttempt(a *DeliveryAttempt) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if existing, ok := s.attempts[a.ID]; ok {
		existing.Status = a.Status
		existing.AttemptNo = a.AttemptNo
		existing.NextAttemptAt = a.NextAttemptAt
		existing.UpdatedAt = time.Now().UTC()
	}
}
