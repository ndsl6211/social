package delete_admin_test

import (
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"mashu.example/internal/entity"
	"mashu.example/internal/entity/enums/group_permission"
	"mashu.example/internal/usecase/group/delete_admin"
	"mashu.example/internal/usecase/repository/mock"
	"testing"
	"time"
)

func TestDeleteAdminByOwnerSucceed(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	memberId := uuid.New()
	groupId := uuid.New()
	ownerId := uuid.New()

	member := entity.NewUser(
		memberId,
		"member",
		"member display name",
		"member@email.com",
		false,
	)
	owner := entity.NewUser(
		ownerId,
		"owner",
		"owner display name",
		"owner@email.com",
		false,
	)

	group := &entity.Group{
		ID:             groupId,
		Name:           "first group",
		Owner:          owner,
		Permission:     group_permission.UNPUBLIC,
		Admins:         []uuid.UUID{memberId},
		CreatedAt:      time.Time{},
		Members:        []uuid.UUID{memberId},
		JoinRequests:   nil,
		InviteRequests: nil,
	}

	userRepo := mock.NewMockUserRepo(mockCtrl)
	userRepo.EXPECT().GetUserById(memberId).Return(member, nil)
	userRepo.EXPECT().GetUserById(ownerId).Return(owner, nil)
	groupRepo := mock.NewMockGroupRepo(mockCtrl)
	groupRepo.EXPECT().GetGroupById(groupId).Return(group, nil)

	groupRepo.EXPECT().Save(gomock.AssignableToTypeOf(&entity.Group{})).Do(
		func(arg *entity.Group) { group = arg },
	)

	req := delete_admin.NewDeleteAdminUseCaseReq(
		memberId,
		groupId,
		ownerId,
	)

	res := delete_admin.NewDeleteAdminUseCaseRes()
	gc := delete_admin.NewDeleteAdminUseCase(userRepo, groupRepo, &req, &res)
	gc.Execute()

	if res.Err != nil {
		t.Errorf("failed to execute usecase")
	}

	assert.Equal(t, len(group.Admins), 0)
	assert.Equal(t, len(group.Members), 1)
}
