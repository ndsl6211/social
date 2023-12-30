package handle_follow_request

import (
	"errors"

	"github.com/google/uuid"
	"golang.org/x/exp/slices"
	"mashu.example/internal/model"
	"mashu.example/internal/usecase"
	"mashu.example/internal/usecase/repository"
)

type HandleFollowRequestAction string

const (
	ACCEPT_FOLLOW_REQUEST HandleFollowRequestAction = "ACCEPT"
	REJECT_FOLLOW_REQUEST HandleFollowRequestAction = "REJECT"
)

type HandleFollowRequestUseCaseReq struct {
	userId     uuid.UUID
	followerId uuid.UUID
	action     HandleFollowRequestAction
}

type HandleFollowRequestUseCaseRes struct {
	Err error
}

type HandleFollowRequestUseCase struct {
	userRepo repository.UserRepo
	req      *HandleFollowRequestUseCaseReq
	res      *HandleFollowRequestUseCaseRes
}

func (uc *HandleFollowRequestUseCase) Execute() {
	user, err := uc.userRepo.GetUserById(uc.req.userId)
	if err != nil {
		uc.res.Err = err
		return
	}
	follower, err := uc.userRepo.GetUserById(uc.req.followerId)
	if err != nil {
		uc.res.Err = err
		return
	}

	var targetFollowReqIdx int
	var targetFollowReq *model.FollowRequest = nil
	for idx, followReq := range user.FollowRequests {
		if followReq.From == uc.req.followerId && followReq.To == uc.req.userId {
			targetFollowReq = followReq
			targetFollowReqIdx = idx
			break
		}
	}
	if targetFollowReq == nil {
		uc.res.Err = errors.New("follow request not found in followee")
		return
	}
	user.FollowRequests = slices.Delete(user.FollowRequests, targetFollowReqIdx, targetFollowReqIdx+1)

	targetFollowReq = nil
	for idx, followReq := range follower.FollowRequests {
		if followReq.From == uc.req.followerId && followReq.To == uc.req.userId {
			targetFollowReq = followReq
			targetFollowReqIdx = idx
			break
		}
	}
	if targetFollowReq == nil {
		uc.res.Err = errors.New("follow request not found in follower")
		return
	}
	follower.FollowRequests = slices.Delete(follower.FollowRequests, targetFollowReqIdx, targetFollowReqIdx+1)

	if uc.req.action == ACCEPT_FOLLOW_REQUEST {
		user.AddFollower(uc.req.followerId)
		follower.AddFollowing(uc.req.userId)
	}

	uc.userRepo.Save(user)
	uc.userRepo.Save(follower)

	uc.res.Err = nil
}

func NewHandleFollowRequestUseCase(
	userRepo repository.UserRepo,
	req *HandleFollowRequestUseCaseReq,
	res *HandleFollowRequestUseCaseRes,
) usecase.UseCase {
	return &HandleFollowRequestUseCase{userRepo, req, res}
}

func NewHandleFollowRequestUseCaseReq(
	userId uuid.UUID,
	followerId uuid.UUID,
	action HandleFollowRequestAction,
) HandleFollowRequestUseCaseReq {
	return HandleFollowRequestUseCaseReq{userId, followerId, action}
}

func NewHandleFollowRequestUsecaseRes() HandleFollowRequestUseCaseRes {
	return HandleFollowRequestUseCaseRes{}
}
