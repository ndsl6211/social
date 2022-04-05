package get_following

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"mashu.example/internal/usecase"
	"mashu.example/internal/usecase/repository"
	"mashu.example/internal/usecase/types"
)

type GetFollowingUseCaseReq struct {
	userId uuid.UUID
}

type GetFollowingUseCaseRes struct {
	Users []*types.FollowingInfo
	Err   error
}

type GetFollowingUseCase struct {
	userRepo repository.UserRepo
	req      *GetFollowingUseCaseReq
	res      *GetFollowingUseCaseRes
}

func (uc *GetFollowingUseCase) Execute() {
	user, err := uc.userRepo.GetUserById(uc.req.userId)
	if err != nil {
		errMsg := fmt.Sprintf("failed to get user (userId: %s)", uc.req.userId)
		logrus.Errorf(errMsg)
		uc.res.Err = errors.New(errMsg)
		return
	}

	followingInfos := []*types.FollowingInfo{}
	for _, followingId := range user.Followings {
		followingUser, err := uc.userRepo.GetUserById(followingId)
		if err != nil {
			errMsg := fmt.Sprintf("failed to get user's following user (userId: %s)", uc.req.userId)
			logrus.Errorf(errMsg)
			uc.res.Err = errors.New(errMsg)
			return
		}
		followingInfos = append(followingInfos, &types.FollowingInfo{
			ID:          followingUser.ID,
			UserName:    followingUser.UserName,
			DisplayName: followingUser.DisplayName,
			Email:       followingUser.Email,
			Public:      followingUser.Public,
		})
	}

	uc.res.Users = followingInfos
	uc.res.Err = nil
}

func NewGetFollowingUseCase(
	userRepo repository.UserRepo,
	req *GetFollowingUseCaseReq,
	res *GetFollowingUseCaseRes,
) usecase.UseCase {
	return &GetFollowingUseCase{userRepo, req, res}
}

func NewGetFollowingUseCaseReq(userId uuid.UUID) *GetFollowingUseCaseReq {
	return &GetFollowingUseCaseReq{userId}
}

func NewGetFollowingUseCaseRes() *GetFollowingUseCaseRes {
	return &GetFollowingUseCaseRes{}
}
