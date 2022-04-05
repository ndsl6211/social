package get_following

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"mashu.example/internal/usecase/repository"
)

type GetFollowingUseCaseReq struct {
	userId uuid.UUID
}

type followingInfo struct {
	ID          uuid.UUID
	UserName    string
	DisplayName string
	Email       string
	Public      bool
}

type GetFollowingUseCaseRes struct {
	Users []*followingInfo
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

	followingInfos := []*followingInfo{}
	for _, followingId := range user.Followings {
		followingUser, err := uc.userRepo.GetUserById(followingId)
		if err != nil {
			errMsg := fmt.Sprintf("failed to get user's following user (userId: %s)", uc.req.userId)
			logrus.Errorf(errMsg)
			uc.res.Err = errors.New(errMsg)
			return
		}
		followingInfos = append(followingInfos, &followingInfo{
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
) *GetFollowingUseCase {
	return &GetFollowingUseCase{userRepo, req, res}
}

func NewGetFollowingUseCaseReq(userId uuid.UUID) *GetFollowingUseCaseReq {
	return &GetFollowingUseCaseReq{userId}
}

func NewGetFollowingUseCaseRes() *GetFollowingUseCaseRes {
	return &GetFollowingUseCaseRes{}
}
