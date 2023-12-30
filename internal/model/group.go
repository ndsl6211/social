package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/exp/slices"
	entity_enums "mashu.example/internal/entity/enums"
)

type JoinRequest struct {
	Requester uuid.UUID
	timestamp time.Time
}

type InviteRequest struct {
	Invitee uuid.UUID
	Inviter uuid.UUID
}

// struct to represent an admin in group
// note:
// - the promoter can only be admin or owner
type GroupAdmin struct {
	UserId     uuid.UUID
	PromotedBy uuid.UUID
	PromotedAt time.Time
}

// struct to represent a member in group
// rules:
// - the inviter can be any member in this group
// - the approver can only be admin/owner of this group
// - if the user actively send the join request, the `invitedBy` will be nil
// - if the user is invited by group member, the the `invitedBy` will be the member
// - if the group is public, anyone can send join request
// - if the group is private, user can only be invited into this group

type GroupMember struct {
	UserId     uuid.UUID
	InvitedBy  uuid.UUID
	ApprovedBy uuid.UUID
	JoinAt     time.Time
}

type Group struct {
	ID         uuid.UUID
	Name       string
	Owner      *User
	Permission entity_enums.GroupPrivacy
	CreatedAt  time.Time

	Admins         []*GroupAdmin
	Members        []*GroupMember
	JoinRequests   []*JoinRequest
	InviteRequests []*InviteRequest
}

func (g *Group) Inspect() {
	fmt.Printf("#{g}\n")
}

func (g *Group) EditName(name string) {
	g.Name = name
}

func (g *Group) AddMember(
	userId uuid.UUID,
	inviter uuid.UUID,
	approver uuid.UUID,
) {
	g.Members = append(g.Members, &GroupMember{
		UserId:     userId,
		InvitedBy:  inviter,
		ApprovedBy: approver,
		JoinAt:     time.Now(),
	})
}

func (g *Group) RemoveMember(userId uuid.UUID) {
	idx := slices.IndexFunc(g.Members, func(member *GroupMember) bool {
		return member.UserId == userId
	})

	if idx == -1 {
		return
	}

	g.Members = slices.Delete(g.Members, idx, idx+1)
}

func (g *Group) AddAdmin(
	userId uuid.UUID,
	promoter uuid.UUID,
) {
	g.Admins = append(g.Admins, &GroupAdmin{
		UserId:     userId,
		PromotedBy: promoter,
		PromotedAt: time.Now(),
	})
}

func (g *Group) RemoveAdmin(userId uuid.UUID) {
	idx := slices.IndexFunc(g.Admins, func(admin *GroupAdmin) bool {
		return admin.UserId == userId
	})

	if idx == -1 {
		return
	}

	g.Admins = slices.Delete(g.Admins, idx, idx+1)
}

func (g *Group) IsAdmin(userId uuid.UUID) bool {
	return slices.IndexFunc(g.Admins, func(admin *GroupAdmin) bool {
		return admin.UserId == userId
	}) != -1
}

func (g *Group) IsOwner(userId uuid.UUID) bool {
	return userId == g.Owner.ID
}

func (g *Group) AddJoinRequest(requesterId uuid.UUID) {
	g.JoinRequests = append(g.JoinRequests, &JoinRequest{requesterId, time.Now()})
}

func (g *Group) FindJoinRequest(requesterId uuid.UUID) *JoinRequest {
	if idx := slices.IndexFunc(g.JoinRequests, func(req *JoinRequest) bool {
		return requesterId == req.Requester
	}); idx != -1 {
		return g.JoinRequests[idx]
	} else {
		return nil
	}
}

func (g *Group) RemoveJoinRequest(requesterId uuid.UUID) {
	idx := slices.IndexFunc(g.JoinRequests, func(req *JoinRequest) bool {
		return requesterId == req.Requester
	})
	if idx != -1 {
		g.JoinRequests = slices.Delete(g.JoinRequests, idx, idx+1)
	}
}

func (g *Group) AddInviteRequest(invitee uuid.UUID, inviter uuid.UUID) {
	g.InviteRequests = append(g.InviteRequests, &InviteRequest{invitee, inviter})
}

func (g *Group) FindInvitationByInvitee(inviteeId uuid.UUID) *InviteRequest {
	if idx := slices.IndexFunc(g.InviteRequests, func(req *InviteRequest) bool {
		return inviteeId == req.Invitee
	}); idx != -1 {
		return g.InviteRequests[idx]
	} else {
		return nil
	}
}

func (g *Group) RemoveInvitation(inviteeId uuid.UUID) {
	idx := slices.IndexFunc(g.InviteRequests, func(req *InviteRequest) bool {
		return inviteeId == req.Invitee
	})
	if idx != -1 {
		g.InviteRequests = slices.Delete(g.InviteRequests, idx, idx+1)
	}
}

func NewGroup(
	id uuid.UUID,
	name string,
	owner *User,
	permission entity_enums.GroupPrivacy,
) *Group {
	return &Group{
		ID:         id,
		Name:       name,
		Owner:      owner,
		Permission: permission,
		CreatedAt:  time.Now(),
	}
}
