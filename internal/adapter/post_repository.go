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
	postData := datamapper.PostDataMapper{ID: postId}
	if err := pr.db.
		Preload("Owner").
		First(&postData).Error; err != nil {
		return nil, err
	}

	return postData.ToPost(), nil
}

func (pr *postRepo) GetPostByUserId(userId uuid.UUID) ([]*entity.Post, error) {
	postDataMappers := []*datamapper.PostDataMapper{}
	if err := pr.db.
		Where("posts.owner_id = ?", userId).
		Find(postDataMappers).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	posts := []*entity.Post{}
	for _, post := range postDataMappers {
		posts = append(posts, post.ToPost())
	}

	return posts, nil
}

func (pr *postRepo) Save(post *entity.Post) error {
	postDataMapper := datamapper.NewPostDataMapper(post)
	if err := pr.db.Save(postDataMapper).Error; err != nil {
		return err
	}

	return nil
}

func (pr *postRepo) Delete(postId uuid.UUID) error {
	if err := pr.db.Delete(&entity.Post{ID: postId}).Error; err != nil {
		return err
	}

	return nil
}

func NewPostRepository(db *gorm.DB) repository.PostRepo {
	if err := db.AutoMigrate(&datamapper.PostDataMapper{}); err != nil {
		fmt.Println(err.Error())
	}
	if err := db.AutoMigrate(&datamapper.CommentDataMapper{}); err != nil {
		fmt.Println(err.Error())
	}

	return &postRepo{db}
}
