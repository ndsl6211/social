package get_follower

import (
	"github.com/google/uuid"
	"mashu.example/internal/usecase"
	"mashu.example/internal/usecase/dto"
	"mashu.example/internal/usecase/repository"
)

type GetFollowerUseCaseReq struct {
	userId uuid.UUID
}

type GetFollowerUseCaseRes struct {
	Users []*dto.FollowingInfo
	Err   error
}

type GetFollowerUseCase struct {
	userRepo repository.UserRepo
	req      *GetFollowerUseCaseReq
	res      *GetFollowerUseCaseRes
}

func (uc *GetFollowerUseCase) Execute() {
	user, err := uc.userRepo.GetUserById(uc.req.userId)
	if err != nil {
		uc.res.Err = err
		return
	}

	followerInfos := []*dto.FollowingInfo{}
	for _, followerId := range user.Followers {
		followerUser, err := uc.userRepo.GetUserById(followerId)
		if err != nil {
			uc.res.Err = err
			return
		}

		followerInfos = append(followerInfos, &dto.FollowingInfo{
			ID:       followerUser.ID,
			UserName: followerUser.UserName,
			Email:    followerUser.Email,
			Public:   followerUser.Public,
		})
	}

	uc.res.Users = followerInfos
	uc.res.Err = nil
}

func NewGetFollowerUsecase(
	userRepo repository.UserRepo,
	req *GetFollowerUseCaseReq,
	res *GetFollowerUseCaseRes,
) usecase.UseCase {
	return &GetFollowerUseCase{userRepo, req, res}
}

func NewGetFollowerUseCaseReq(userId uuid.UUID) *GetFollowerUseCaseReq {
	return &GetFollowerUseCaseReq{userId}
}

func NewGetFollowerUseCaseRes() *GetFollowerUseCaseRes {
	return &GetFollowerUseCaseRes{}
}
