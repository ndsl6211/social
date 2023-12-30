package model

import "github.com/google/uuid"

type FollowRequest struct {
	From uuid.UUID
	To   uuid.UUID
}
