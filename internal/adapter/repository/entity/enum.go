package entity

type UserFollowsStatus string

const (
	USER_FOLLOWS_REQUESTED UserFollowsStatus = "REQUESTED"
	USER_FOLLOWS_FOLLOWING UserFollowsStatus = "FOLLOWING"
)
