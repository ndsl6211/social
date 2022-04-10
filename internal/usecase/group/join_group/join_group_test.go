package join_group_test

import (
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"mashu.example/internal/entity"
	"mashu.example/internal/entity/enums/group_permission"
	"mashu.example/internal/usecase/group/join_group"
	"mashu.example/internal/usecase/repository/mock"
	"testing"
)

func TestJoinPublicGroup(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	userId := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	groupId := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	ownerId := uuid.MustParse("11111111-1111-1111-1111-111111111111")

	owner := entity.NewUser(
		ownerId,
		"owner",
		"owner display name",
		"owner@email.com",
		false,
	)
	user := entity.NewUser(
		userId,
		"user",
		"user display name",
		"user@email.com",
		false,
	)
	group := entity.NewGroup(
		groupId,
		"group",
		owner,
		group_permission.PUBLIC,
	)

	groupRepo := mock.NewMockGroupRepo(mockCtrl)
	userRepo := mock.NewMockUserRepo(mockCtrl)
	groupRepo.EXPECT().GetGroupById(groupId).Return(group, nil)
	userRepo.EXPECT().GetUserById(userId).Return(user, nil)
	groupRepo.EXPECT().Save(gomock.AssignableToTypeOf(group)).Do(
		func(arg *entity.Group) { group = arg },
	)

	req := join_group.NewJoinGroupUseCaseReq(userId, groupId)
	res := join_group.NewJoinGroupUseCaseRes()
	gc := join_group.NewJoinGroupUseCase(userRepo, groupRepo, &req, &res)

	gc.Execute()

	assert.Equal(t, len(group.Members), 1)
	assert.Equal(t, len(group.JoinRequests), 0)
	assert.Equal(t, len(group.Admins), 0)
	assert.Equal(t, group.Owner, owner)
}

func TestJoinUnpublicGroup(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	userId := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	groupId := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	ownerId := uuid.MustParse("11111111-1111-1111-1111-111111111111")

	owner := entity.NewUser(
		ownerId,
		"owner",
		"owner display name",
		"owner@email.com",
		false,
	)
	user := entity.NewUser(
		userId,
		"user",
		"user display name",
		"user@email.com",
		false,
	)
	group := entity.NewGroup(
		groupId,
		"group",
		owner,
		group_permission.UNPUBLIC,
	)

	groupRepo := mock.NewMockGroupRepo(mockCtrl)
	userRepo := mock.NewMockUserRepo(mockCtrl)
	groupRepo.EXPECT().GetGroupById(groupId).Return(group, nil)
	userRepo.EXPECT().GetUserById(userId).Return(user, nil)
	groupRepo.EXPECT().Save(gomock.AssignableToTypeOf(group)).Do(
		func(arg *entity.Group) { group = arg },
	)

	req := join_group.NewJoinGroupUseCaseReq(userId, groupId)
	res := join_group.NewJoinGroupUseCaseRes()
	gc := join_group.NewJoinGroupUseCase(userRepo, groupRepo, &req, &res)

	gc.Execute()

	assert.Equal(t, len(group.JoinRequests), 1)
	assert.Equal(t, len(group.Members), 0)
	assert.Equal(t, len(group.Admins), 0)
	assert.Equal(t, group.Owner, owner)
}

func TestJoinPrivateGroup(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	userId := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	groupId := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	ownerId := uuid.MustParse("11111111-1111-1111-1111-111111111111")

	owner := entity.NewUser(
		ownerId,
		"owner",
		"user display name",
		"owner@email.com",
		false,
	)
	user := entity.NewUser(
		userId,
		"user",
		"user display name",
		"user@email.com",
		false,
	)
	group := entity.NewGroup(
		groupId,
		"group",
		owner,
		group_permission.PRIVATE,
	)

	groupRepo := mock.NewMockGroupRepo(mockCtrl)
	userRepo := mock.NewMockUserRepo(mockCtrl)
	groupRepo.EXPECT().GetGroupById(groupId).Return(group, nil)
	userRepo.EXPECT().GetUserById(userId).Return(user, nil)

	req := join_group.NewJoinGroupUseCaseReq(userId, groupId)
	res := join_group.NewJoinGroupUseCaseRes()
	gc := join_group.NewJoinGroupUseCase(userRepo, groupRepo, &req, &res)

	gc.Execute()

	assert.Equal(t, len(group.JoinRequests), 0)
	assert.Equal(t, len(group.Members), 0)
	assert.Equal(t, len(group.Admins), 0)
	assert.Equal(t, group.Owner, owner)
}
