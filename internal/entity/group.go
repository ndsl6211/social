package entity

import (
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/exp/slices"
	"mashu.example/internal/entity/enums/group_permission"
	"time"
)

type JoinRequest struct {
	User  uuid.UUID
	Group uuid.UUID
}

type InviteRequest struct {
	Invitee uuid.UUID
	Group   uuid.UUID
	Inviter uuid.UUID
}

type Group struct {
	ID             uuid.UUID
	Name           string
	Owner          *User
	Permission     group_permission.GroupPermission
	Admins         []uuid.UUID
	CreatedAt      time.Time
	Members        []uuid.UUID
	JoinRequests   []*JoinRequest
	InviteRequests []*InviteRequest
}

func NewGroup(
	id uuid.UUID,
	name string,
	owner *User,
	permission group_permission.GroupPermission,
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

func (g *Group) AddInviteRequests(req *InviteRequest) {
	g.InviteRequests = append(g.InviteRequests, req)
}
