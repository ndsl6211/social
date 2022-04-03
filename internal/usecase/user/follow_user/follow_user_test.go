package follow_user_test

import (
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"mashu.example/internal/entity"
	"mashu.example/internal/usecase/repository/mock"
	"mashu.example/internal/usecase/user/follow_user"
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
	repo.EXPECT().GetUserById(followerId).Return(follower, nil)
	repo.EXPECT().GetUserById(followeeId).Return(followee, nil)

	// should call with follower first
	repo.EXPECT().Save(gomock.AssignableToTypeOf(&entity.User{})).Do(
		func(arg *entity.User) { follower = arg },
	)

	// and then call with followee
	repo.EXPECT().Save(gomock.AssignableToTypeOf(&entity.User{})).Do(
		func(arg *entity.User) { followee = arg },
	)

	req := follow_user.NewFollowUserUseCaseReq(followerId.String(), followeeId.String())
	res := follow_user.NewFollowUserUseCaseRes()
	uc := follow_user.NewFollowUserUseCase(repo, &req, &res)

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
	repo.EXPECT().GetUserById(followerId).Return(follower, nil)
	repo.EXPECT().GetUserById(followeeId).Return(followee, nil)
	repo.EXPECT().Save(gomock.AssignableToTypeOf(follower)).Do(
		func(arg *entity.User) { follower = arg },
	)
	repo.EXPECT().Save(gomock.AssignableToTypeOf(followee)).Do(
		func(arg *entity.User) { followee = arg },
	)

	req := follow_user.NewFollowUserUseCaseReq(followerId.String(), followeeId.String())
	res := follow_user.NewFollowUserUseCaseRes()
	uc := follow_user.NewFollowUserUseCase(repo, &req, &res)

	uc.Execute()

	assert.Equal(t, len(follower.Followings), 1)
	assert.Equal(t, len(follower.Followers), 0)
	assert.Equal(t, len(follower.FollowRequests), 0)

	assert.Equal(t, len(followee.Followings), 0)
	assert.Equal(t, len(followee.Followers), 1)
	assert.Equal(t, len(followee.FollowRequests), 0)
}
