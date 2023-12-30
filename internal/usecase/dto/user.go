package dto

import "github.com/google/uuid"

type FollowingInfo struct {
	ID       uuid.UUID
	UserName string
	Email    string
	Public   bool
}
