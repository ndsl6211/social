package follow_user

import (
	"github.com/google/uuid"
	"mashu.example/internal/entity"
	"mashu.example/internal/usecase"
	"mashu.example/internal/usecase/repository"
)

type FollowUserUseCaseReq struct {
	followerId uuid.UUID
	followeeId uuid.UUID
}

type FollowUserUseCaseRes struct {
	Err error
}

type FollowUserUseCase struct {
	userRepo repository.UserRepo

	Req *FollowUserUseCaseReq
	Res *FollowUserUseCaseRes
}

func (uc *FollowUserUseCase) Execute() {
	follower, err := uc.userRepo.GetUserById(uc.Req.followerId)
	followee, err := uc.userRepo.GetUserById(uc.Req.followeeId)
	if err != nil {
		uc.Res.Err = err
		return
	}

	if followee.Public {
		followee.AddFollower(uc.Req.followerId)
		follower.AddFollowing(uc.Req.followeeId)
	} else {
		followReq := &entity.FollowRequest{From: uc.Req.followerId, To: uc.Req.followeeId}
		follower.AddFollowRequest(followReq)
		followee.AddFollowRequest(followReq)
	}

	// save system state
	uc.userRepo.Save(follower)
	uc.userRepo.Save(followee)

	uc.Res.Err = nil
}

func NewFollowUserUseCase(
	userRepo repository.UserRepo,
	req *FollowUserUseCaseReq,
	res *FollowUserUseCaseRes,
) usecase.UseCase {
	return &FollowUserUseCase{userRepo: userRepo, Req: req, Res: res}
}

func NewFollowUserUseCaseReq(from, to uuid.UUID) FollowUserUseCaseReq {
	return FollowUserUseCaseReq{from, to}
}

func NewFollowUserUseCaseRes() FollowUserUseCaseRes {
	return FollowUserUseCaseRes{}
}
