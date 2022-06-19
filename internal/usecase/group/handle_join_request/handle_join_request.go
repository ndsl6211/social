package handle_join_request

import (
	"errors"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"mashu.example/internal/usecase"
	"mashu.example/internal/usecase/repository"
)

var (
	ErrApproverHasNoPermission = errors.New("the approver doesn't have permission")
	ErrJoinRequestNotFound     = errors.New("the join request not found")
)

type HandleJoinRequestAction string

const (
	ACCEPT_JOIN_REQUEST HandleJoinRequestAction = "ACCEPT"
	REJECT_JOIN_REQUEST HandleJoinRequestAction = "REJECT"
)

type HandleJoinRequestUseCaseReq struct {
	requesterId uuid.UUID
	groupId     uuid.UUID
	action      HandleJoinRequestAction
	approverId  uuid.UUID
}

type HandleJoinRequestUseCaseRes struct {
	Err error
}

type HandleJoinRequestUseCase struct {
	userRepo  repository.UserRepo
	groupRepo repository.GroupRepo
	req       *HandleJoinRequestUseCaseReq
	res       *HandleJoinRequestUseCaseRes
}

func (uc *HandleJoinRequestUseCase) Execute() {
	_, err := uc.userRepo.GetUserById(uc.req.requesterId)
	if err != nil {
		uc.res.Err = err
		return
	}
	group, err := uc.groupRepo.GetGroupById(uc.req.groupId)
	if err != nil {
		uc.res.Err = err
		return
	}
	_, err = uc.userRepo.GetUserById(uc.req.approverId)
	if err != nil {
		uc.res.Err = err
		return
	}

	if !group.IsAdmin(uc.req.approverId) && !group.IsOwner(uc.req.approverId) {
		uc.res.Err = ErrApproverHasNoPermission
		logrus.Error(uc.res.Err)
		return
	}

	joinRequest := group.FindJoinRequest(uc.req.requesterId)
	if joinRequest == nil {
		uc.res.Err = ErrJoinRequestNotFound
		logrus.Error(uc.res.Err)
		return
	}

	if uc.req.action == ACCEPT_JOIN_REQUEST {
		group.AddMember(uc.req.requesterId, uuid.Nil, uc.req.approverId)
	}

	group.RemoveJoinRequest(uc.req.requesterId)

	uc.groupRepo.Save(group)
	uc.res.Err = nil
}

func NewHandleJoinRequestUseCase(
	userRepo repository.UserRepo,
	groupRepo repository.GroupRepo,
	req *HandleJoinRequestUseCaseReq,
	res *HandleJoinRequestUseCaseRes,
) usecase.UseCase {
	return &HandleJoinRequestUseCase{userRepo, groupRepo, req, res}
}

func NewHandleJoinRequestUseCaseReq(
	requesterId uuid.UUID,
	groupId uuid.UUID,
	action HandleJoinRequestAction,
	approverId uuid.UUID,
) HandleJoinRequestUseCaseReq {
	return HandleJoinRequestUseCaseReq{requesterId, groupId, action, approverId}
}

func NewHandleJoinRequestUseCaseRes() HandleJoinRequestUseCaseRes {
	return HandleJoinRequestUseCaseRes{}
}
