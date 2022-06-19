package handle_join_request_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"mashu.example/internal/entity"
	entity_enums "mashu.example/internal/entity/enums"
	"mashu.example/internal/usecase/group/handle_join_request"
	"mashu.example/internal/usecase/repository/mock"
)

func setup(t *testing.T) (*mock.MockUserRepo, *mock.MockGroupRepo) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	return mock.NewMockUserRepo(mockCtrl), mock.NewMockGroupRepo(mockCtrl)
}

func TestAcceptJoinRequestByGroupOwner(t *testing.T) {
	userRepo, groupRepo := setup(t)

	requesterId := uuid.New()
	ownerId := uuid.New()
	groupId := uuid.New()

	requester := entity.NewUser(requesterId, "requester", "Requester", "requester@email.com", false)
	owner := entity.NewUser(ownerId, "owner", "Owner", "owner@email.com", false)
	group := entity.NewGroup(groupId, "my group", owner, entity_enums.GROUP_PUBLIC)
	group.AddJoinRequest(requesterId)

	userRepo.EXPECT().GetUserById(requesterId).Return(requester, nil)
	userRepo.EXPECT().GetUserById(ownerId).Return(owner, nil)
	groupRepo.EXPECT().GetGroupById(groupId).Return(group, nil)
	groupRepo.EXPECT().Save(gomock.AssignableToTypeOf(&entity.Group{})).Do(
		func(arg *entity.Group) { group = arg },
	)

	req := handle_join_request.NewHandleJoinRequestUseCaseReq(
		requesterId,
		groupId,
		handle_join_request.ACCEPT_JOIN_REQUEST,
		ownerId,
	)
	res := handle_join_request.NewHandleJoinRequestUseCaseRes()
	uc := handle_join_request.NewHandleJoinRequestUseCase(userRepo, groupRepo, &req, &res)
	uc.Execute()

	assert.Nil(t, res.Err)
	assert.Len(t, group.JoinRequests, 0)
	assert.Len(t, group.Members, 1)
	assert.Equal(t, group.Members[0].UserId, requesterId)
	assert.Equal(t, group.Members[0].ApprovedBy, ownerId)
	assert.Equal(t, group.Members[0].InvitedBy, uuid.Nil)
}

func TestAcceptJoinRequestByGroupAdmin(t *testing.T) {
	userRepo, groupRepo := setup(t)

	requesterId := uuid.New()
	approverId := uuid.New()
	ownerId := uuid.New()
	groupId := uuid.New()

	requester := entity.NewUser(requesterId, "requester", "Requester", "requester@email.com", false)
	approver := entity.NewUser(approverId, "approver", "Approver", "approver@email.com", false)
	owner := entity.NewUser(ownerId, "owner", "Owner", "owner@email.com", false)
	group := entity.NewGroup(groupId, "my group", owner, entity_enums.GROUP_PUBLIC)
	group.AddAdmin(approverId, ownerId)
	group.AddJoinRequest(requesterId)

	userRepo.EXPECT().GetUserById(requesterId).Return(requester, nil)
	userRepo.EXPECT().GetUserById(approverId).Return(approver, nil)
	userRepo.EXPECT().GetUserById(ownerId).Return(owner, nil)
	groupRepo.EXPECT().GetGroupById(groupId).Return(group, nil)
	groupRepo.EXPECT().Save(gomock.AssignableToTypeOf(&entity.Group{})).Do(
		func(arg *entity.Group) { group = arg },
	)

	req := handle_join_request.NewHandleJoinRequestUseCaseReq(
		requesterId,
		groupId,
		handle_join_request.ACCEPT_JOIN_REQUEST,
		approverId,
	)
	res := handle_join_request.NewHandleJoinRequestUseCaseRes()
	uc := handle_join_request.NewHandleJoinRequestUseCase(userRepo, groupRepo, &req, &res)
	uc.Execute()

	assert.Nil(t, res.Err)
	assert.Len(t, group.JoinRequests, 0)
	assert.Len(t, group.Members, 1)
	assert.Equal(t, group.Members[0].UserId, requesterId)
	assert.Equal(t, group.Members[0].ApprovedBy, approverId)
	assert.Equal(t, group.Members[0].InvitedBy, uuid.Nil)
}

