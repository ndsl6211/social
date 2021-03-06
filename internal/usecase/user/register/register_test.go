package register_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"mashu.example/internal/entity"
	"mashu.example/internal/usecase/repository/mock"
	"mashu.example/internal/usecase/user/register"
)

func setup(t *testing.T) *mock.MockUserRepo {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	return mock.NewMockUserRepo(mockCtrl)
}

func TestRegisterNewUser(t *testing.T) {
	userRepo := setup(t)

	var user *entity.User
	userRepo.EXPECT().Save(gomock.AssignableToTypeOf(&entity.User{})).Do(
		func(arg *entity.User) { user = arg },
	)

	req := register.NewRegisterUseCaseReq("userA", "User A", "userA@email.com")
	res := register.NewRegisterUseCaseRes()
	uc := register.NewRegisterUseCase(userRepo, req, res)

	uc.Execute()

	assert.Nil(t, res.Err)
	assert.Equal(t, "userA", user.UserName)
	assert.Equal(t, "User A", user.DisplayName)
	assert.Equal(t, "userA@email.com", user.Email)
}

func TestRegisterDuplicateUser(t *testing.T) {
	userRepo := setup(t)

	userRepo.
		EXPECT().
		Save(gomock.AssignableToTypeOf(&entity.User{})).
		Return(errors.New(""))

	req := register.NewRegisterUseCaseReq("userA", "User A", "userA@email.com")
	res := register.NewRegisterUseCaseRes()
	uc := register.NewRegisterUseCase(userRepo, req, res)

	uc.Execute()

	fmt.Println(res.Err.Error())

	assert.Error(t, res.Err)
	assert.Equal(t, "user userA already exist", res.Err.Error())
}
