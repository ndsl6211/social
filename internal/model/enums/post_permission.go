package entity_enums

type PostPermission string

const (
	POST_PUBLIC        PostPermission = "PUBLIC"
	POST_FOLLOWER_ONLY PostPermission = "FOLLER_ONLY"
	POST_PRIVATE       PostPermission = "PRIVATE"
)
