package login_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"mashu.example/internal/entity"
	"mashu.example/internal/usecase/repository/mock"
	"mashu.example/internal/usecase/user/login"
)

func setup(t *testing.T) *mock.MockUserRepo {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	return mock.NewMockUserRepo(mockCtrl)
}

func TestLoginAsValidUser(t *testing.T) {
	userRepo := setup(t)

	user := entity.NewUser(uuid.New(), "mashu6211", "Mashu", "mashu@email.com", false)

	userRepo.EXPECT().GetUserByUserName(user.UserName).Return(user, nil)

	req := login.NewLoginUseCaseReq("mashu6211")
	res := login.NewLoginUseCaseRes()
	uc := login.NewLoginUseCase(userRepo, req, res)

	uc.Execute()

	assert.Nil(t, res.Err)
	assert.NotEmpty(t, res.AccessToken)
}

func TestLoginAsNonExistUser(t *testing.T) {
	userRepo := setup(t)
	userRepo.EXPECT().GetUserByUserName("mashu6211").Return(nil, gorm.ErrRecordNotFound)

	req := login.NewLoginUseCaseReq("mashu6211")
	res := login.NewLoginUseCaseRes()
	uc := login.NewLoginUseCase(userRepo, req, res)

	uc.Execute()

	assert.Error(t, res.Err)
	assert.Empty(t, res.AccessToken)
}
