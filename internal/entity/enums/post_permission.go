package entity_enums

type PostPermission int

const (
	POST_PUBLIC        PostPermission = iota
	POST_FOLLOWER_ONLY PostPermission = iota
	POST_PRIVATE       PostPermission = iota
)
