package follow_user

import (
	"github.com/google/uuid"
	"mashu.example/internal/model"
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

	req *FollowUserUseCaseReq
	res *FollowUserUseCaseRes
}

func (uc *FollowUserUseCase) Execute() {
	follower, err := uc.userRepo.GetUserById(uc.req.followerId)
	if err != nil {
		uc.res.Err = err
		return
	}

	followee, err := uc.userRepo.GetUserById(uc.req.followeeId)
	if err != nil {
		uc.res.Err = err
		return
	}

	if followee.Public {
		followee.AddFollower(uc.req.followerId)
		follower.AddFollowing(uc.req.followeeId)
	} else {
		followReq := &model.FollowRequest{From: uc.req.followerId, To: uc.req.followeeId}
		follower.AddFollowRequest(followReq)
		followee.AddFollowRequest(followReq)
	}

	// save system state

	if err := uc.userRepo.Save(follower); err != nil {
		uc.res.Err = err
		return
	}
	if err := uc.userRepo.Save(followee); err != nil {
		uc.res.Err = err
		return
	}

	uc.res.Err = nil
}

func NewFollowUserUseCase(
	userRepo repository.UserRepo,
	req *FollowUserUseCaseReq,
	res *FollowUserUseCaseRes,
) usecase.UseCase {
	return &FollowUserUseCase{userRepo, req, res}
}

func NewFollowUserUseCaseReq(from, to uuid.UUID) FollowUserUseCaseReq {
	return FollowUserUseCaseReq{from, to}
}

func NewFollowUserUseCaseRes() FollowUserUseCaseRes {
	return FollowUserUseCaseRes{}
}
