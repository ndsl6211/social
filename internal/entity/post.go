package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"mashu.example/internal/entity/enums/post_permission"
)

type Comment struct {
	ID        uuid.UUID
	Owner     *User
	Post      *Post
	Content   string
	CreatedAt time.Time
}

func NewComment(id uuid.UUID, owner *User, post *Post, content string) *Comment {
	return &Comment{id, owner, post, content, time.Now()}
}

type Post struct {
	ID         uuid.UUID
	Title      string
	Content    string
	Owner      *User
	Permission post_permission.PostPermission

	Comments []*Comment
}

func (p *Post) Inspect() {
	fmt.Printf("%+v\n", p)
}

func (p *Post) SetOwner(u *User) {
	p.Owner = u
}

func NewPost(
	id uuid.UUID,
	title string,
	content string,
	owner *User,
	permission post_permission.PostPermission,
) *Post {
	return &Post{
		ID:         id,
		Title:      title,
		Content:    content,
		Owner:      owner,
		Permission: permission,
		Comments:   []*Comment{},
	}
}
