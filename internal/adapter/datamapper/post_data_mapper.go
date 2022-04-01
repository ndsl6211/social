package datamapper

import (
	"github.com/google/uuid"
	"mashu.example/internal/entity"
)

type PostDataMapper struct {
	ID      uuid.UUID `gorm:"primaryKey"`
	Title   string    `gorm:"column:title"`
	Content string    `gorm:"column:content"`
	OwnerId uuid.UUID
	Owner   UserDataMapper `gorm:"foreignKey:OwnerId"`
	Public  bool           `gorm:"public"`
}

func (PostDataMapper) TableName() string {
	return "posts"
}

func (p PostDataMapper) ToPost() *entity.Post {
	return &entity.Post{
		ID:      p.ID,
		Title:   p.Title,
		Content: p.Content,
		Owner:   p.Owner.ToUser(),
		Public:  p.Public,
	}
}

func NewPostDataMapper(post *entity.Post) *PostDataMapper {
	return &PostDataMapper{
		ID:      post.ID,
		Title:   post.Title,
		Content: post.Content,
		OwnerId: post.Owner.ID,
		Owner:   *NewUserDataMapper(post.Owner),
		Public:  post.Public,
	}
}
