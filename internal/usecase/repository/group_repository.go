package repository

import (
	"github.com/google/uuid"
	"mashu.example/internal/entity"
)

//go:generate mockgen -destination=./mock/group_mock.go -package=mock . GroupRepo
type GroupRepo interface {
	GetGroupById(groupId uuid.UUID) (*entity.Group, error)
	GetGroupByName(groupname string) (*entity.Group, error)
	Save(group *entity.Group) error
	Delete(groupId uuid.UUID) error
}
