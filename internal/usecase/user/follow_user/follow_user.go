package follow_user

import (
	"github.com/google/uuid"
	"mashu.example/internal/entity"
	"mashu.example/internal/usecase"
	"mashu.example/internal/usecase/repository"
)

type FollowUserUseCaseReq struct {
	followerId string
	followeeId string
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
	followerId := uuid.MustParse(uc.Req.followerId)
	followeeId := uuid.MustParse(uc.Req.followeeId)

	follower, err := uc.userRepo.GetUserById(followerId)
	followee, err := uc.userRepo.GetUserById(followeeId)
	if err != nil {
		uc.Res.Err = err
		return
	}

	if followee.Public {
		followee.AddFollower(followerId)
		follower.AddFollowing(followeeId)
	} else {
		followReq := &entity.FollowRequest{From: followerId, To: followeeId}
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

func NewFollowUserUseCaseReq(from, to string) FollowUserUseCaseReq {
	return FollowUserUseCaseReq{from, to}
}

func NewFollowUserUseCaseRes() FollowUserUseCaseRes {
	return FollowUserUseCaseRes{}
}
