package entity

import (
	"github.com/google/uuid"
	"mashu.example/internal/model"
)

type UserEntity struct {
	ID            uuid.UUID           `gorm:"primaryKey;column:id;type:varchat(36)" json:"id"`
	UserName      string              `gorm:"column:name;unique" json:"username"`
	Email         string              `gorm:"column:email" json:"email"`
	Public        bool                `gorm:"column:public" json:"public"`
	DiscordUserId string              `gorm:"column:discord_user_id;unique" json:"discordUserId"`
	Follows       []UserFollowsEntity `gorm:"foreignKey:user_id,follower_id;references:id,id" json:"-"`
}

func (UserEntity) TableName() string {
	return "users"
}

func (u UserEntity) ToUser() *model.User {
	followReqs := []*model.FollowRequest{}
	followers := []uuid.UUID{}
	followings := []uuid.UUID{}

	for _, follow := range u.Follows {
		if follow.Status == USER_FOLLOWS_FOLLOWING && follow.User == u.ID {
			followers = append(followers, follow.Follower)
		} else if follow.Status == USER_FOLLOWS_FOLLOWING && follow.Follower == u.ID {
			followings = append(followings, follow.User)
		} else if follow.Status == USER_FOLLOWS_REQUESTED {
			followReqs = append(followReqs, &model.FollowRequest{From: follow.Follower, To: follow.User})
		}
	}

	return &model.User{
		ID:             u.ID,
		UserName:       u.UserName,
		Email:          u.Email,
		Public:         u.Public,
		DiscordUserId:  u.DiscordUserId,
		FollowRequests: followReqs,
		Followers:      followers,
		Followings:     followings,
	}
}

func NewUserEntity(user *model.User) *UserEntity {
	return &UserEntity{
		ID:            user.ID,
		UserName:      user.UserName,
		Email:         user.Email,
		Public:        user.Public,
		DiscordUserId: user.DiscordUserId,
		Follows:       []UserFollowsEntity{},
	}
}
