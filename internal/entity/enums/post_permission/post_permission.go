package post_permission

type PostPermission string

const (
	PUBLIC        PostPermission = "PUBLIC"
	FOLLOWER_ONLY PostPermission = "FOLLOWER_ONLY"
	PRIVATE       PostPermission = "PRIVATE"
)
