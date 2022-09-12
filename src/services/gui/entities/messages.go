package entities

import (
	"fmt"
	"strings"
)

type Message struct {
	Type    MessageType
	Content string
}

// Flash types based on https://getbootstrap.com/docs/5.1/components/alerts/
type MessageType string

const (
	MessageTypeError   = "danger"
	MessageTypeSuccess = "success"
)

func (m Message) String() string {
	return fmt.Sprintf("%s: %s", m.Type, m.Content)
}

func ParseMessage(input string) Message {
	parts := strings.SplitN(input, ": ", 2)
	if len(parts) < 2 {
		return Message{Content: input}
	}
	return Message{
		Type:    MessageType(parts[0]),
		Content: parts[1],
	}
}
