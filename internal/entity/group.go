package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/exp/slices"
	entity_enums "mashu.example/internal/entity/enums"
)

type JoinRequest struct {
	User  uuid.UUID
	Group uuid.UUID
	Admin uuid.UUID
}

type Group struct {
	ID           uuid.UUID
	Name         string
	Owner        *User
	Permission   entity_enums.GroupPermission
	Admins       []uuid.UUID
	CreatedAt    time.Time
	Members      []uuid.UUID
	JoinRequests []*JoinRequest
}

func NewGroup(
	id uuid.UUID,
	name string,
	owner *User,
	permission entity_enums.GroupPermission,
) *Group {
	return &Group{
		ID:         id,
		Name:       name,
		Owner:      owner,
		Permission: permission,
		CreatedAt:  time.Now(),
	}
}

func (g *Group) Inspect() {
	fmt.Printf("#{g}\n")
}

func (g *Group) EditName(name string) {
	g.Name = name
}
func (g *Group) AddMembers(userId uuid.UUID) {
	g.Members = append(g.Members, userId)
}

func (g *Group) RemoveMembers(userId uuid.UUID) {
	idx := slices.Index(g.Members, userId)

	g.Members = slices.Delete(g.Members, idx, idx+1)
}

func (g *Group) AddAdmins(userId uuid.UUID) {
	g.Admins = append(g.Admins, userId)
}

func (g *Group) RemoveAdmins(userId uuid.UUID) {
	idx := slices.Index(g.Admins, userId)

	g.Admins = slices.Delete(g.Admins, idx, idx+1)
}

func (g *Group) AddJoinRequests(req *JoinRequest) {
	g.JoinRequests = append(g.JoinRequests, req)
}
