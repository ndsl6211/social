package group_permission

type GroupPermission string

const (
	PUBLIC   GroupPermission = "PUBLIC"
	UNPUBLIC GroupPermission = "UNPUBLIC"
	PRIVATE  GroupPermission = "PRIVATE"
)
