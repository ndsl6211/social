package delete_admin

import (
	"errors"
	"github.com/google/uuid"
	"golang.org/x/exp/slices"
	"mashu.example/internal/usecase"
	"mashu.example/internal/usecase/repository"
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

func (gc *DeleteAdminUseCase) Execute() {
	user, err := gc.userRepo.GetUserById(gc.req.userId)
	group, err := gc.groupRepo.GetGroupById(gc.req.groupId)
	admin, err := gc.userRepo.GetUserById(gc.req.adminId)
	if err != nil {
		gc.res.Err = err
		return
	}
	if !slices.Contains(group.Admins, admin.ID) && admin == group.Owner {
		errMsg := "permission denied"
		gc.res.Err = errors.New(errMsg)
		return
	}

	idx := slices.IndexFunc(group.Admins, func(req uuid.UUID) bool {
		return req == user.ID
	})

	if idx < 0 {
		errMsg := "user not found"
		gc.res.Err = errors.New(errMsg)
		return
	}

	group.Admins = slices.Delete(group.Admins, idx, idx+1)

	gc.groupRepo.Save(group)
	gc.res.Err = nil
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
