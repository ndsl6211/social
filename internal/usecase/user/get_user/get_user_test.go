package get_user_test

import (
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"mashu.example/internal/entity"
	"mashu.example/internal/usecase/repository/mock"
	"mashu.example/internal/usecase/user/get_user"
)

func TestGetUserUseCase(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	userId := uuid.MustParse("00000000-0000-0000-0000-000000000000")
	user := &entity.User{
		ID:          userId,
		UserName:    "username",
		DisplayName: "displayName",
		Email:       "test@email.com",
		Public:      true,
		Followers: []uuid.UUID{
			uuid.MustParse("11111111-1111-1111-1111-111111111111"),
			uuid.MustParse("22222222-2222-2222-2222-222222222222"),
		},
		Followings: []uuid.UUID{
			uuid.MustParse("11111111-1111-1111-1111-111111111111"),
			uuid.MustParse("33333333-3333-3333-3333-333333333333"),
		},
		FollowRequests: []*entity.FollowRequest{
			{
				From: uuid.MustParse("44444444-4444-4444-4444-444444444444"),
				To:   uuid.MustParse("11111111-1111-1111-1111-111111111111"),
			},
			{
				From: uuid.MustParse("55555555-5555-5555-5555-555555555555"),
				To:   uuid.MustParse("11111111-1111-1111-1111-111111111111"),
			},
			{
				From: uuid.MustParse("11111111-1111-1111-1111-111111111111"),
				To:   uuid.MustParse("55555555-5555-5555-5555-555555555555"),
			},
		},
	}
	repo := mock.NewMockUserRepo(mockCtrl)
	repo.EXPECT().GetUserById(userId).Return(user, nil)

	req := get_user.NewGetUserUseCaseReq(userId)
	res := get_user.NewGetUserUseCaseRes()
	usecase := get_user.NewGetUserUseCase(repo, &req, &res)
	usecase.Execute()

	if res.Err != nil {
		t.Errorf("test failed! %s", res.Err.Error())
	}

	assert.Equal(t, res.ID, user.ID)
	assert.Equal(t, res.UserName, user.UserName)
	assert.Equal(t, res.DisplayName, user.DisplayName)
	assert.Equal(t, res.Email, user.Email)
	assert.Equal(t, res.Public, user.Public)
}
