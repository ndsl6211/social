package add_admin

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
	ErrNotGroupOwnerOrAdmin = errors.New("only the group owner or admin can add admin")
	ErrNotGroupMember       = errors.New("user is not the group member")
	ErrIsAlreadyAdmin       = errors.New("the user is already admin")
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

func (uc *AddAdminUseCase) Execute() {
	_, err := uc.userRepo.GetUserById(uc.req.memberId)
	group, err := uc.groupRepo.GetGroupById(uc.req.groupId)
	_, err = uc.userRepo.GetUserById(uc.req.adminId)
	if err != nil {
		uc.res.Err = err
		return
	}

	// check whether `uc.req.adminId` is admin
	isNotAdmin := (slices.IndexFunc(group.Admins, func(admin *entity.GroupAdmin) bool {
		return admin.UserId == uc.req.adminId
	}) == -1)
	if uc.req.adminId != group.Owner.ID && isNotAdmin {
		uc.res.Err = ErrNotGroupOwnerOrAdmin
		logrus.Error(uc.res.Err)
		return
	}

	// if `uc.req.memberId` is already admin
	if slices.IndexFunc(group.Admins, func(admin *entity.GroupAdmin) bool {
		return admin.UserId == uc.req.memberId
	}) != -1 {
		uc.res.Err = ErrIsAlreadyAdmin
		logrus.Error(uc.res.Err)
		return
	}

	// check whether `uc.req.memberId` is member
	if slices.IndexFunc(group.Members, func(mem *entity.GroupMember) bool {
		return mem.UserId == uc.req.memberId
	}) == -1 {
		uc.res.Err = ErrNotGroupMember
		logrus.Error(uc.res.Err)
		return
	}

	group.AddAdmin(uc.req.memberId, uc.req.adminId)
	uc.groupRepo.Save(group)
	uc.res.Err = nil
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
