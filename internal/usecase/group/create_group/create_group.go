package create_group

import (
	"github.com/google/uuid"
	"mashu.example/internal/entity"
	"mashu.example/internal/entity/enums/group_permission"
	"mashu.example/internal/usecase"
	"mashu.example/internal/usecase/repository"
)

type CreateGroupUseCaseReq struct {
	Name       string
	OwnerId    uuid.UUID
	Permission group_permission.GroupPermission
}

type CreateGroupUseCaseRes struct {
	Err error
}

type CreateGroupUseCase struct {
	groupRepo repository.GroupRepo
	userRepo  repository.UserRepo

	req *CreateGroupUseCaseReq
	res *CreateGroupUseCaseRes
}

func (gc *CreateGroupUseCase) Execute() {
	owner, err := gc.userRepo.GetUserById(gc.req.OwnerId)
	if err != nil {
		gc.res.Err = err
		return
	}

	group := entity.NewGroup(uuid.New(), gc.req.Name, owner, gc.req.Permission)
	gc.groupRepo.Save(group)

	gc.res.Err = nil
}

func NewCreateGroupUseCase(
	groupRepo repository.GroupRepo,
	userRepo repository.UserRepo,
	req *CreateGroupUseCaseReq,
	res *CreateGroupUseCaseRes,
) usecase.UseCase {
	return &CreateGroupUseCase{groupRepo, userRepo, req, res}
}

func NewCreateGroupUseCaseReq(
	name string,
	ownerId uuid.UUID,
	permission group_permission.GroupPermission,
) *CreateGroupUseCaseReq {
	return &CreateGroupUseCaseReq{name, ownerId, permission}
}

func NewCreateGroupUseCaseRes() *CreateGroupUseCaseRes {
	return &CreateGroupUseCaseRes{}
}
