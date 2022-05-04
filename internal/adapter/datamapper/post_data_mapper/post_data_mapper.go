package post_data_mapper

import (
	"time"

	"mashu.example/internal/adapter/datamapper/user_data_mapper"
	entity_enums "mashu.example/internal/entity/enums"

	"github.com/google/uuid"
	"mashu.example/internal/entity"
)

type CommentDataMapper struct {
	ID uuid.UUID `gorm:"primaryKey;column:id"`

	OwnerId uuid.UUID
	Owner   *user_data_mapper.UserDataMapper

	PostId uuid.UUID
	Post   *PostDataMapper

	Content   string    `gorm:"column:content"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (CommentDataMapper) TableName() string {
	return "comments"
}

func (c CommentDataMapper) ToComment() *entity.Comment {
	return &entity.Comment{
		ID:        c.ID,
		Owner:     c.Owner.ToUser(),
		Post:      c.Post.ToPost(),
		Content:   c.Content,
		CreatedAt: c.CreatedAt,
	}
}

func NewCommentDataMapper(comment *entity.Comment) *CommentDataMapper {
	return &CommentDataMapper{
		ID:        comment.ID,
		OwnerId:   comment.Owner.ID,
		PostId:    comment.Post.ID,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt,
	}
}

type PostDataMapper struct {
	ID         uuid.UUID                   `gorm:"primaryKey"`
	Title      string                      `gorm:"column:title"`
	Content    string                      `gorm:"column:content"`
	Permission entity_enums.PostPermission `gorm:"public"`

	OwnerId uuid.UUID
	Owner   *user_data_mapper.UserDataMapper `gorm:"foreignKey:OwnerId"`

	Comments []*CommentDataMapper `gorm:"foreignKey:PostId"`
	CreateAt time.Time            `gorm:"column:created_at"`
}

func (PostDataMapper) TableName() string {
	return "posts"
}

func (p PostDataMapper) ToPost() *entity.Post {
	return &entity.Post{
		ID:         p.ID,
		Title:      p.Title,
		Content:    p.Content,
		Owner:      p.Owner.ToUser(),
		Permission: p.Permission,
	}
}

func NewPostDataMapper(post *entity.Post) *PostDataMapper {
	var comments []*CommentDataMapper
	for _, comment := range post.Comments {
		comments = append(comments, NewCommentDataMapper(comment))
	}

	return &PostDataMapper{
		ID:         post.ID,
		Title:      post.Title,
		Content:    post.Content,
		Permission: post.Permission,
		OwnerId:    post.Owner.ID,
		Owner:      user_data_mapper.NewUserDataMapper(post.Owner),
		Comments:   comments,
		CreateAt:   post.CreatedAt,
	}
}
