package get_following_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"mashu.example/internal/entity"
	"mashu.example/internal/usecase/repository/mock"
	"mashu.example/internal/usecase/user/get_following"
)

func setup(t *testing.T) *mock.MockUserRepo {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	return mock.NewMockUserRepo(mockCtrl)
}

func TestGetFollowingUser(t *testing.T) {
	userRepo := setup(t)

	user := entity.NewUser(uuid.New(), "mashu6211", "mashu", "mashu@email.com", true)
	following1 := entity.NewUser(uuid.New(), "following1", "following 1", "following1@email.com", true)
	following2 := entity.NewUser(uuid.New(), "following2", "following 2", "following2@email.com", true)
	following3 := entity.NewUser(uuid.New(), "following3", "following 3", "following3@email.com", true)

	user.Followings = append(user.Followings, following1.ID)
	user.Followings = append(user.Followings, following2.ID)
	user.Followings = append(user.Followings, following3.ID)

	userRepo.EXPECT().GetUserById(user.ID).Return(user, nil)
	userRepo.EXPECT().GetUserById(following1.ID).Return(following1, nil)
	userRepo.EXPECT().GetUserById(following2.ID).Return(following2, nil)
	userRepo.EXPECT().GetUserById(following3.ID).Return(following3, nil)

	req := get_following.NewGetFollowingUseCaseReq(user.ID)
	res := get_following.NewGetFollowingUseCaseRes()
	uc := get_following.NewGetFollowingUseCase(userRepo, req, res)

	uc.Execute()

	assert.Nil(t, res.Err)
	assert.Equal(t, 3, len(res.Users))
	assert.Equal(t, "following1", res.Users[0].UserName)
	assert.Equal(t, "following2", res.Users[1].UserName)
	assert.Equal(t, "following3", res.Users[2].UserName)
}
