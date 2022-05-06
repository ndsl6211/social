package create_group

import (
	"github.com/google/uuid"
	"mashu.example/internal/entity"
	entity_enums "mashu.example/internal/entity/enums"
	"mashu.example/internal/usecase"
	"mashu.example/internal/usecase/repository"
)

type CreateGroupUseCaseReq struct {
	name       string
	ownerId    uuid.UUID
	permission entity_enums.GroupPermission
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
	owner, err := gc.userRepo.GetUserById(gc.req.ownerId)
	if err != nil {
		gc.res.Err = err
		return
	}

	group := entity.NewGroup(uuid.New(), gc.req.name, owner, gc.req.permission)
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
	permission entity_enums.GroupPermission,
) *CreateGroupUseCaseReq {
	return &CreateGroupUseCaseReq{name, ownerId, permission}
}

func NewCreateGroupUseCaseRes() *CreateGroupUseCaseRes {
	return &CreateGroupUseCaseRes{}
}
