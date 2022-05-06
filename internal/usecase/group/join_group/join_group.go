package join_group

import (
	"github.com/google/uuid"
	"mashu.example/internal/entity"
	entity_enums "mashu.example/internal/entity/enums"
	"mashu.example/internal/usecase"
	"mashu.example/internal/usecase/repository"
)

type JoinGroupUseCaseReq struct {
	userId  uuid.UUID
	groupId uuid.UUID
}

type JoinGroupUseCaseRes struct {
	Err error
}

type JoinGroupUseCase struct {
	userRepo  repository.UserRepo
	groupRepo repository.GroupRepo

	Req *JoinGroupUseCaseReq
	Res *JoinGroupUseCaseRes
}

func (gc *JoinGroupUseCase) Execute() {
	user, err := gc.userRepo.GetUserById(gc.Req.userId)
	group, err := gc.groupRepo.GetGroupById(gc.Req.groupId)
	if err != nil {
		gc.Res.Err = err
		return
	}

	if group.Permission == entity_enums.GROUP_PUBLIC {
		group.AddMembers(user.ID)
	} else if group.Permission == group_permission.UNPUBLIC {
		joinReq := &entity.JoinRequest{Group: group.ID, Requester: user.ID}
		group.AddJoinRequests(joinReq)
	} else {
		return
	}

	gc.groupRepo.Save(group)
	gc.Res.Err = nil
}

func NewJoinGroupUseCase(
	userRepo repository.UserRepo,
	groupRepo repository.GroupRepo,
	req *JoinGroupUseCaseReq,
	res *JoinGroupUseCaseRes,
) usecase.UseCase {
	return &JoinGroupUseCase{userRepo, groupRepo, req, res}
}

func NewJoinGroupUseCaseReq(user, group uuid.UUID) JoinGroupUseCaseReq {
	return JoinGroupUseCaseReq{userId: user, groupId: group}
}

func NewJoinGroupUseCaseRes() JoinGroupUseCaseRes {
	return JoinGroupUseCaseRes{}
}
