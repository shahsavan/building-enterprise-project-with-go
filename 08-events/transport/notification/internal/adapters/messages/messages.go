package messages

// NotificationIssued is a generated struct.
type NotificationIssued struct {
	RecipientID string `avro:"recipientId"`
	Channel     string `avro:"channel"`
	Message     string `avro:"message"`
	EventType   string `avro:"eventType"`
	Timestamp   string `avro:"timestamp"`
}
