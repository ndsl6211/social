package handle_join_request

import (
	"github.com/google/uuid"
	"mashu.example/internal/usecase/repository"
)

type HandleJoinRequestAction string

const (
	ACCEPT_JOIN_REQUEST HandleJoinRequestAction = "ACCEPT"
	REJECT_JOIN_REQUEST HandleJoinRequestAction = "REJECT"
)

type HandleJoinRequestUseCaseReq struct {
	userId   uuid.UUID
	groupId  uuid.UUID
	action   HandleJoinRequestAction
	approver uuid.UUID
}

type HandleJoinRequestUseCaseRes struct {
	Err error
}

type HandleJoinRequestUseCase struct {
	userRepo  repository.UserRepo
	groupRepo repository.GroupRepo
	Req       *HandleJoinRequestUseCaseReq
	Res       *HandleJoinRequestUseCaseRes
}

func (gc *HandleJoinRequestUseCase) Execute() {
	user, err := gc.userRepo.GetUserById(gc.Req.userId)
	group, err := gc.groupRepo.GetGroupById(gc.Req.groupId)
	if err != nil {
		gc.Res.Err = err
		return
	}

}
