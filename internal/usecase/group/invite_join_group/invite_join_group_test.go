package invite_join_group_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"mashu.example/internal/entity"
	entity_enums "mashu.example/internal/entity/enums"
	"mashu.example/internal/usecase/group/invite_join_group"
	"mashu.example/internal/usecase/repository/mock"
)

func setup(t *testing.T) (*mock.MockUserRepo, *mock.MockGroupRepo) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	return mock.NewMockUserRepo(mockCtrl), mock.NewMockGroupRepo(mockCtrl)
}

func TestInviteJoinGroupByOwner(t *testing.T) {
	userRepo, groupRepo := setup(t)

	inviteeId := uuid.New()
	groupId := uuid.New()
	ownerId := uuid.New()

	owner := entity.NewUser(ownerId, "owner", "Owner", "owner@email.com", false)
	invitee := entity.NewUser(inviteeId, "invitee", "Invitee", "invitee@email.com", false)
	group := entity.NewGroup(groupId, "group", owner, entity_enums.GROUP_PUBLIC)

	groupRepo.EXPECT().GetGroupById(groupId).Return(group, nil)
	userRepo.EXPECT().GetUserById(inviteeId).Return(invitee, nil)
	userRepo.EXPECT().GetUserById(ownerId).Return(owner, nil)
	groupRepo.EXPECT().Save(gomock.AssignableToTypeOf(group)).Do(
		func(arg *entity.Group) { group = arg },
	)

	req := invite_join_group.NewInviteJoinGroupUseCaseReq(invitee.ID, group.ID, owner.ID)
	res := invite_join_group.NewInviteJoinGroupUseCaseRes()
	uc := invite_join_group.NewInviteJoinGroupUseCase(userRepo, groupRepo, &req, &res)

	uc.Execute()

	assert.Nil(t, res.Err)
	assert.Len(t, group.InviteRequests, 1)
	assert.Equal(t, group.InviteRequests[0].Invitee, inviteeId)
	assert.Equal(t, group.InviteRequests[0].Inviter, ownerId)
}

func TestInviteJoinGroupByMember(t *testing.T) {
	userRepo, groupRepo := setup(t)

	inviteeId := uuid.New()
	memberId := uuid.New()
	ownerId := uuid.New()
	groupId := uuid.New()

	owner := entity.NewUser(ownerId, "owner", "owner display name", "owner@email.com", false)
	invitee := entity.NewUser(inviteeId, "invitee", "invitee display name", "invitee@email.com", false)
	member := entity.NewUser(memberId, "member", "Member", "member@email.com", false)
	group := entity.NewGroup(groupId, "group", owner, entity_enums.GROUP_PUBLIC)
	group.AddMember(memberId, uuid.Nil, ownerId)

	groupRepo.EXPECT().GetGroupById(groupId).Return(group, nil)
	userRepo.EXPECT().GetUserById(inviteeId).Return(invitee, nil)
	userRepo.EXPECT().GetUserById(memberId).Return(member, nil)
	groupRepo.EXPECT().Save(gomock.AssignableToTypeOf(group)).Do(
		func(arg *entity.Group) { group = arg },
	)

	req := invite_join_group.NewInviteJoinGroupUseCaseReq(invitee.ID, group.ID, memberId)
	res := invite_join_group.NewInviteJoinGroupUseCaseRes()
	uc := invite_join_group.NewInviteJoinGroupUseCase(userRepo, groupRepo, &req, &res)

	uc.Execute()

	assert.Nil(t, res.Err)
	assert.Len(t, group.InviteRequests, 1)
	assert.Equal(t, group.InviteRequests[0].Invitee, inviteeId)
	assert.Equal(t, group.InviteRequests[0].Inviter, memberId)
}

