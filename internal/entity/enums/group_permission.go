package entity_enums

type GroupPrivacy string

// GROUP_PUBLIC - all users can view the posts and members of the group
// GROUP_PRIVATE - only the group member can view the posts and members of this
//                 group
const (
	GROUP_PUBLIC  GroupPrivacy = "PUBLIC"
	GROUP_PRIVATE GroupPrivacy = "PRIVATE"
)

type GroupVisibility string

// GROUP_VISIBLE - all users can find the group
// GROUP_HIDDEN - only the group members can find the group
const (
	GROUP_VISIBLE GroupVisibility = "VISIBLE"
	GROUP_HIDDEN  GroupVisibility = "HIDDEN"
)
