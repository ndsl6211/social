package join_group

import (
	"github.com/google/uuid"
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

func (uc *JoinGroupUseCase) Execute() {
	user, err := uc.userRepo.GetUserById(uc.Req.userId)
	if err != nil {
		uc.Res.Err = err
		return
	}
	group, err := uc.groupRepo.GetGroupById(uc.Req.groupId)
	if err != nil {
		uc.Res.Err = err
		return
	}

	group.AddJoinRequest(user.ID)

	uc.groupRepo.Save(group)
	uc.Res.Err = nil
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
