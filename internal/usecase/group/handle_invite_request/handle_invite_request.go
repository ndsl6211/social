package handle_invite_request

import (
	"errors"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"mashu.example/internal/usecase"
	"mashu.example/internal/usecase/repository"
)

var (
	ErrApproverHasNoPermission = errors.New("the approver doesn't have permission")
	ErrInvitationNotFound      = errors.New("the invitation not found")
)

type HandleInviteRequestAction string

const (
	ACCEPT_INVITE_REQUEST HandleInviteRequestAction = "ACCEPT"
	REJECT_INVITE_REQUEST HandleInviteRequestAction = "REJECT"
)

type HandleInviteRequestUseCaseReq struct {
	inviteeId  uuid.UUID
	inviterId  uuid.UUID
	groupId    uuid.UUID
	action     HandleInviteRequestAction
	approverId uuid.UUID
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

func (uc *HandleInviteRequestUseCase) Execute() {
	if _, err := uc.userRepo.GetUserById(uc.req.inviteeId); err != nil {
		uc.res.Err = err
		return
	}
	if _, err := uc.userRepo.GetUserById(uc.req.inviterId); err != nil {
		uc.res.Err = err
		return
	}
	if _, err := uc.userRepo.GetUserById(uc.req.approverId); err != nil {
		uc.res.Err = err
		return
	}
	group, err := uc.groupRepo.GetGroupById(uc.req.groupId)
	if err != nil {
		uc.res.Err = err
		return
	}

	if !group.IsAdmin(uc.req.approverId) && !group.IsOwner(uc.req.approverId) {
		uc.res.Err = ErrApproverHasNoPermission
		logrus.Error(uc.res.Err)
		return
	}

	invitation := group.FindInvitationByInvitee(uc.req.inviteeId)
	if invitation == nil {
		uc.res.Err = ErrInvitationNotFound
		logrus.Error(uc.res.Err)
		return
	}

	if uc.req.action == ACCEPT_INVITE_REQUEST {
		group.AddMember(uc.req.inviteeId, uc.req.inviterId, uc.req.approverId)
	}

	group.RemoveInvitation(invitation.Invitee)

	uc.groupRepo.Save(group)
	uc.res.Err = nil
}

func NewHandleInviteRequestUseCase(
	userRepo repository.UserRepo,
	groupRepo repository.GroupRepo,
	req *HandleInviteRequestUseCaseReq,
	res *HandleInviteRequestUseCaseRes,
) usecase.UseCase {
	return &HandleInviteRequestUseCase{userRepo, groupRepo, req, res}
}

func NewHandleInviteRequestUseCaseReq(
	inviteeId uuid.UUID,
	inviterId uuid.UUID,
	groupId uuid.UUID,
	action HandleInviteRequestAction,
	approverId uuid.UUID,
) HandleInviteRequestUseCaseReq {
	return HandleInviteRequestUseCaseReq{inviteeId, inviterId, groupId, action, approverId}
}

func NewHandleInviteRequestUseCaseRes() HandleInviteRequestUseCaseRes {
	return HandleInviteRequestUseCaseRes{}
}
