package zulip

import "time"

type BaseResponse struct {
	Result string `json:"result"`
	Msg    string `json:"msg"`
}

type RegisterResponse struct {
	BaseResponse
	QueueID      string `json:"queue_id"`
	MaxMessageID int64  `json:"max_message_id"`
	LastEventID  int64  `json:"last_event_id"`
}

type MessageResponse struct {
	Content          string      `json:"content"`
	Subject          string      `json:"subject"`
	ContentType      string      `json:"content_type"`
	Client           string      `json:"client"`
	GravatarHash     string      `json:"gravatar_hash"`
	AvatarUrl        string      `json:"avatar_url"`
	ID               int64       `json:"id"`
	SenderID         int         `json:"sender_id"`
	RecipientID      int         `json:"recipient_id"`
	SenderEmail      string      `json:"sender_email"`
	SenderDomain     string      `json:"sender_domain"`
	SenderFullName   string      `json:"sender_full_name"`
	SenderShortName  string      `json:"sender_short_name"`
	DisplayRecipient interface{} `json:"display_recipient"`
	SubjectLinks     []string    `json:"subject_links"`
	Reactions        []string    `json:"reactions"`
	Mentioned        bool        `json:"is_mentioned"`
	Type             string      `json:"type"`
	Timestamp        int64       `json:"timestamp"`
}

type EventResponse struct {
	Flags   []string         `json:"flags"`
	Message *MessageResponse `json:"message"`
	Type    string           `json:"type"`
	ID      int64            `json:"id"`
}

type EventsResponse struct {
	BaseResponse
	QueueID string           `json:"queue_id"`
	Events  []*EventResponse `json:"events"`
}

func (m *MessageResponse) IsPrivate() bool {
	return m.Type == "private"
}

func (m *MessageResponse) GetStreamName() string {
	if m.IsPrivate() {
		return ""
	}

	return m.DisplayRecipient.(string)
}

func (m *MessageResponse) GetTopicName() string {
	if m.IsPrivate() {
		return ""
	}

	return m.Subject
}

func (m *MessageResponse) GetTime() time.Time {
	return time.Unix(m.Timestamp, 0)
}

func (m *MessageResponse) GetRecipients() []string {
	var res []string

	if m.IsPrivate() {
		recipients := m.DisplayRecipient.([]map[string]interface{})
		for _, recipient := range recipients {
			email := recipient["email"].(string)
			res = append(res, email)
		}
	}

	return res
}

func (e *EventResponse) IsMentioned() bool {
	for _, flag := range e.Flags {
		if flag == "mentioned" {
			return true
		}
	}
	return false
}
