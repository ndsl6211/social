package handle_join_request_test

import (
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"mashu.example/internal/entity"
	"mashu.example/internal/entity/enums/group_permission"
	"mashu.example/internal/usecase/group/handle_join_request"
	"mashu.example/internal/usecase/repository/mock"
	"testing"
	"time"
)

func TestAcceptJoinRequest(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	requesterId := uuid.New()
	groupId := uuid.New()
	approverId := uuid.New()

	requester := entity.NewUser(
		requesterId,
		"requester",
		"requester display name",
		"requester@email.com",
		false,
	)

	approver := entity.NewUser(
		approverId,
		"approver",
		"approver display name",
		"approver@email.com",
		false,
	)

	group := &entity.Group{
		ID:         groupId,
		Name:       "first group",
		Owner:      approver,
		Permission: group_permission.UNPUBLIC,
		Admins:     nil,
		CreatedAt:  time.Time{},
		Members:    nil,
		JoinRequests: []*entity.JoinRequest{
			{
				Requester: requesterId,
				Group:     groupId,
			},
		},
		InviteRequests: nil,
	}

	userRepo := mock.NewMockUserRepo(mockCtrl)
	userRepo.EXPECT().GetUserById(requesterId).Return(requester, nil)
	userRepo.EXPECT().GetUserById(approverId).Return(approver, nil)
	groupRepo := mock.NewMockGroupRepo(mockCtrl)
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
	gc := handle_join_request.NewHandleJoinRequestUseCase(userRepo, groupRepo, &req, &res)
	gc.Execute()

	if res.Err != nil {
		t.Errorf("failed to execute usecase")
	}

	assert.Equal(t, len(group.JoinRequests), 0)
	assert.Equal(t, len(group.Members), 1)
	assert.Equal(t, group.Owner.ID, approver.ID)
}

func TestRejectJoinRequest(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	requesterId := uuid.New()
	groupId := uuid.New()
	approverId := uuid.New()
	ownerId := uuid.New()

	owner := entity.NewUser(
		ownerId,
		"owner",
		"owner display name",
		"owner@email.com",
		false,
	)

	requester := entity.NewUser(
		requesterId,
		"requester",
		"requester display name",
		"requester@email.com",
		false,
	)

	approver := entity.NewUser(
		approverId,
		"approver",
		"approver display name",
		"approver@email.com",
		false,
	)

	group := &entity.Group{
		ID:         groupId,
		Name:       "second group",
		Owner:      owner,
		Permission: group_permission.UNPUBLIC,
		Admins:     []uuid.UUID{approverId},
		CreatedAt:  time.Time{},
		Members:    nil,
		JoinRequests: []*entity.JoinRequest{
			{
				Requester: requesterId,
				Group:     groupId,
			},
		},
		InviteRequests: nil,
	}

	userRepo := mock.NewMockUserRepo(mockCtrl)
	userRepo.EXPECT().GetUserById(requesterId).Return(requester, nil)
	userRepo.EXPECT().GetUserById(approverId).Return(approver, nil)
	groupRepo := mock.NewMockGroupRepo(mockCtrl)
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
	gc := handle_join_request.NewHandleJoinRequestUseCase(userRepo, groupRepo, &req, &res)
	gc.Execute()

	if res.Err != nil {
		t.Errorf("failed to execute usecase")
	}

	assert.Equal(t, len(group.JoinRequests), 0)
	assert.Equal(t, len(group.Members), 0)
	assert.Equal(t, group.Owner.ID, owner.ID)
}

func TestDenyJoinRequest(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	requesterId := uuid.New()
	groupId := uuid.New()
	approverId := uuid.New()
	ownerId := uuid.New()

	owner := entity.NewUser(
		ownerId,
		"owner",
		"owner display name",
		"owner@email.com",
		false,
	)

	requester := entity.NewUser(
		requesterId,
		"requester",
		"requester display name",
		"requester@email.com",
		false,
	)

	approver := entity.NewUser(
		approverId,
		"approver",
		"approver display name",
		"approver@email.com",
		false,
	)

	group := &entity.Group{
		ID:         groupId,
		Name:       "second group",
		Owner:      owner,
		Permission: group_permission.UNPUBLIC,
		Admins:     nil,
		CreatedAt:  time.Time{},
		Members:    nil,
		JoinRequests: []*entity.JoinRequest{
			{
				Requester: requesterId,
				Group:     groupId,
			},
		},
		InviteRequests: nil,
	}

	userRepo := mock.NewMockUserRepo(mockCtrl)
	userRepo.EXPECT().GetUserById(requesterId).Return(requester, nil)
	userRepo.EXPECT().GetUserById(approverId).Return(approver, nil)
	groupRepo := mock.NewMockGroupRepo(mockCtrl)
	groupRepo.EXPECT().GetGroupById(groupId).Return(group, nil)

	req := handle_join_request.NewHandleJoinRequestUseCaseReq(
		requesterId,
		groupId,
		handle_join_request.ACCEPT_JOIN_REQUEST,
		approverId,
	)

	res := handle_join_request.NewHandleJoinRequestUseCaseRes()
	gc := handle_join_request.NewHandleJoinRequestUseCase(userRepo, groupRepo, &req, &res)
	gc.Execute()

	assert.Equal(t, len(group.JoinRequests), 1)
	assert.Equal(t, res.Err.Error(), "permission denied")
	assert.Equal(t, len(group.Members), 0)
}
