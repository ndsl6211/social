package invite_join_group

import (
	"errors"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
	"mashu.example/internal/entity"
	"mashu.example/internal/usecase"
	"mashu.example/internal/usecase/repository"
)

var (
	ErrInviteeIsAlreadyMember  = errors.New("the invitee is already a member of the group")
	ErrInviterIsNotMember      = errors.New("the inviter is not a member of the group")
	ErrInviteeIsAlreadyInvited = errors.New("the invitee has been invited")
)

type InviteJoinGroupUseCaseReq struct {
	invitee uuid.UUID
	group   uuid.UUID
	inviter uuid.UUID
}

type InviteJoinGroupUseCaseRes struct {
	Err error
}

type InviteJoinGroupUseCase struct {
	userRepo  repository.UserRepo
	groupRepo repository.GroupRepo

	req *InviteJoinGroupUseCaseReq
	res *InviteJoinGroupUseCaseRes
}

func (uc *InviteJoinGroupUseCase) Execute() {
	_, err := uc.userRepo.GetUserById(uc.req.invitee)
	if err != nil {
		uc.res.Err = err
		return
	}
	_, err = uc.userRepo.GetUserById(uc.req.inviter)
	if err != nil {
		uc.res.Err = err
		return
	}
	group, err := uc.groupRepo.GetGroupById(uc.req.group)
	if err != nil {
		uc.res.Err = err
		return
	}

	if !uc.isMember(group, uc.req.inviter) && uc.req.inviter != group.Owner.ID {
		uc.res.Err = ErrInviterIsNotMember
		logrus.Error(uc.res.Err)
		return
	}

	if uc.isMember(group, uc.req.invitee) {
		uc.res.Err = ErrInviteeIsAlreadyMember
		logrus.Error(uc.res.Err)
		return
	}

	if uc.isInviting(group, uc.req.invitee) {
		uc.res.Err = ErrInviteeIsAlreadyInvited
		logrus.Error(uc.res.Err)
		return
	}

	group.AddInviteRequest(uc.req.invitee, uc.req.inviter)

	uc.groupRepo.Save(group)
	uc.res.Err = nil
}

func (uc *InviteJoinGroupUseCase) isMember(group *entity.Group, userId uuid.UUID) bool {
	return slices.IndexFunc(group.Members, func(member *entity.GroupMember) bool {
		return member.UserId == userId
	}) != -1
}

func (uc *InviteJoinGroupUseCase) isInviting(group *entity.Group, inviteeId uuid.UUID) bool {
	return slices.IndexFunc(group.InviteRequests, func(invite *entity.InviteRequest) bool {
		return invite.Invitee == inviteeId
	}) != -1
}

func NewInviteJoinGroupUseCase(
	userRepo repository.UserRepo,
	groupRepo repository.GroupRepo,
	req *InviteJoinGroupUseCaseReq,
	res *InviteJoinGroupUseCaseRes,
) usecase.UseCase {
	return &InviteJoinGroupUseCase{userRepo, groupRepo, req, res}
}

func NewInviteJoinGroupUseCaseReq(invitee, group, inviter uuid.UUID) InviteJoinGroupUseCaseReq {
	return InviteJoinGroupUseCaseReq{invitee, group, inviter}
}

func NewInviteJoinGroupUseCaseRes() InviteJoinGroupUseCaseRes {
	return InviteJoinGroupUseCaseRes{}
}
