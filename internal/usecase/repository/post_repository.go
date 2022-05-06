package repository

import (
	"fmt"

	"github.com/google/uuid"
	"mashu.example/internal/entity"
)

type ErrPostNotFound struct {
	PostId uuid.UUID
}

func (e *ErrPostNotFound) Error() string {
	return fmt.Sprintf("Post %s not found", e.PostId)
}

//go:generate mockgen -destination=./mock/post_mock.go -package=mock . PostRepo
type PostRepo interface {
	GetPostById(postId uuid.UUID) (*entity.Post, error)
	GetPostByUserId(userId uuid.UUID) ([]*entity.Post, error)
	Save(post *entity.Post) error
	Delete(postId uuid.UUID) error
}
