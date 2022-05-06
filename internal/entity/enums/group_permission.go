package entity_enums

type GroupPermission string

const (
	GROUP_PUBLIC   GroupPermission = "PUBLIC"
	GROUP_UNPUBLIC GroupPermission = "UNPUBLIC"
	GROUP_PRIVATE  GroupPermission = "PRIVATE"
)
