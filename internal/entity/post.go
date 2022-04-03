package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
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
	ID      uuid.UUID
	Title   string
	Content string
	Owner   *User
	Public  bool

	Comments []*Comment
}

func (p *Post) Inspect() {
	fmt.Printf("%+v\n", p)
}

func (p *Post) SetOwner(u *User) {
	p.Owner = u
}

func NewPost(id uuid.UUID, title, content string, owner *User, public bool) *Post {
	return &Post{
		ID:       id,
		Title:    title,
		Content:  content,
		Owner:    owner,
		Public:   public,
		Comments: []*Comment{},
	}
}
