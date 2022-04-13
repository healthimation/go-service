package eventbus

type Message struct {
	Key     string            `json:"key"`
	Headers map[string]string `json:"headers"`
	Body    []byte            `json:"body"`
}

type MessageHandler func(message *Message)
