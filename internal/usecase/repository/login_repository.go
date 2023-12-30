package repository

import "github.com/google/uuid"

//go:generate mockgen -destination=./mock/login_mock.go -package=mock . LoginRepo
type LoginRepo interface {
	SaveUserSess(userId uuid.UUID) error
	ClearUserSess(userId uuid.UUID) error
}
