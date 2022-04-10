package create_group_test

import (
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"mashu.example/internal/entity"
	"mashu.example/internal/entity/enums/group_permission"
	"mashu.example/internal/usecase/group/create_group"
	"mashu.example/internal/usecase/repository/mock"
	"testing"
)

func setup(t *testing.T) (*mock.MockGroupRepo, *mock.MockUserRepo) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	return mock.NewMockGroupRepo(mockCtrl), mock.NewMockUserRepo(mockCtrl)
}

func TestCreatePublicGroup(t *testing.T) {
	groupRepo, userRepo := setup(t)

	ownerId := uuid.MustParse("00000000-0000-0000-0000-000000000001")
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
		group_permission.PUBLIC,
	)
	res := create_group.NewCreateGroupUseCaseRes()
	gc := create_group.NewCreateGroupUseCase(groupRepo, userRepo, req, res)

	gc.Execute()

	if res.Err != nil {
		t.Errorf("failed to execute usecase")
	}

	assert.Equal(t, resultGroup.Name, "First Group")
	assert.Equal(t, resultGroup.Owner, owner)
	assert.Equal(t, resultGroup.Permission, group_permission.PUBLIC)
}

func TestCreateUnpublicGroup(t *testing.T) {
	groupRepo, userRepo := setup(t)

	ownerId := uuid.MustParse("00000000-0000-0000-0000-000000000002")
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
		"Second Group",
		ownerId,
		group_permission.UNPUBLIC,
	)
	res := create_group.NewCreateGroupUseCaseRes()
	gc := create_group.NewCreateGroupUseCase(groupRepo, userRepo, req, res)

	gc.Execute()

	if res.Err != nil {
		t.Errorf("failed to execute usecase")
	}

	assert.Equal(t, resultGroup.Name, "Second Group")
	assert.Equal(t, resultGroup.Owner, owner)
	assert.Equal(t, resultGroup.Permission, group_permission.UNPUBLIC)
}

func TestCreatePrivateGroup(t *testing.T) {
	groupRepo, userRepo := setup(t)

	ownerId := uuid.MustParse("00000000-0000-0000-0000-000000000003")
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
		"Third Group",
		ownerId,
		group_permission.PRIVATE,
	)
	res := create_group.NewCreateGroupUseCaseRes()
	gc := create_group.NewCreateGroupUseCase(groupRepo, userRepo, req, res)

	gc.Execute()

	if res.Err != nil {
		t.Errorf("failed to execute usecase")
	}

	assert.Equal(t, resultGroup.Name, "Third Group")
	assert.Equal(t, resultGroup.Owner, owner)
	assert.Equal(t, resultGroup.Permission, group_permission.PRIVATE)
}
