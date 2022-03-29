package adapter

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"mashu.example/internal/adapter/datamapper"
	"mashu.example/internal/entity"
	"mashu.example/internal/usecase/repository"
)

type postRepo struct {
	db *gorm.DB
}

func (pr *postRepo) GetPostById(postId uuid.UUID) (*entity.Post, error) {
	// get post
	postData := datamapper.PostDataMapper{}
	if err := pr.db.
		Where("posts.id = ?", postId).
		Find(postData).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	return postData.ToPost(), nil
}

func (pr *postRepo) Save(post *entity.Post) error {
	return nil
}

func NewPostRepository(db *gorm.DB) repository.PostRepo {
	err := db.AutoMigrate(&datamapper.PostDataMapper{})
	if err != nil {
		fmt.Println(err.Error())
	}

	return &postRepo{db}
}
