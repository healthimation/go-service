package eventbus

type Message struct {
	ID      string            `json:"id"`
	Key     string            `json:"key"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

type MessageHandler func(message *Message)
