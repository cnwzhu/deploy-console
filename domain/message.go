package domain

//socket消息
type Message struct {
	Message string `json:"message"`
}

func (it *Message) String() string {
	return it.Message
}

func NewMessage(msg string) *Message {
	return &Message{Message: msg}
}
