package handle_invite_request_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"mashu.example/internal/entity"
	entity_enums "mashu.example/internal/entity/enums"
	"mashu.example/internal/usecase/group/handle_invite_request"
	"mashu.example/internal/usecase/repository/mock"
)

func setup(t *testing.T) (*mock.MockUserRepo, *mock.MockGroupRepo) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	return mock.NewMockUserRepo(mockCtrl), mock.NewMockGroupRepo(mockCtrl)
}

func TestAcceptInvitationRequestByGroupOwner(t *testing.T) {
	userRepo, groupRepo := setup(t)

	inviteeId := uuid.New()
	inviterId := uuid.New() // the inviter may not be the member
	ownerId := uuid.New()
	groupId := uuid.New()

	invitee := entity.NewUser(inviteeId, "invitee", "Invitee", "invitee@email.com", false)
	inviter := entity.NewUser(inviterId, "inviter", "Inviter", "inviter@email.com", false)
	owner := entity.NewUser(ownerId, "owner", "Owner", "owner@email.com", false)
	group := entity.NewGroup(groupId, "my group", owner, entity_enums.GROUP_PUBLIC)
	group.AddInviteRequest(inviteeId, inviterId)

	userRepo.EXPECT().GetUserById(inviteeId).Return(invitee, nil)
	userRepo.EXPECT().GetUserById(inviterId).Return(inviter, nil)
	userRepo.EXPECT().GetUserById(ownerId).Return(owner, nil)
	groupRepo.EXPECT().GetGroupById(groupId).Return(group, nil)
	groupRepo.EXPECT().Save(gomock.AssignableToTypeOf(&entity.Group{})).Do(
		func(arg *entity.Group) { group = arg },
	)

	req := handle_invite_request.NewHandleInviteRequestUseCaseReq(
		inviteeId,
		inviterId,
		groupId,
		handle_invite_request.ACCEPT_INVITE_REQUEST,
		ownerId,
	)
	res := handle_invite_request.NewHandleInviteRequestUseCaseRes()
	uc := handle_invite_request.NewHandleInviteRequestUseCase(userRepo, groupRepo, &req, &res)
	uc.Execute()

	assert.Nil(t, res.Err)
	assert.Len(t, group.InviteRequests, 0)
	assert.Len(t, group.Members, 1)
}

func TestAcceptInvitationByGroupAdmin(t *testing.T) {
	userRepo, groupRepo := setup(t)

	inviteeId := uuid.New()
	inviterId := uuid.New() // the inviter may not be the member
	adminId := uuid.New()
	ownerId := uuid.New()
	groupId := uuid.New()

	invitee := entity.NewUser(inviteeId, "invitee", "Invitee", "invitee@email.com", false)
	inviter := entity.NewUser(inviterId, "inviter", "Inviter", "inviter@email.com", false)
	admin := entity.NewUser(adminId, "admin", "Admin", "admin@email.com", false)
	owner := entity.NewUser(ownerId, "owner", "Owner", "owner@email.com", false)
	group := entity.NewGroup(groupId, "my group", owner, entity_enums.GROUP_PUBLIC)
	group.AddAdmin(adminId, ownerId)
	group.AddInviteRequest(inviteeId, inviterId)

	userRepo.EXPECT().GetUserById(inviteeId).Return(invitee, nil)
	userRepo.EXPECT().GetUserById(inviterId).Return(inviter, nil)
	userRepo.EXPECT().GetUserById(adminId).Return(admin, nil)
	userRepo.EXPECT().GetUserById(ownerId).Return(owner, nil)
	groupRepo.EXPECT().GetGroupById(groupId).Return(group, nil)
	groupRepo.EXPECT().Save(gomock.AssignableToTypeOf(&entity.Group{})).Do(
		func(arg *entity.Group) { group = arg },
	)

	req := handle_invite_request.NewHandleInviteRequestUseCaseReq(
		inviteeId,
		inviterId,
		groupId,
		handle_invite_request.ACCEPT_INVITE_REQUEST,
		adminId,
	)
	res := handle_invite_request.NewHandleInviteRequestUseCaseRes()
	uc := handle_invite_request.NewHandleInviteRequestUseCase(userRepo, groupRepo, &req, &res)
	uc.Execute()

	assert.Nil(t, res.Err)
	assert.Len(t, group.InviteRequests, 0)
	assert.Len(t, group.Members, 1)
}

