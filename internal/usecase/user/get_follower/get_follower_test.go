package get_follower_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"mashu.example/internal/entity"
	"mashu.example/internal/usecase/repository/mock"
	"mashu.example/internal/usecase/user/get_follower"
)

func setup(t *testing.T) *mock.MockUserRepo {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	return mock.NewMockUserRepo(mockCtrl)
}

func TestGetFollower(t *testing.T) {
	userRepo := setup(t)

	user := entity.NewUser(uuid.New(), "mashu6211", "mashu", "mashu@email.com", true)
	follower1 := entity.NewUser(uuid.New(), "follower1", "follower 1", "follower1@email.com", true)
	follower2 := entity.NewUser(uuid.New(), "follower2", "follower 2", "follower2@email.com", true)
	follower3 := entity.NewUser(uuid.New(), "follower3", "follower 3", "follower3@email.com", true)

	user.Followers = append(user.Followers, follower1.ID)
	user.Followers = append(user.Followers, follower2.ID)
	user.Followers = append(user.Followers, follower3.ID)

	userRepo.EXPECT().GetUserById(user.ID).Return(user, nil)
	userRepo.EXPECT().GetUserById(follower1.ID).Return(follower1, nil)
	userRepo.EXPECT().GetUserById(follower2.ID).Return(follower2, nil)
	userRepo.EXPECT().GetUserById(follower3.ID).Return(follower3, nil)

	req := get_follower.NewGetFollowerUseCaseReq(user.ID)
	res := get_follower.NewGetFollowerUseCaseRes()
	uc := get_follower.NewGetFollowerUsecase(userRepo, req, res)

	uc.Execute()

	assert.Nil(t, res.Err)
	assert.Equal(t, 3, len(res.Users))
	assert.Equal(t, follower1.UserName, res.Users[0].UserName)
	assert.Equal(t, follower2.UserName, res.Users[1].UserName)
	assert.Equal(t, follower3.UserName, res.Users[2].UserName)
}
