package user_data_mapper

import (
	"github.com/google/uuid"
	"mashu.example/internal/entity"
)

type UserDataMapper struct {
	ID          uuid.UUID          `gorm:"primaryKey;column:id;type:varchat(36)" json:"id"`
	UserName    string             `gorm:"column:name;unique" json:"username"`
	DisplayName string             `gorm:"column:display_name" json:"displayName"`
	Email       string             `gorm:"column:email" json:"email"`
	Public      bool               `gorm:"column:public" json:"public"`
	Follows     []FollowDataMapper `gorm:"foreignKey:user_id,follower_id;references:id,id" json:"-"`
}

func (UserDataMapper) TableName() string {
	return "users"
}

func (u UserDataMapper) ToUser() *entity.User {
	followReqs := []*entity.FollowRequest{}
	followers := []uuid.UUID{}
	followings := []uuid.UUID{}

	for _, follow := range u.Follows {
		if follow.Status == FOLLOWING && follow.User == u.ID {
			// find my followers
			followers = append(followers, follow.Follower)
		} else if follow.Status == FOLLOWING && follow.Follower == u.ID {
			// find my following users
			followings = append(followings, follow.User)
		} else if follow.Status == REQUESTED {
			followReqs = append(followReqs, &entity.FollowRequest{From: follow.Follower, To: follow.User})
		}
	}

	return &entity.User{
		ID:             u.ID,
		UserName:       u.UserName,
		DisplayName:    u.DisplayName,
		Email:          u.Email,
		Public:         u.Public,
		FollowRequests: followReqs,
		Followers:      followers,
		Followings:     followings,
	}
}

func NewUserDataMapper(user *entity.User) *UserDataMapper {
	return &UserDataMapper{user.ID, user.UserName, user.DisplayName, user.Email, user.Public, []FollowDataMapper{}}
}

type FollowStatus string

const (
	REQUESTED FollowStatus = "REQUESTED"
	FOLLOWING FollowStatus = "FOLLOWING"
)

type FollowDataMapper struct {
	User     uuid.UUID `gorm:"primaryKey;column:user_id;"`
	Follower uuid.UUID `gorm:"primaryKey;column:follower_id;"`
	Status   FollowStatus
}

func (FollowDataMapper) TableName() string {
	return "follows"
}

func NewFollowDataMapper(
	userId uuid.UUID,
	followerId uuid.UUID,
	status FollowStatus,
) *FollowDataMapper {
	return &FollowDataMapper{userId, followerId, status}
}
