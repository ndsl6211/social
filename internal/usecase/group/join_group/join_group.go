package join_group

import "github.com/google/uuid"

type JoinGroupUseCaseReq struct {
	userId  uuid.UUID
	groupId uuid.UUID
}
