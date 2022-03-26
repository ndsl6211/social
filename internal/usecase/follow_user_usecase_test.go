package usecase_test

import (
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"mashu.example/internal/entity"
	"mashu.example/internal/usecase"
	"mashu.example/internal/usecase/repository/mock"
)

func TestFollowPrivateUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	followerId := uuid.MustParse("00000000-0000-0000-0000-000000000000")
	followeeId := uuid.MustParse("11111111-1111-1111-1111-111111111111")

	follower := entity.NewUser(
		followerId,
		"follower",
		"follower display name",
		"follower@email.com",
		false,
	)

	followee := entity.NewUser(
		followeeId,
		"followee",
		"followee display name",
		"folowee@email.com",
		false,
	)

	repo := mock.NewMockUserRepo(mockCtrl)
	repo.EXPECT().GetUserById(followerId.String()).Return(follower, nil)
	repo.EXPECT().GetUserById(followeeId.String()).Return(followee, nil)
	repo.EXPECT().Save(gomock.AssignableToTypeOf(follower)).Do(
		func(arg *entity.User) { follower = arg },
	)
	repo.EXPECT().Save(gomock.AssignableToTypeOf(followee)).Do(
		func(arg *entity.User) { followee = arg },
	)

	req := usecase.NewFollowUserUseCaseReq(followerId.String(), followeeId.String())
	res := usecase.NewFollowUserUseCaseRes()
	uc := usecase.NewFollowUserUseCase(repo, &req, &res)

	uc.Execute()

	assert.Equal(t, len(follower.FollowRequests), 1)
	assert.Equal(t, len(followee.FollowRequests), 1)
}

func TestFollowPublicUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	followerId := uuid.MustParse("00000000-0000-0000-0000-000000000000")
	followeeId := uuid.MustParse("11111111-1111-1111-1111-111111111111")

	follower := entity.NewUser(
		followerId,
		"follower",
		"follower display name",
		"follower@email.com",
		false,
	)

	followee := entity.NewUser(
		followeeId,
		"followee",
		"followee display name",
		"folowee@email.com",
		true,
	)

	repo := mock.NewMockUserRepo(mockCtrl)
	repo.EXPECT().GetUserById(followerId.String()).Return(follower, nil)
	repo.EXPECT().GetUserById(followeeId.String()).Return(followee, nil)
	repo.EXPECT().Save(gomock.AssignableToTypeOf(follower)).Do(
		func(arg *entity.User) { follower = arg },
	)
	repo.EXPECT().Save(gomock.AssignableToTypeOf(followee)).Do(
		func(arg *entity.User) { followee = arg },
	)

	req := usecase.NewFollowUserUseCaseReq(followerId.String(), followeeId.String())
	res := usecase.NewFollowUserUseCaseRes()
	uc := usecase.NewFollowUserUseCase(repo, &req, &res)

	uc.Execute()

	assert.Equal(t, len(follower.Followings), 1)
	assert.Equal(t, len(follower.Followers), 0)
	assert.Equal(t, len(follower.FollowRequests), 0)

	assert.Equal(t, len(followee.Followings), 0)
	assert.Equal(t, len(followee.Followers), 1)
	assert.Equal(t, len(followee.FollowRequests), 0)
}