func TestRejectJoinRequestByOwner(t *testing.T) {
	userRepo, groupRepo := setup(t)

	requesterId := uuid.New()
	ownerId := uuid.New()
	groupId := uuid.New()

	requester := entity.NewUser(requesterId, "requester", "Requester", "requester@email.com", false)
	owner := entity.NewUser(ownerId, "owner", "Owner", "owner@email.com", false)
	group := entity.NewGroup(groupId, "my group", owner, entity_enums.GROUP_PUBLIC)
	group.AddJoinRequest(requesterId)

	userRepo.EXPECT().GetUserById(requesterId).Return(requester, nil)
	userRepo.EXPECT().GetUserById(ownerId).Return(owner, nil)
	groupRepo.EXPECT().GetGroupById(groupId).Return(group, nil)
	groupRepo.EXPECT().Save(gomock.AssignableToTypeOf(&entity.Group{})).Do(
		func(arg *entity.Group) { group = arg },
	)

	req := handle_join_request.NewHandleJoinRequestUseCaseReq(
		requesterId,
		groupId,
		handle_join_request.REJECT_JOIN_REQUEST,
		ownerId,
	)
	res := handle_join_request.NewHandleJoinRequestUseCaseRes()
	uc := handle_join_request.NewHandleJoinRequestUseCase(userRepo, groupRepo, &req, &res)
	uc.Execute()

	assert.Nil(t, res.Err)
	assert.Len(t, group.JoinRequests, 0)
	assert.Len(t, group.Members, 0)
}

func TestRejectJoinRequestByGroupAdmin(t *testing.T) {
	userRepo, groupRepo := setup(t)

	requesterId := uuid.New()
	approverId := uuid.New()
	ownerId := uuid.New()
	groupId := uuid.New()

	requester := entity.NewUser(requesterId, "requester", "Requester", "requester@email.com", false)
	approver := entity.NewUser(approverId, "approver", "Approver", "approver@email.com", false)
	owner := entity.NewUser(ownerId, "owner", "Owner", "owner@email.com", false)
	group := entity.NewGroup(groupId, "my group", owner, entity_enums.GROUP_PUBLIC)
	group.AddAdmin(approverId, ownerId)
	group.AddJoinRequest(requesterId)

	userRepo.EXPECT().GetUserById(requesterId).Return(requester, nil)
	userRepo.EXPECT().GetUserById(approverId).Return(approver, nil)
	userRepo.EXPECT().GetUserById(ownerId).Return(owner, nil)
	groupRepo.EXPECT().GetGroupById(groupId).Return(group, nil)
	groupRepo.EXPECT().Save(gomock.AssignableToTypeOf(&entity.Group{})).Do(
		func(arg *entity.Group) { group = arg },
	)

	req := handle_join_request.NewHandleJoinRequestUseCaseReq(
		requesterId,
		groupId,
		handle_join_request.REJECT_JOIN_REQUEST,
		approverId,
	)
	res := handle_join_request.NewHandleJoinRequestUseCaseRes()
	uc := handle_join_request.NewHandleJoinRequestUseCase(userRepo, groupRepo, &req, &res)
	uc.Execute()

	assert.Nil(t, res.Err)
	assert.Len(t, group.JoinRequests, 0)
	assert.Len(t, group.Members, 0)
}

func TestTryToAcceptJoinRequestWithoutPermission(t *testing.T) {
	userRepo, groupRepo := setup(t)

	requesterId := uuid.New()
	approverId := uuid.New()
	ownerId := uuid.New()
	groupId := uuid.New()

	requester := entity.NewUser(requesterId, "requester", "Requester", "requester@email.com", false)
	approver := entity.NewUser(approverId, "approver", "Approver", "approver@email.com", false)
	owner := entity.NewUser(ownerId, "owner", "Owner", "owner@email.com", false)
	group := entity.NewGroup(groupId, "my group", owner, entity_enums.GROUP_PUBLIC)
	group.AddMember(approverId, uuid.Nil, ownerId)
	group.AddJoinRequest(requesterId)

	userRepo.EXPECT().GetUserById(requesterId).Return(requester, nil)
	userRepo.EXPECT().GetUserById(approverId).Return(approver, nil)
	userRepo.EXPECT().GetUserById(ownerId).Return(owner, nil)
	groupRepo.EXPECT().GetGroupById(groupId).Return(group, nil)
	groupRepo.EXPECT().Save(gomock.AssignableToTypeOf(&entity.Group{})).Do(
		func(arg *entity.Group) { group = arg },
	)

	req := handle_join_request.NewHandleJoinRequestUseCaseReq(
		requesterId,
		groupId,
		handle_join_request.ACCEPT_JOIN_REQUEST,
		approverId,
	)
	res := handle_join_request.NewHandleJoinRequestUseCaseRes()
	uc := handle_join_request.NewHandleJoinRequestUseCase(userRepo, groupRepo, &req, &res)
	uc.Execute()

	assert.ErrorIs(t, res.Err, handle_join_request.ErrApproverHasNoPermission)
	assert.Len(t, group.JoinRequests, 1)
}

