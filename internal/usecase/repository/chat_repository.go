package repository

import (
	"fmt"

	"github.com/google/uuid"
	chat "mashu.example/internal/entity/chat"
)

type ErrDMNotFound struct {
	dmId uuid.UUID
}

func (err *ErrDMNotFound) Error() string {
	return fmt.Sprintf("Direct message %s not found", err.dmId.String())
}

//go:generate mockgen -destination=./mock/chat_mock.go -package=mock . ChatRepo
type ChatRepo interface {
	GetDirectMessage(dmId uuid.UUID) (*chat.DirectMessage, error)
	GetDMByUserId(userA uuid.UUID, userB uuid.UUID) (*chat.DirectMessage, error)
	SaveDirectMessage(dm *chat.DirectMessage) error
}
