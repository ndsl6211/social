package entity

import (
	"time"

	"github.com/google/uuid"
	"mashu.example/internal/entity"
)

type Message struct {
	ID        uuid.UUID
	Owner     *entity.User
	Content   string
	Timestamp time.Time
}

func NewMessageWithTime(
	id uuid.UUID,
	owner *entity.User,
	content string,
	time time.Time,
) *Message {
	return &Message{id, owner, content, time}
}

func NewMessage(
	id uuid.UUID,
	owner *entity.User,
	content string,
) *Message {
	return &Message{id, owner, content, time.Now()}
}

type DirectMessage struct {
	ID        uuid.UUID
	Creator   *entity.User
	Receiver  *entity.User
	Messages  []*Message
	CreatedAt time.Time
}

func NewDirectMessage(
	id uuid.UUID,
	creator *entity.User,
	receiver *entity.User,
) *DirectMessage {
	return &DirectMessage{id, creator, receiver, []*Message{}, time.Now()}
}
