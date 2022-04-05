package get_follower

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"mashu.example/internal/usecase"
	"mashu.example/internal/usecase/repository"
	"mashu.example/internal/usecase/types"
)

type GetFollowerUseCaseReq struct {
	userId uuid.UUID
}

type GetFollowerUseCaseRes struct {
	Users []*types.FollowingInfo
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
		errMsg := fmt.Sprintf("failed to get user (userId: %s)", uc.req.userId)
		logrus.Errorf(errMsg)
		uc.res.Err = errors.New(errMsg)
		return
	}

	followerInfos := []*types.FollowingInfo{}
	for _, followerId := range user.Followers {
		followerUser, err := uc.userRepo.GetUserById(followerId)
		if err != nil {
			errMsg := fmt.Sprintf("failed to get user's following user (userId: %s)", uc.req.userId)
			logrus.Errorf(errMsg)
			uc.res.Err = errors.New(errMsg)
			return
		}
		followerInfos = append(followerInfos, &types.FollowingInfo{
			ID:          followerUser.ID,
			UserName:    followerUser.UserName,
			DisplayName: followerUser.DisplayName,
			Email:       followerUser.Email,
			Public:      followerUser.Public,
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
