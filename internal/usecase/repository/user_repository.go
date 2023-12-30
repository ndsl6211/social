package repository

import (
	"fmt"

	"github.com/google/uuid"
	"mashu.example/internal/model"
)

type ErrUserNotFound struct {
	UserId uuid.UUID
}

func (err *ErrUserNotFound) Error() string {
	return fmt.Sprintf("User %s not found", err.UserId.String())
}

func NewErrUserNotFound(userId uuid.UUID) *ErrUserNotFound {
	return &ErrUserNotFound{userId}
}

//go:generate mockgen -destination=./mock/user_mock.go -package=mock . UserRepo
type UserRepo interface {
	GetUserById(userId uuid.UUID) (*model.User, error)
	GetUserByUserName(userName string) (*model.User, error)
	GetUserByDiscordUserId(discordUserId string) (*model.User, error)
	Save(user *model.User) error
}