func TestInviteJoinGroupByStranger(t *testing.T) {
	userRepo, groupRepo := setup(t)

	inviteeId := uuid.New()
	nonMemberId := uuid.New()
	ownerId := uuid.New()
	groupId := uuid.New()

	owner := entity.NewUser(ownerId, "owner", "Owner", "owner@email.com", false)
	invitee := entity.NewUser(inviteeId, "invitee", "Invitee", "invitee@email.com", false)
	nonMember := entity.NewUser(nonMemberId, "non_member", "NonMember", "stranger@email.com", false)
	group := entity.NewGroup(groupId, "group", owner, entity_enums.GROUP_PUBLIC)

	groupRepo.EXPECT().GetGroupById(groupId).Return(group, nil)
	userRepo.EXPECT().GetUserById(inviteeId).Return(invitee, nil)
	userRepo.EXPECT().GetUserById(nonMemberId).Return(nonMember, nil)

	req := invite_join_group.NewInviteJoinGroupUseCaseReq(invitee.ID, group.ID, nonMemberId)
	res := invite_join_group.NewInviteJoinGroupUseCaseRes()
	uc := invite_join_group.NewInviteJoinGroupUseCase(userRepo, groupRepo, &req, &res)

	uc.Execute()

	assert.ErrorIs(t, res.Err, invite_join_group.ErrInviterIsNotMember)
	assert.Len(t, group.InviteRequests, 0)
	assert.Len(t, group.JoinRequests, 0)
}

func TestInviteUserWhoIsAlreadyMember(t *testing.T) {
	userRepo, groupRepo := setup(t)

	inviteeId := uuid.New()
	ownerId := uuid.New()
	groupId := uuid.New()

	owner := entity.NewUser(ownerId, "owner", "Owner", "owner@email.com", false)
	invitee := entity.NewUser(inviteeId, "invitee", "Invitee", "invitee@email.com", false)
	group := entity.NewGroup(groupId, "group", owner, entity_enums.GROUP_PUBLIC)
	group.AddMember(inviteeId, uuid.Nil, ownerId)

	userRepo.EXPECT().GetUserById(inviteeId).Return(invitee, nil)
	userRepo.EXPECT().GetUserById(ownerId).Return(owner, nil)
	groupRepo.EXPECT().GetGroupById(groupId).Return(group, nil)

	req := invite_join_group.NewInviteJoinGroupUseCaseReq(invitee.ID, group.ID, ownerId)
	res := invite_join_group.NewInviteJoinGroupUseCaseRes()
	uc := invite_join_group.NewInviteJoinGroupUseCase(userRepo, groupRepo, &req, &res)

	uc.Execute()

	assert.ErrorIs(t, res.Err, invite_join_group.ErrInviteeIsAlreadyMember)
	assert.Len(t, group.InviteRequests, 0)
}

func TestInviteUserWhoIsAlreadyInvited(t *testing.T) {
	userRepo, groupRepo := setup(t)

	inviteeId := uuid.New()
	ownerId := uuid.New()
	groupId := uuid.New()

	owner := entity.NewUser(ownerId, "owner", "Owner", "owner@email.com", false)
	invitee := entity.NewUser(inviteeId, "invitee", "Invitee", "invitee@email.com", false)
	group := entity.NewGroup(groupId, "group", owner, entity_enums.GROUP_PUBLIC)
	group.AddInviteRequest(inviteeId, ownerId)

	userRepo.EXPECT().GetUserById(inviteeId).Return(invitee, nil)
	userRepo.EXPECT().GetUserById(ownerId).Return(owner, nil)
	groupRepo.EXPECT().GetGroupById(groupId).Return(group, nil)

	req := invite_join_group.NewInviteJoinGroupUseCaseReq(invitee.ID, group.ID, ownerId)
	res := invite_join_group.NewInviteJoinGroupUseCaseRes()
	uc := invite_join_group.NewInviteJoinGroupUseCase(userRepo, groupRepo, &req, &res)

	uc.Execute()

	assert.ErrorIs(t, res.Err, invite_join_group.ErrInviteeIsAlreadyInvited)
	assert.Len(t, group.InviteRequests, 1)
}
