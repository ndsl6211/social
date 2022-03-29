package entity

import (
	"fmt"

	"github.com/google/uuid"
)

type Post struct {
	ID      uuid.UUID
	Title   string
	Content string
	Owner   *User

	Public bool
}

func (p *Post) Inspect() {
	fmt.Printf("%+v\n", p)
}

func (p *Post) SetOwner(u *User) {
	p.Owner = u
}

func NewPost(id uuid.UUID, title, content string, owner *User, public bool) *Post {
	return &Post{
		ID:      id,
		Title:   title,
		Content: content,
		Owner:   owner,
		Public:  public,
	}
}
