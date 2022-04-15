package handle_invite_request

import (
	"github.com/google/uuid"
	"mashu.example/internal/usecase/repository"
)

type HandleInviteRequestAction string

const (
	ACCEPT_INVITE_REQUEST HandleInviteRequestAction = "ACCEPT"
	REJECT_INVITE_REQUEST HandleInviteRequestAction = "REJECT"
)

type HandleInviteRequestUseCaseReq struct {
	inviterId uuid.UUID
	groupId   uuid.UUID
	action    HandleInviteRequestAction
	inviteeId uuid.UUID
}

type HandleInviteRequestUseCaseRes struct {
	Err error
}

type HandleInviteRequestUseCase struct {
	userRepo  repository.UserRepo
	groupRepo repository.GroupRepo
	req       *HandleInviteRequestUseCaseReq
	res       *HandleInviteRequestUseCaseRes
}

func (gc *HandleInviteRequestUseCase) Execute() {

}
