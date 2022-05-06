package handle_follow_request_test

import (
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"mashu.example/internal/entity"
	"mashu.example/internal/usecase/repository/mock"
	"mashu.example/internal/usecase/user/handle_follow_request"
)

func TestAcceptFollowRequest(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	followerId := uuid.New()
	followeeId := uuid.New()

	followee := &entity.User{
		ID:          followeeId,
		UserName:    "followee",
		DisplayName: "followee display name",
		Email:       "followee@email.com",
		Public:      true,
		Followers:   []uuid.UUID{},
		Followings:  []uuid.UUID{},
		FollowRequests: []*entity.FollowRequest{
			{
				From: followerId,
				To:   followeeId,
			},
		},
	}

	follower := &entity.User{
		ID:          followerId,
		UserName:    "follower",
		DisplayName: "follower display name",
		Email:       "follower@email.com",
		Public:      false,
		Followers:   []uuid.UUID{},
		Followings:  []uuid.UUID{},
		FollowRequests: []*entity.FollowRequest{
			{
				From: followerId,
				To:   followeeId,
			},
		},
	}

	repo := mock.NewMockUserRepo(mockCtrl)
	repo.EXPECT().GetUserById(followerId).Return(follower, nil)
	repo.EXPECT().GetUserById(followeeId).Return(followee, nil)

	// should save the followee first
	repo.EXPECT().Save(gomock.AssignableToTypeOf(&entity.User{})).Do(
		func(arg *entity.User) { followee = arg },
	)

	// and then save follower
	repo.EXPECT().Save(gomock.AssignableToTypeOf(&entity.User{})).Do(
		func(arg *entity.User) { follower = arg },
	)

	req := handle_follow_request.NewHandleFollowRequestUseCaseReq(
		followeeId,
		followerId,
		handle_follow_request.ACCEPT_FOLLOW_REQUEST,
	)
	res := handle_follow_request.NewHandleFollowRequestUsecaseRes()
	uc := handle_follow_request.NewHandleFollowRequestUseCase(repo, &req, &res)

	uc.Execute()

	if res.Err != nil {
		t.Errorf("failed to execute usecase")
	}

	assert.Equal(t, len(follower.Followings), 1)
	assert.Equal(t, len(follower.FollowRequests), 0)

	assert.Equal(t, len(followee.Followers), 1)
	assert.Equal(t, len(followee.FollowRequests), 0)
}

func TestRejectFollowRequest(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	followerId := uuid.New()
	followeeId := uuid.New()

	followee := &entity.User{
		ID:          followeeId,
		UserName:    "followee",
		DisplayName: "followee display name",
		Email:       "followee@email.com",
		Public:      true,
		Followers:   []uuid.UUID{},
		Followings:  []uuid.UUID{},
		FollowRequests: []*entity.FollowRequest{
			{
				From: followerId,
				To:   followeeId,
			},
		},
	}

	follower := &entity.User{
		ID:          followerId,
		UserName:    "follower",
		DisplayName: "follower display name",
		Email:       "follower@email.com",
		Public:      false,
		Followers:   []uuid.UUID{},
		Followings:  []uuid.UUID{},
		FollowRequests: []*entity.FollowRequest{
			{
				From: followerId,
				To:   followeeId,
			},
		},
	}

	repo := mock.NewMockUserRepo(mockCtrl)
	repo.EXPECT().GetUserById(followerId).Return(follower, nil)
	repo.EXPECT().GetUserById(followeeId).Return(followee, nil)

	// should save the followee first
	repo.EXPECT().Save(gomock.AssignableToTypeOf(&entity.User{})).Do(
		func(arg *entity.User) { followee = arg },
	)

	// and then save follower
	repo.EXPECT().Save(gomock.AssignableToTypeOf(&entity.User{})).Do(
		func(arg *entity.User) { follower = arg },
	)

	req := handle_follow_request.NewHandleFollowRequestUseCaseReq(
		followeeId,
		followerId,
		handle_follow_request.REJECT_FOLLOW_REQUEST,
	)
	res := handle_follow_request.NewHandleFollowRequestUsecaseRes()
	uc := handle_follow_request.NewHandleFollowRequestUseCase(repo, &req, &res)

	uc.Execute()

	if res.Err != nil {
		t.Errorf("failed to execute usecase")
	}

	assert.Equal(t, len(follower.Followings), 0)
	assert.Equal(t, len(follower.FollowRequests), 0)

	assert.Equal(t, len(followee.Followers), 0)
	assert.Equal(t, len(followee.FollowRequests), 0)
}
