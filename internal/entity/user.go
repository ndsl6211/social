package entity

import (
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/exp/slices"
)

type FollowRequest struct {
	From uuid.UUID
	To   uuid.UUID
}

type User struct {
	ID          uuid.UUID
	UserName    string
	DisplayName string
	Email       string
	Public      bool

	Followers      []uuid.UUID
	Followings     []uuid.UUID
	FollowRequests []*FollowRequest
}

func (u *User) Inspect() {
	fmt.Printf("%+v\n", u)
}

func (u *User) AddFollower(userId uuid.UUID) {
	u.Followers = append(u.Followers, userId)
}

func (u *User) RemoveFollower(userId uuid.UUID) {
	idx := slices.Index(u.Followers, userId)

	u.Followers = slices.Delete(u.Followers, idx, idx+1)
}

func (u *User) AddFollowing(userId uuid.UUID) {
	u.Followings = append(u.Followings, userId)
}

func (u *User) RemoveFollowing(userId uuid.UUID) {
	idx := slices.Index(u.Followings, userId)

	u.Followings = slices.Delete(u.Followings, idx, idx+1)
}

func (u *User) AddFollowRequest(req *FollowRequest) {
	u.FollowRequests = append(u.FollowRequests, req)
}

func NewUser(id uuid.UUID, userName, displayName, email string, public bool) *User {
	return &User{
		ID:             id,
		UserName:       userName,
		DisplayName:    displayName,
		Email:          email,
		Public:         public,
		Followers:      []uuid.UUID{},
		Followings:     []uuid.UUID{},
		FollowRequests: []*FollowRequest{},
	}
}
