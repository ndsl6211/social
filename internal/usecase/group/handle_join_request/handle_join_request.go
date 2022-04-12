package handle_join_request

import (
	"errors"
	"github.com/google/uuid"
	"golang.org/x/exp/slices"
	"mashu.example/internal/entity"
	"mashu.example/internal/usecase"
	"mashu.example/internal/usecase/repository"
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
	Req       *HandleJoinRequestUseCaseReq
	Res       *HandleJoinRequestUseCaseRes
}

func (gc *HandleJoinRequestUseCase) Execute() {
	requester, err := gc.userRepo.GetUserById(gc.Req.requesterId)
	group, err := gc.groupRepo.GetGroupById(gc.Req.groupId)
	approver, err := gc.userRepo.GetUserById(gc.Req.approverId)
	if err != nil {
		gc.Res.Err = err
		return
	}
	if !slices.Contains(group.Admins, approver.ID) && approver != group.Owner {
		errMsg := "permission denied"
		gc.Res.Err = errors.New(errMsg)
		return
	}

	idx := slices.IndexFunc(group.JoinRequests, func(req *entity.JoinRequest) bool {
		return req.Group == gc.Req.groupId && req.Requester == gc.Req.requesterId
	})

	if idx < 0 {
		errMsg := "request not found"
		gc.Res.Err = errors.New(errMsg)
		return
	}

	if gc.Req.action == ACCEPT_JOIN_REQUEST {
		group.AddMembers(requester.ID)
	}

	group.JoinRequests = slices.Delete(group.JoinRequests, idx, idx+1)

	gc.groupRepo.Save(group)
	gc.Res.Err = nil
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
