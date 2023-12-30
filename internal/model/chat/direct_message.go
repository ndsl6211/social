package entity

import (
	"time"

	"github.com/google/uuid"
	"mashu.example/internal/entity"
)

type Message struct {
	ID        uuid.UUID
	OwnerId   uuid.UUID
	Content   string
	Timestamp time.Time
}

func NewMessageWithTime(
	id uuid.UUID,
	ownerId uuid.UUID,
	content string,
	time time.Time,
) *Message {
	return &Message{id, ownerId, content, time}
}

func NewMessage(
	id uuid.UUID,
	ownerId uuid.UUID,
	content string,
) *Message {
	return &Message{id, ownerId, content, time.Now()}
}

type DirectMessage struct {
	ID        uuid.UUID
	Creator   *entity.User
	Receiver  *entity.User
	Messages  []*Message
	CreatedAt time.Time
}

func (dm *DirectMessage) AddMessage(senderId uuid.UUID, content string) {
	messageId := uuid.New()
	dm.Messages = append(dm.Messages, NewMessage(messageId, senderId, content))
}

func NewDirectMessage(
	id uuid.UUID,
	creator *entity.User,
	receiver *entity.User,
) *DirectMessage {
	return &DirectMessage{id, creator, receiver, []*Message{}, time.Now()}
}
