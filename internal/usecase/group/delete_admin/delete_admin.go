package delete_admin

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
	ErrNotGroupOwnerOrAdmin = errors.New("only the group owner or admin can remove admin")
	ErrAdminNotFound        = errors.New("admin not found")
)

type DeleteAdminUseCaseReq struct {
	userId  uuid.UUID
	groupId uuid.UUID
	adminId uuid.UUID
}

type DeleteAdminUseCaseRes struct {
	Err error
}

type DeleteAdminUseCase struct {
	userRepo  repository.UserRepo
	groupRepo repository.GroupRepo
	req       *DeleteAdminUseCaseReq
	res       *DeleteAdminUseCaseRes
}

func (uc *DeleteAdminUseCase) Execute() {
	user, err := uc.userRepo.GetUserById(uc.req.userId)
	if err != nil {
		uc.res.Err = err
		logrus.Error(err)
		return
	}
	group, err := uc.groupRepo.GetGroupById(uc.req.groupId)
	if err != nil {
		uc.res.Err = err
		logrus.Error(err)
		return
	}
	_, err = uc.userRepo.GetUserById(uc.req.adminId)
	if err != nil {
		uc.res.Err = err
		logrus.Error(err)
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

	idx := slices.IndexFunc(group.Admins, func(admin *entity.GroupAdmin) bool {
		return admin.UserId == user.ID
	})
	if idx < 0 {
		uc.res.Err = ErrAdminNotFound
		logrus.Error(uc.res.Err)
		return
	}

	group.Admins = slices.Delete(group.Admins, idx, idx+1)

	uc.groupRepo.Save(group)
	uc.res.Err = nil
}

func NewDeleteAdminUseCase(
	userRepo repository.UserRepo,
	groupRepo repository.GroupRepo,
	req *DeleteAdminUseCaseReq,
	res *DeleteAdminUseCaseRes,
) usecase.UseCase {
	return &DeleteAdminUseCase{userRepo, groupRepo, req, res}
}

func NewDeleteAdminUseCaseReq(
	userId uuid.UUID,
	groupId uuid.UUID,
	adminId uuid.UUID,
) DeleteAdminUseCaseReq {
	return DeleteAdminUseCaseReq{userId, groupId, adminId}
}

func NewDeleteAdminUseCaseRes() DeleteAdminUseCaseRes {
	return DeleteAdminUseCaseRes{}
}
