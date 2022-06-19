package delete_admin_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"mashu.example/internal/entity"
	entity_enums "mashu.example/internal/entity/enums"
	"mashu.example/internal/usecase/group/delete_admin"
	"mashu.example/internal/usecase/repository/mock"
)

func setup(t *testing.T) (*mock.MockUserRepo, *mock.MockGroupRepo) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	return mock.NewMockUserRepo(mockCtrl), mock.NewMockGroupRepo(mockCtrl)
}

func TestDeleteAdminByOwnerSucceed(t *testing.T) {
	userRepo, groupRepo := setup(t)

	adminId := uuid.New()
	groupId := uuid.New()
	ownerId := uuid.New()

	admin := entity.NewUser(adminId, "admin", "Admin", "admin@email.com", false)
	owner := entity.NewUser(ownerId, "owner", "Owner", "owner@email.com", false)

	group := entity.NewGroup(groupId, "my group", owner, entity_enums.GROUP_PUBLIC)
	group.AddMember(adminId, uuid.Nil, ownerId)
	group.AddAdmin(adminId, ownerId)

	userRepo.EXPECT().GetUserById(adminId).Return(admin, nil)
	userRepo.EXPECT().GetUserById(ownerId).Return(owner, nil)
	groupRepo.EXPECT().GetGroupById(groupId).Return(group, nil)

	groupRepo.EXPECT().Save(gomock.AssignableToTypeOf(&entity.Group{})).Do(
		func(arg *entity.Group) { group = arg },
	)

	req := delete_admin.NewDeleteAdminUseCaseReq(adminId, groupId, ownerId)
	res := delete_admin.NewDeleteAdminUseCaseRes()
	uc := delete_admin.NewDeleteAdminUseCase(userRepo, groupRepo, &req, &res)
	uc.Execute()

	assert.Nil(t, res.Err)
	assert.Len(t, group.Admins, 0)
	assert.Len(t, group.Members, 1)
}

func TestDeleteAdminByOtherAdmin(t *testing.T) {
	userRepo, groupRepo := setup(t)

	adminId1 := uuid.New()
	adminId2 := uuid.New()
	groupId := uuid.New()
	ownerId := uuid.New()

	member1 := entity.NewUser(adminId1, "member1", "Member1", "member1@email.com", false)
	member2 := entity.NewUser(adminId2, "member2", "Member2", "member2@email.com", false)
	owner := entity.NewUser(ownerId, "owner", "Owner", "owner@email.com", false)

	group := entity.NewGroup(groupId, "my group", owner, entity_enums.GROUP_PUBLIC)
	group.AddMember(adminId1, uuid.Nil, ownerId)
	group.AddMember(adminId2, uuid.Nil, ownerId)
	group.AddAdmin(adminId1, ownerId)
	group.AddAdmin(adminId2, ownerId)

	userRepo.EXPECT().GetUserById(adminId1).Return(member1, nil)
	userRepo.EXPECT().GetUserById(adminId2).Return(member2, nil)
	userRepo.EXPECT().GetUserById(ownerId).Return(owner, nil)
	groupRepo.EXPECT().GetGroupById(groupId).Return(group, nil)

	groupRepo.EXPECT().Save(gomock.AssignableToTypeOf(&entity.Group{})).Do(
		func(arg *entity.Group) { group = arg },
	)

	req := delete_admin.NewDeleteAdminUseCaseReq(adminId1, groupId, adminId2)
	res := delete_admin.NewDeleteAdminUseCaseRes()
	uc := delete_admin.NewDeleteAdminUseCase(userRepo, groupRepo, &req, &res)
	uc.Execute()

	assert.Nil(t, res.Err)
	assert.Len(t, group.Admins, 1)
	assert.Len(t, group.Members, 2)
}

