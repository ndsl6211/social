package repository

import "mashu.example/internal/entity"

//go:generate mockgen -destination=./mock/user_mock.go -package=mock . UserRepo
type UserRepo interface {
	GetUserById(userId string) (*entity.User, error)
	Save(user *entity.User) error
}
