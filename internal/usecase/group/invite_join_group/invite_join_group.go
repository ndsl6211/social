package invite_join_group

import (
	"github.com/google/uuid"
	"golang.org/x/exp/slices"
	"mashu.example/internal/entity"
	"mashu.example/internal/usecase"
	"mashu.example/internal/usecase/repository"
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

	Req *InviteJoinGroupUseCaseReq
	Res *InviteJoinGroupUseCaseRes
}

func (gc *InviteJoinGroupUseCase) Execute() {
	invitee, err := gc.userRepo.GetUserById(gc.Req.invitee)
	group, err := gc.groupRepo.GetGroupById(gc.Req.group)
	inviter, err := gc.userRepo.GetUserById(gc.Req.inviter)
	if err != nil {
		gc.Res.Err = err
		return
	}

	if group.Owner.ID == inviter.ID || slices.Contains(group.Admins, inviter.ID) {
		inviteReq := &entity.InviteRequest{invitee.ID, group.ID, inviter.ID}
		group.AddInviteRequests(inviteReq)
	} else {
		return
	}

	gc.groupRepo.Save(group)
	gc.Res.Err = nil
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
