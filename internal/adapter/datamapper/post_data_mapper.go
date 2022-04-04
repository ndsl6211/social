package datamapper

import (
	"time"

	"github.com/google/uuid"
	"mashu.example/internal/entity"
)

type CommentDataMapper struct {
	ID uuid.UUID `gorm:"primaryKey;column:id"`

	OwnerId uuid.UUID
	Owner   *UserDataMapper

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
	ID      uuid.UUID `gorm:"primaryKey"`
	Title   string    `gorm:"column:title"`
	Content string    `gorm:"column:content"`
	Public  bool      `gorm:"public"`

	OwnerId uuid.UUID
	Owner   *UserDataMapper `gorm:"foreignKey:OwnerId"`

	Comments []*CommentDataMapper `gorm:"foreignKey:PostId"`
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
	var comments []*CommentDataMapper
	for _, comment := range post.Comments {
		comments = append(comments, NewCommentDataMapper(comment))
	}

	return &PostDataMapper{
		ID:       post.ID,
		Title:    post.Title,
		Content:  post.Content,
		Public:   post.Public,
		OwnerId:  post.Owner.ID,
		Owner:    NewUserDataMapper(post.Owner),
		Comments: comments,
	}
}
