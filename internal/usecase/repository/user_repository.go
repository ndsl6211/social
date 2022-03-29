package repository

import (
	"github.com/google/uuid"
	"mashu.example/internal/entity"
)

//go:generate mockgen -destination=./mock/user_mock.go -package=mock . UserRepo
type UserRepo interface {
	GetUserById(userId uuid.UUID) (*entity.User, error)
	Save(user *entity.User) error
}
