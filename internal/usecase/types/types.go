package types

import "github.com/google/uuid"

type FollowingInfo struct {
	ID          uuid.UUID
	UserName    string
	DisplayName string
	Email       string
	Public      bool
}
