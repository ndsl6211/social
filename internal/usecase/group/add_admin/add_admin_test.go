package add_admin_test

import (
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"mashu.example/internal/entity"
	"mashu.example/internal/entity/enums/group_permission"
	"mashu.example/internal/usecase/group/add_admin"
	"mashu.example/internal/usecase/repository/mock"
	"testing"
	"time"
)

func TestAddAdminSucceedByOwner(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	memberId := uuid.New()
	groupId := uuid.New()
	adminId := uuid.New()

	member := entity.NewUser(
		memberId,
		"member",
		"member display name",
		"member@email.com",
		false,
	)
	admin := entity.NewUser(
		adminId,
		"name",
		"name display name",
		"name@email.com",
		false,
	)
	group := &entity.Group{
		ID:             groupId,
		Name:           "first group",
		Owner:          admin,
		Permission:     group_permission.UNPUBLIC,
		Admins:         nil,
		CreatedAt:      time.Time{},
		Members:        []uuid.UUID{memberId},
		JoinRequests:   nil,
		InviteRequests: nil,
	}

	userRepo := mock.NewMockUserRepo(mockCtrl)
	userRepo.EXPECT().GetUserById(memberId).Return(member, nil)
	userRepo.EXPECT().GetUserById(adminId).Return(admin, nil)
	groupRepo := mock.NewMockGroupRepo(mockCtrl)
	groupRepo.EXPECT().GetGroupById(groupId).Return(group, nil)

	groupRepo.EXPECT().Save(gomock.AssignableToTypeOf(&entity.Group{})).Do(
		func(arg *entity.Group) { group = arg },
	)

	req := add_admin.NewAddAdminUseCaseReq(
		memberId,
		groupId,
		adminId,
	)

	res := add_admin.NewAddAdminUseCaseRes()
	gc := add_admin.NewAddAdminUseCase(groupRepo, userRepo, &req, &res)
	gc.Execute()

	if res.Err != nil {
		t.Errorf("failed to execute usecase")
	}

	assert.Equal(t, len(group.Admins), 1)
	assert.Equal(t, len(group.Members), 1)
}

func TestAddAdminSucceedByAdmin(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	memberId := uuid.New()
	groupId := uuid.New()
	adminId := uuid.New()
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
	admin := entity.NewUser(
		adminId,
		"name",
		"name display name",
		"name@email.com",
		false,
	)
	group := &entity.Group{
		ID:             groupId,
		Name:           "first group",
		Owner:          owner,
		Permission:     group_permission.UNPUBLIC,
		Admins:         []uuid.UUID{admin.ID},
		CreatedAt:      time.Time{},
		Members:        []uuid.UUID{member.ID},
		JoinRequests:   nil,
		InviteRequests: nil,
	}

	userRepo := mock.NewMockUserRepo(mockCtrl)
	userRepo.EXPECT().GetUserById(memberId).Return(member, nil)
	userRepo.EXPECT().GetUserById(adminId).Return(admin, nil)
	groupRepo := mock.NewMockGroupRepo(mockCtrl)
	groupRepo.EXPECT().GetGroupById(groupId).Return(group, nil)

	groupRepo.EXPECT().Save(gomock.AssignableToTypeOf(&entity.Group{})).Do(
		func(arg *entity.Group) { group = arg },
	)

	req := add_admin.NewAddAdminUseCaseReq(
		memberId,
		groupId,
		adminId,
	)

	res := add_admin.NewAddAdminUseCaseRes()
	gc := add_admin.NewAddAdminUseCase(groupRepo, userRepo, &req, &res)
	gc.Execute()

	if res.Err != nil {
		t.Errorf("failed to execute usecase")
	}

	assert.Equal(t, len(group.Admins), 2)
	assert.Equal(t, len(group.Members), 1)
}
func TestAddAdminFail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	memberId := uuid.New()
	groupId := uuid.New()
	adminId := uuid.New()

	member := entity.NewUser(
		memberId,
		"member",
		"member display name",
		"member@email.com",
		false,
	)
	admin := entity.NewUser(
		adminId,
		"name",
		"name display name",
		"name@email.com",
		false,
	)
	group := &entity.Group{
		ID:             groupId,
		Name:           "third group",
		Owner:          admin,
		Permission:     group_permission.UNPUBLIC,
		Admins:         nil,
		CreatedAt:      time.Time{},
		Members:        nil,
		JoinRequests:   nil,
		InviteRequests: nil,
	}

	userRepo := mock.NewMockUserRepo(mockCtrl)
	userRepo.EXPECT().GetUserById(memberId).Return(member, nil)
	userRepo.EXPECT().GetUserById(adminId).Return(admin, nil)
	groupRepo := mock.NewMockGroupRepo(mockCtrl)
	groupRepo.EXPECT().GetGroupById(groupId).Return(group, nil)

	req := add_admin.NewAddAdminUseCaseReq(
		memberId,
		groupId,
		adminId,
	)

	res := add_admin.NewAddAdminUseCaseRes()
	gc := add_admin.NewAddAdminUseCase(groupRepo, userRepo, &req, &res)
	gc.Execute()

	assert.Equal(t, len(group.Admins), 0)
}
