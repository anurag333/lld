package notification_service_outbox

import "time"

// --- Outbox / Worker wiring ---

type OutboxEvent struct {
	NotificationID string
	CreatedAt      time.Time
}

type Notification struct {
	ID         string                 `json:"id"`
	Title      string                 `json:"title,omitempty"`
	Body       string                 `json:"body,omitempty"`
	Channels   []Channel              `json:"channels"`
	Data       map[string]interface{} `json:"data,omitempty"`
	Recipients []Recipient            `json:"recipients,omitempty"`
	Status     NotificationStatus     `json:"status"`
	CreatedAt  time.Time              `json:"createdAt"`
}

type Recipient struct {
	ID string `json:"id"`
	//Addresses map[string]string `json:"addresses,omitempty"`
}

type DeliveryAttempt struct {
	ID             int64         `json:"id"`
	NotificationID string        `json:"notificationId"`
	RecipientID    string        `json:"recipientId"`
	Channel        Channel       `json:"channel"`
	Status         AttemptStatus `json:"status"`
	AttemptNo      int           `json:"attemptNo"`
	NextAttemptAt  time.Time     `json:"nextAttemptAt,omitempty"`
	UpdatedAt      time.Time     `json:"updatedAt"`
}
