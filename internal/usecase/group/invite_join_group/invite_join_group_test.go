package invite_join_group_test

import (
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"golang.org/x/exp/slices"
	"mashu.example/internal/entity"
	"mashu.example/internal/entity/enums/group_permission"
	"mashu.example/internal/usecase/group/invite_join_group"
	"mashu.example/internal/usecase/repository/mock"
	"testing"
)

func TestInviteJoinGroupByOwner(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	inviteeId := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	groupId := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	ownerId := uuid.MustParse("11111111-1111-1111-1111-111111111111")

	owner := entity.NewUser(
		ownerId,
		"owner",
		"owner display name",
		"owner@email.com",
		false,
	)
	invitee := entity.NewUser(
		inviteeId,
		"invitee",
		"invitee display name",
		"invitee@email.com",
		false,
	)
	group := entity.NewGroup(
		groupId,
		"group",
		owner,
		group_permission.PRIVATE,
	)

	userRepo := mock.NewMockUserRepo(mockCtrl)
	groupRepo := mock.NewMockGroupRepo(mockCtrl)
	groupRepo.EXPECT().GetGroupById(groupId).Return(group, nil)
	userRepo.EXPECT().GetUserById(inviteeId).Return(invitee, nil)
	userRepo.EXPECT().GetUserById(ownerId).Return(owner, nil)
	groupRepo.EXPECT().Save(gomock.AssignableToTypeOf(group)).Do(
		func(arg *entity.Group) { group = arg },
	)

	req := invite_join_group.NewInviteJoinGroupUseCaseReq(invitee.ID, group.ID, owner.ID)
	res := invite_join_group.NewInviteJoinGroupUseCaseRes()
	gc := invite_join_group.NewInviteJoinGroupUseCase(userRepo, groupRepo, &req, &res)

	gc.Execute()

	assert.Equal(t, len(group.InviteRequests), 1)
	assert.Equal(t, len(group.JoinRequests), 0)
	assert.Equal(t, group.Owner, owner)
}
func TestInviteJoinGroupByAdmin(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	inviteeId := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	groupId := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	ownerId := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	adminId := uuid.MustParse("00000000-0000-0000-0000-000000000002")

	owner := entity.NewUser(
		ownerId,
		"owner",
		"owner display name",
		"owner@email.com",
		false,
	)
	invitee := entity.NewUser(
		inviteeId,
		"invitee",
		"invitee display name",
		"invitee@email.com",
		false,
	)
	admin := entity.NewUser(
		adminId,
		"admin",
		"admin display name",
		"admin@email.com",
		false,
	)
	group := entity.NewGroup(
		groupId,
		"group",
		owner,
		group_permission.PRIVATE,
	)
	group.AddAdmins(adminId)

	userRepo := mock.NewMockUserRepo(mockCtrl)
	groupRepo := mock.NewMockGroupRepo(mockCtrl)
	groupRepo.EXPECT().GetGroupById(groupId).Return(group, nil)
	userRepo.EXPECT().GetUserById(inviteeId).Return(invitee, nil)
	userRepo.EXPECT().GetUserById(adminId).Return(admin, nil)
	groupRepo.EXPECT().Save(gomock.AssignableToTypeOf(group)).Do(
		func(arg *entity.Group) { group = arg },
	)

	req := invite_join_group.NewInviteJoinGroupUseCaseReq(invitee.ID, group.ID, adminId)
	res := invite_join_group.NewInviteJoinGroupUseCaseRes()
	gc := invite_join_group.NewInviteJoinGroupUseCase(userRepo, groupRepo, &req, &res)

	gc.Execute()

	assert.Equal(t, len(group.InviteRequests), 1)
	assert.Equal(t, len(group.JoinRequests), 0)
	assert.Equal(t, group.Owner, owner)
	assert.Equal(t, slices.Contains(group.Admins, adminId), true)
}
func TestInviteJoinGroupByStranger(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	inviteeId := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	groupId := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	ownerId := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	strangerId := uuid.MustParse("00000000-0000-0000-0000-000000000002")

	owner := entity.NewUser(
		ownerId,
		"owner",
		"owner display name",
		"owner@email.com",
		false,
	)
	invitee := entity.NewUser(
		inviteeId,
		"invitee",
		"invitee display name",
		"invitee@email.com",
		false,
	)
	stranger := entity.NewUser(
		strangerId,
		"stranger",
		"stranger display name",
		"stranger@email.com",
		false,
	)
	group := entity.NewGroup(
		groupId,
		"group",
		owner,
		group_permission.PRIVATE,
	)

	userRepo := mock.NewMockUserRepo(mockCtrl)
	groupRepo := mock.NewMockGroupRepo(mockCtrl)
	groupRepo.EXPECT().GetGroupById(groupId).Return(group, nil)
	userRepo.EXPECT().GetUserById(inviteeId).Return(invitee, nil)
	userRepo.EXPECT().GetUserById(strangerId).Return(stranger, nil)

	req := invite_join_group.NewInviteJoinGroupUseCaseReq(invitee.ID, group.ID, strangerId)
	res := invite_join_group.NewInviteJoinGroupUseCaseRes()
	gc := invite_join_group.NewInviteJoinGroupUseCase(userRepo, groupRepo, &req, &res)

	gc.Execute()

	assert.Equal(t, len(group.InviteRequests), 0)
	assert.Equal(t, len(group.JoinRequests), 0)
	assert.Equal(t, group.Owner, owner)
	assert.Equal(t, slices.Contains(group.Admins, strangerId), false)
}
