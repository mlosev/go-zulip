package zulip

import (
	"encoding/json"
	"net/url"
)

type messageType string

const (
	privateMessage messageType = "private"
	streamMessage  messageType = "stream"
)

type message struct {
	Type                  messageType
	Content               string
	StreamName, TopicName string
	Usernames             []string
}

func NewPrivateMessage(recipients []string) *message {
	return &message{Type: privateMessage, Usernames: recipients}
}

func NewStreamMessage(stream, topic string) *message {
	return &message{Type: streamMessage, StreamName: stream, TopicName: topic}
}

func Reply(msg *MessageResponse) *message {
	if msg.IsPrivate() {
		return NewPrivateMessage(msg.GetRecipients())
	} else {
		return NewStreamMessage(msg.GetStreamName(), msg.GetTopicName())
	}
}

func (z *Zulip) SendMessage(msg *message) {
	v := url.Values{}
	v.Set("type", string(msg.Type))
	v.Set("content", msg.Content)
	switch msg.Type {
	case privateMessage:
		recipients, _ := json.Marshal(msg.Usernames)
		v.Set("to", string(recipients))
	case streamMessage:
		v.Set("to", msg.StreamName)
		v.Set("subject", msg.TopicName)
	}

	z.api("api/v1/messages", "POST", v)
}
