package create_group_test

import (
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"mashu.example/internal/entity"
	entity_enums "mashu.example/internal/entity/enums"
	"mashu.example/internal/usecase/group/create_group"
	"mashu.example/internal/usecase/repository/mock"
)

func setup(t *testing.T) (*mock.MockGroupRepo, *mock.MockUserRepo) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	return mock.NewMockGroupRepo(mockCtrl), mock.NewMockUserRepo(mockCtrl)
}

func TestCreateGroup(t *testing.T) {
	groupRepo, userRepo := setup(t)

	ownerId := uuid.MustParse("00000000-0000-0000-0000-000000000000")
	owner := entity.NewUser(
		ownerId,
		"owner",
		"owner display name",
		"owner@email.com",
		false,
	)
	userRepo.EXPECT().GetUserById(ownerId).Return(owner, nil)
	var resultGroup *entity.Group
	groupRepo.EXPECT().Save(gomock.AssignableToTypeOf(&entity.Group{})).Do(
		func(arg *entity.Group) { resultGroup = arg },
	)

	req := create_group.NewCreateGroupUseCaseReq(
		"First Group",
		ownerId,
		entity_enums.GROUP_PUBLIC,
	)
	res := create_group.NewCreateGroupUseCaseRes()
	gc := create_group.NewCreateGroupUseCase(groupRepo, userRepo, req, res)

	gc.Execute()

	if res.Err != nil {
		t.Errorf("failed to execute usecase")
	}

	assert.Equal(t, resultGroup.Name, "First Group")
	assert.Equal(t, resultGroup.Owner, owner)
	assert.Equal(t, resultGroup.Permission, entity_enums.GROUP_PUBLIC)
}