func TestTryToRejectJoinRequestWithoutPermission(t *testing.T) {
	userRepo, groupRepo := setup(t)

	requesterId := uuid.New()
	approverId := uuid.New()
	ownerId := uuid.New()
	groupId := uuid.New()

	requester := entity.NewUser(requesterId, "requester", "Requester", "requester@email.com", false)
	approver := entity.NewUser(approverId, "approver", "Approver", "approver@email.com", false)
	owner := entity.NewUser(ownerId, "owner", "Owner", "owner@email.com", false)
	group := entity.NewGroup(groupId, "my group", owner, entity_enums.GROUP_PUBLIC)
	group.AddMember(approverId, uuid.Nil, ownerId)
	group.AddJoinRequest(requesterId)

	userRepo.EXPECT().GetUserById(requesterId).Return(requester, nil)
	userRepo.EXPECT().GetUserById(approverId).Return(approver, nil)
	userRepo.EXPECT().GetUserById(ownerId).Return(owner, nil)
	groupRepo.EXPECT().GetGroupById(groupId).Return(group, nil)
	groupRepo.EXPECT().Save(gomock.AssignableToTypeOf(&entity.Group{})).Do(
		func(arg *entity.Group) { group = arg },
	)

	req := handle_join_request.NewHandleJoinRequestUseCaseReq(
		requesterId,
		groupId,
		handle_join_request.REJECT_JOIN_REQUEST,
		approverId,
	)
	res := handle_join_request.NewHandleJoinRequestUseCaseRes()
	uc := handle_join_request.NewHandleJoinRequestUseCase(userRepo, groupRepo, &req, &res)
	uc.Execute()

	assert.ErrorIs(t, res.Err, handle_join_request.ErrApproverHasNoPermission)
	assert.Len(t, group.JoinRequests, 1)
}

func TestAcceptNonExistJoinRequest(t *testing.T) {
	userRepo, groupRepo := setup(t)

	requesterId := uuid.New()
	ownerId := uuid.New()
	groupId := uuid.New()

	requester := entity.NewUser(requesterId, "requester", "Requester", "requester@email.com", false)
	owner := entity.NewUser(ownerId, "owner", "Owner", "owner@email.com", false)
	group := entity.NewGroup(groupId, "my group", owner, entity_enums.GROUP_PUBLIC)

	userRepo.EXPECT().GetUserById(requesterId).Return(requester, nil)
	userRepo.EXPECT().GetUserById(ownerId).Return(owner, nil)
	groupRepo.EXPECT().GetGroupById(groupId).Return(group, nil)
	groupRepo.EXPECT().Save(gomock.AssignableToTypeOf(&entity.Group{})).Do(
		func(arg *entity.Group) { group = arg },
	)

	req := handle_join_request.NewHandleJoinRequestUseCaseReq(
		requesterId,
		groupId,
		handle_join_request.REJECT_JOIN_REQUEST,
		ownerId,
	)
	res := handle_join_request.NewHandleJoinRequestUseCaseRes()
	uc := handle_join_request.NewHandleJoinRequestUseCase(userRepo, groupRepo, &req, &res)
	uc.Execute()

	assert.ErrorIs(t, res.Err, handle_join_request.ErrJoinRequestNotFound)
	assert.Len(t, group.JoinRequests, 0)
	assert.Len(t, group.Members, 0)
}
