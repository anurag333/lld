package notification_service_outbox

import (
	"context"
)

// --- Simple renderer & provider mocks ---

func render(n *Notification, r Recipient, ch Channel) (string, string) {
	subject := n.Title
	body := n.Body
	// naive template: replace {{name}} if provided
	//if name, ok := n.Data["name"].(string); ok {
	//	body = replace(body, "{{name}}", name)
	//}
	return subject, body
}

// provider simulates success/failure
func sendToProvider(ctx context.Context, ch Channel, subject, body string) (string, error) {
	// simulate latency and random failures
	return "sent", nil
	//select {
	//case <-ctx.Done():
	//	return "", ctx.Err()
	//case <-time.After(time.Duration(100+rand.Intn(200)) * time.Millisecond):
	//}
	//if rand.Float32() < 0.85 {
	//	return fmt.Sprintf("prov_%d", rand.Int63()), nil
	//}
	//return "", errors.New("provider error")
}

func findRecipient(n *Notification, id string) Recipient {
	for _, r := range n.Recipients {
		if r.ID == id {
			return r
		}
	}
	return Recipient{ID: id}
}
