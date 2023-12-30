package entity

import (
	"github.com/google/uuid"
)

type UserFollowsEntity struct {
	User     uuid.UUID `gorm:"primaryKey;column:user_id;"`
	Follower uuid.UUID `gorm:"primaryKey;column:follower_id;"`
	Status   UserFollowsStatus
}

func (UserFollowsEntity) TableName() string {
	return "user_follows"
}

func NewUserFollowsEntity(
	userId uuid.UUID,
	followerId uuid.UUID,
	status UserFollowsStatus,
) *UserFollowsEntity {
	return &UserFollowsEntity{userId, followerId, status}
}
