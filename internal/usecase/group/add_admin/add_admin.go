package add_admin

import (
	"errors"
	"github.com/google/uuid"
	"golang.org/x/exp/slices"
	"mashu.example/internal/usecase"
	"mashu.example/internal/usecase/repository"
)

type AddAdminUseCaseReq struct {
	memberId uuid.UUID
	groupId  uuid.UUID
	adminId  uuid.UUID
}

type AddAdminUseCaseRes struct {
	Err error
}

type AddAdminUseCase struct {
	groupRepo repository.GroupRepo
	userRepo  repository.UserRepo

	req *AddAdminUseCaseReq
	res *AddAdminUseCaseRes
}

func (gc *AddAdminUseCase) Execute() {
	member, err := gc.userRepo.GetUserById(gc.req.memberId)
	group, err := gc.groupRepo.GetGroupById(gc.req.groupId)
	admin, err := gc.userRepo.GetUserById(gc.req.adminId)
	if err != nil {
		gc.res.Err = err
		return
	}

	if !slices.Contains(group.Admins, admin.ID) && admin != group.Owner {
		errMsg := "permission denied"
		gc.res.Err = errors.New(errMsg)
		return
	}

	if !slices.Contains(group.Members, member.ID) {
		errMsg := "permission denied"
		gc.res.Err = errors.New(errMsg)
		return
	}
	group.AddAdmins(member.ID)
	gc.groupRepo.Save(group)
	gc.res.Err = nil
}

func NewAddAdminUseCase(
	groupRepo repository.GroupRepo,
	userRepo repository.UserRepo,
	req *AddAdminUseCaseReq,
	res *AddAdminUseCaseRes,
) usecase.UseCase {
	return &AddAdminUseCase{groupRepo, userRepo, req, res}
}

func NewAddAdminUseCaseReq(
	memberId uuid.UUID,
	groupId uuid.UUID,
	adminId uuid.UUID,
) AddAdminUseCaseReq {
	return AddAdminUseCaseReq{memberId, groupId, adminId}
}

func NewAddAdminUseCaseRes() AddAdminUseCaseRes {
	return AddAdminUseCaseRes{}
}