func TestDeleteNonAdmin(t *testing.T) {
	userRepo, groupRepo := setup(t)

	memberId := uuid.New()
	groupId := uuid.New()
	ownerId := uuid.New()

	member := entity.NewUser(memberId, "member", "Member", "member@email.com", false)
	owner := entity.NewUser(ownerId, "owner", "Owner", "owner@email.com", false)

	group := entity.NewGroup(groupId, "my group", owner, entity_enums.GROUP_PUBLIC)
	group.AddMember(memberId, uuid.Nil, ownerId)

	userRepo.EXPECT().GetUserById(memberId).Return(member, nil)
	userRepo.EXPECT().GetUserById(ownerId).Return(owner, nil)
	groupRepo.EXPECT().GetGroupById(groupId).Return(group, nil)

	groupRepo.EXPECT().Save(gomock.AssignableToTypeOf(&entity.Group{})).Do(
		func(arg *entity.Group) { group = arg },
	)

	req := delete_admin.NewDeleteAdminUseCaseReq(memberId, groupId, ownerId)
	res := delete_admin.NewDeleteAdminUseCaseRes()
	uc := delete_admin.NewDeleteAdminUseCase(userRepo, groupRepo, &req, &res)
	uc.Execute()

	assert.ErrorIs(t, res.Err, delete_admin.ErrAdminNotFound)
	assert.Len(t, group.Admins, 0)
	assert.Len(t, group.Members, 1)
}

func TestDeleteAdminByNonAdmin(t *testing.T) {
	userRepo, groupRepo := setup(t)

	memberId := uuid.New()
	adminId := uuid.New()
	groupId := uuid.New()
	ownerId := uuid.New()

	member := entity.NewUser(memberId, "member", "Member", "member@email.com", false)
	admin := entity.NewUser(adminId, "admin", "Admin", "admin@email.com", false)
	owner := entity.NewUser(ownerId, "owner", "Owner", "owner@email.com", false)

	group := entity.NewGroup(groupId, "my group", owner, entity_enums.GROUP_PUBLIC)
	group.AddMember(memberId, uuid.Nil, ownerId)
	group.AddMember(adminId, uuid.Nil, ownerId)
	group.AddAdmin(adminId, ownerId)

	userRepo.EXPECT().GetUserById(memberId).Return(member, nil)
	userRepo.EXPECT().GetUserById(adminId).Return(admin, nil)
	userRepo.EXPECT().GetUserById(ownerId).Return(owner, nil)
	groupRepo.EXPECT().GetGroupById(groupId).Return(group, nil)

	groupRepo.EXPECT().Save(gomock.AssignableToTypeOf(&entity.Group{})).Do(
		func(arg *entity.Group) { group = arg },
	)

	req := delete_admin.NewDeleteAdminUseCaseReq(adminId, groupId, memberId)
	res := delete_admin.NewDeleteAdminUseCaseRes()
	uc := delete_admin.NewDeleteAdminUseCase(userRepo, groupRepo, &req, &res)
	uc.Execute()

	assert.ErrorIs(t, res.Err, delete_admin.ErrNotGroupOwnerOrAdmin)
	assert.Len(t, group.Admins, 1)
	assert.Len(t, group.Members, 2)
}

func TestDeleteAdminByNonMember(t *testing.T) {
	userRepo, groupRepo := setup(t)

	nonMemberId := uuid.New()
	adminId := uuid.New()
	groupId := uuid.New()
	ownerId := uuid.New()

	member := entity.NewUser(nonMemberId, "non_member", "NonMember", "non_member@email.com", false)
	admin := entity.NewUser(adminId, "admin", "Admin", "admin@email.com", false)
	owner := entity.NewUser(ownerId, "owner", "Owner", "owner@email.com", false)

	group := entity.NewGroup(groupId, "my group", owner, entity_enums.GROUP_PUBLIC)
	group.AddMember(adminId, uuid.Nil, ownerId)
	group.AddAdmin(adminId, ownerId)

	userRepo.EXPECT().GetUserById(nonMemberId).Return(member, nil)
	userRepo.EXPECT().GetUserById(adminId).Return(admin, nil)
	userRepo.EXPECT().GetUserById(ownerId).Return(owner, nil)
	groupRepo.EXPECT().GetGroupById(groupId).Return(group, nil)

	groupRepo.EXPECT().Save(gomock.AssignableToTypeOf(&entity.Group{})).Do(
		func(arg *entity.Group) { group = arg },
	)

	req := delete_admin.NewDeleteAdminUseCaseReq(adminId, groupId, nonMemberId)
	res := delete_admin.NewDeleteAdminUseCaseRes()
	uc := delete_admin.NewDeleteAdminUseCase(userRepo, groupRepo, &req, &res)
	uc.Execute()

	assert.ErrorIs(t, res.Err, delete_admin.ErrNotGroupOwnerOrAdmin)
	assert.Len(t, group.Admins, 1)
	assert.Len(t, group.Members, 1)
}
