package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	entity_enums "mashu.example/internal/entity/enums"
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
	Permission entity_enums.PostPermission

	Comments []*Comment

	CreatedAt time.Time
	UpdatedAt time.Time
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
	permission entity_enums.PostPermission,
) *Post {
	return &Post{
		ID:         id,
		Title:      title,
		Content:    content,
		Owner:      owner,
		Permission: permission,
		Comments:   []*Comment{},
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}