func TestRejectInvitationRequestByGroupOwner(t *testing.T) {
	userRepo, groupRepo := setup(t)

	inviteeId := uuid.New()
	inviterId := uuid.New() // the inviter may not be the member
	ownerId := uuid.New()
	groupId := uuid.New()

	invitee := entity.NewUser(inviteeId, "invitee", "Invitee", "invitee@email.com", false)
	inviter := entity.NewUser(inviterId, "inviter", "Inviter", "inviter@email.com", false)
	owner := entity.NewUser(ownerId, "owner", "Owner", "owner@email.com", false)
	group := entity.NewGroup(groupId, "my group", owner, entity_enums.GROUP_PUBLIC)
	group.AddInviteRequest(inviteeId, inviterId)

	userRepo.EXPECT().GetUserById(inviteeId).Return(invitee, nil)
	userRepo.EXPECT().GetUserById(inviterId).Return(inviter, nil)
	userRepo.EXPECT().GetUserById(ownerId).Return(owner, nil)
	groupRepo.EXPECT().GetGroupById(groupId).Return(group, nil)
	groupRepo.EXPECT().Save(gomock.AssignableToTypeOf(&entity.Group{})).Do(
		func(arg *entity.Group) { group = arg },
	)

	req := handle_invite_request.NewHandleInviteRequestUseCaseReq(
		inviteeId,
		inviterId,
		groupId,
		handle_invite_request.REJECT_INVITE_REQUEST,
		ownerId,
	)
	res := handle_invite_request.NewHandleInviteRequestUseCaseRes()
	uc := handle_invite_request.NewHandleInviteRequestUseCase(userRepo, groupRepo, &req, &res)
	uc.Execute()

	assert.Nil(t, res.Err)
	assert.Len(t, group.InviteRequests, 0)
	assert.Len(t, group.Members, 0)
}

func TestRejectInvitationByGroupAdmin(t *testing.T) {
	userRepo, groupRepo := setup(t)

	inviteeId := uuid.New()
	inviterId := uuid.New() // the inviter may not be the member
	adminId := uuid.New()
	ownerId := uuid.New()
	groupId := uuid.New()

	invitee := entity.NewUser(inviteeId, "invitee", "Invitee", "invitee@email.com", false)
	inviter := entity.NewUser(inviterId, "inviter", "Inviter", "inviter@email.com", false)
	admin := entity.NewUser(adminId, "admin", "Admin", "admin@email.com", false)
	owner := entity.NewUser(ownerId, "owner", "Owner", "owner@email.com", false)
	group := entity.NewGroup(groupId, "my group", owner, entity_enums.GROUP_PUBLIC)
	group.AddAdmin(adminId, ownerId)
	group.AddInviteRequest(inviteeId, inviterId)

	userRepo.EXPECT().GetUserById(inviteeId).Return(invitee, nil)
	userRepo.EXPECT().GetUserById(inviterId).Return(inviter, nil)
	userRepo.EXPECT().GetUserById(adminId).Return(admin, nil)
	userRepo.EXPECT().GetUserById(ownerId).Return(owner, nil)
	groupRepo.EXPECT().GetGroupById(groupId).Return(group, nil)
	groupRepo.EXPECT().Save(gomock.AssignableToTypeOf(&entity.Group{})).Do(
		func(arg *entity.Group) { group = arg },
	)

	req := handle_invite_request.NewHandleInviteRequestUseCaseReq(
		inviteeId,
		inviterId,
		groupId,
		handle_invite_request.REJECT_INVITE_REQUEST,
		adminId,
	)
	res := handle_invite_request.NewHandleInviteRequestUseCaseRes()
	uc := handle_invite_request.NewHandleInviteRequestUseCase(userRepo, groupRepo, &req, &res)
	uc.Execute()

	assert.Nil(t, res.Err)
	assert.Len(t, group.InviteRequests, 0)
	assert.Len(t, group.Members, 0)
}

func TestTryToAcceptInvitationWithoutPermission(t *testing.T) {
	userRepo, groupRepo := setup(t)

	inviteeId := uuid.New()
	inviterId := uuid.New() // the inviter may not be the member
	ownerId := uuid.New()
	groupId := uuid.New()

	invitee := entity.NewUser(inviteeId, "invitee", "Invitee", "invitee@email.com", false)
	inviter := entity.NewUser(inviterId, "inviter", "Inviter", "inviter@email.com", false)
	owner := entity.NewUser(ownerId, "owner", "Owner", "owner@email.com", false)
	group := entity.NewGroup(groupId, "my group", owner, entity_enums.GROUP_PUBLIC)
	group.AddInviteRequest(inviteeId, inviterId)

	userRepo.EXPECT().GetUserById(inviteeId).Return(invitee, nil)
	userRepo.EXPECT().GetUserById(inviterId).AnyTimes().Return(inviter, nil)
	userRepo.EXPECT().GetUserById(ownerId).Return(owner, nil)
	groupRepo.EXPECT().GetGroupById(groupId).Return(group, nil)
	groupRepo.EXPECT().Save(gomock.AssignableToTypeOf(&entity.Group{})).Do(
		func(arg *entity.Group) { group = arg },
	)

	req := handle_invite_request.NewHandleInviteRequestUseCaseReq(
		inviteeId,
		inviterId,
		groupId,
		handle_invite_request.ACCEPT_INVITE_REQUEST,
		inviterId,
	)
	res := handle_invite_request.NewHandleInviteRequestUseCaseRes()
	uc := handle_invite_request.NewHandleInviteRequestUseCase(userRepo, groupRepo, &req, &res)
	uc.Execute()

	assert.ErrorIs(t, res.Err, handle_invite_request.ErrApproverHasNoPermission)
	assert.Len(t, group.InviteRequests, 1)
	assert.Len(t, group.Members, 0)
}
