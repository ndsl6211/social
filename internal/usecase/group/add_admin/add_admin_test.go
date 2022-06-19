package add_admin_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"mashu.example/internal/entity"
	entity_enums "mashu.example/internal/entity/enums"
	"mashu.example/internal/usecase/group/add_admin"
	"mashu.example/internal/usecase/repository/mock"
)

func setup(t *testing.T) (*mock.MockUserRepo, *mock.MockGroupRepo) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	return mock.NewMockUserRepo(mockCtrl), mock.NewMockGroupRepo(mockCtrl)
}

func TestAddAdminSucceedByOwner(t *testing.T) {
	userRepo, groupRepo := setup(t)

	memberId := uuid.New()
	groupId := uuid.New()
	ownerId := uuid.New()

	member := entity.NewUser(memberId, "member", "Member", "member@email.com", false)
	admin := entity.NewUser(ownerId, "owner", "Owner", "owner@email.com", false)

	group := entity.NewGroup(groupId, "first group", admin, entity_enums.GROUP_PUBLIC)
	group.AddMember(memberId, uuid.Nil, ownerId)

	userRepo.EXPECT().GetUserById(memberId).Return(member, nil)
	userRepo.EXPECT().GetUserById(ownerId).Return(admin, nil)
	groupRepo.EXPECT().GetGroupById(groupId).Return(group, nil)

	groupRepo.EXPECT().Save(gomock.AssignableToTypeOf(&entity.Group{})).Do(
		func(arg *entity.Group) { group = arg },
	)

	req := add_admin.NewAddAdminUseCaseReq(memberId, groupId, ownerId)
	res := add_admin.NewAddAdminUseCaseRes()
	uc := add_admin.NewAddAdminUseCase(groupRepo, userRepo, &req, &res)
	uc.Execute()

	assert.Nil(t, res.Err)
	assert.Equal(t, len(group.Admins), 1)
	assert.Equal(t, len(group.Members), 1)
}

func TestAddAdminSucceedByAdmin(t *testing.T) {
	userRepo, groupRepo := setup(t)

	memberId := uuid.New()
	groupId := uuid.New()
	adminId := uuid.New()
	ownerId := uuid.New()

	member := entity.NewUser(memberId, "member", "Member", "member@email.com", false)
	owner := entity.NewUser(ownerId, "owner", "Owner", "owner@email.com", false)
	admin := entity.NewUser(adminId, "admin", "Admin", "name@email.com", false)

	group := entity.NewGroup(groupId, "my group", owner, entity_enums.GROUP_PUBLIC)
	group.AddAdmin(adminId, ownerId)
	group.AddMember(memberId, uuid.Nil, ownerId)

	userRepo.EXPECT().GetUserById(memberId).Return(member, nil)
	userRepo.EXPECT().GetUserById(adminId).Return(admin, nil)
	groupRepo.EXPECT().GetGroupById(groupId).Return(group, nil)

	groupRepo.EXPECT().Save(gomock.AssignableToTypeOf(&entity.Group{})).Do(
		func(arg *entity.Group) { group = arg },
	)

	req := add_admin.NewAddAdminUseCaseReq(memberId, groupId, adminId)
	res := add_admin.NewAddAdminUseCaseRes()
	uc := add_admin.NewAddAdminUseCase(groupRepo, userRepo, &req, &res)
	uc.Execute()

	if res.Err != nil {
		t.Errorf("failed to execute usecase: %s", res.Err)
	}

	assert.Equal(t, len(group.Admins), 2)
	assert.Equal(t, len(group.Members), 1)
}

func TestAddNonMemberAsAdmin(t *testing.T) {
	userRepo, groupRepo := setup(t)

	memberId := uuid.New()
	groupId := uuid.New()
	ownerId := uuid.New()

	owner := entity.NewUser(ownerId, "owner", "Owner", "owner@email.com", false)
	member := entity.NewUser(ownerId, "member", "Member", "member@email.com", false)
	group := entity.NewGroup(groupId, "my group", owner, entity_enums.GROUP_PUBLIC)

	userRepo.EXPECT().GetUserById(memberId).Return(member, nil)
	userRepo.EXPECT().GetUserById(ownerId).Return(owner, nil)
	groupRepo.EXPECT().GetGroupById(groupId).Return(group, nil)

	req := add_admin.NewAddAdminUseCaseReq(memberId, groupId, ownerId)
	res := add_admin.NewAddAdminUseCaseRes()
	uc := add_admin.NewAddAdminUseCase(groupRepo, userRepo, &req, &res)
	uc.Execute()

	assert.ErrorIs(t, res.Err, add_admin.ErrNotGroupMember)
	assert.Len(t, group.Admins, 0)
}

func TestAddMemberAsAdminByNonAdmin(t *testing.T) {
	userRepo, groupRepo := setup(t)

	memberId := uuid.New()
	memberId2 := uuid.New()
	groupId := uuid.New()
	ownerId := uuid.New()

	member := entity.NewUser(memberId, "member", "Member", "member@email.com", false)
	member2 := entity.NewUser(memberId, "member2", "Member2", "member2@email.com", false)
	owner := entity.NewUser(ownerId, "owner", "Owner", "owner@email.com", false)
	group := entity.NewGroup(groupId, "my group", owner, entity_enums.GROUP_PUBLIC)
	group.AddMember(memberId, uuid.Nil, ownerId)
	group.AddMember(memberId2, uuid.Nil, ownerId)

	userRepo.EXPECT().GetUserById(memberId).Return(member, nil)
	userRepo.EXPECT().GetUserById(memberId2).Return(member2, nil)
	userRepo.EXPECT().GetUserById(ownerId).Return(owner, nil)
	groupRepo.EXPECT().GetGroupById(groupId).Return(group, nil)

	req := add_admin.NewAddAdminUseCaseReq(memberId, groupId, memberId2)
	res := add_admin.NewAddAdminUseCaseRes()
	uc := add_admin.NewAddAdminUseCase(groupRepo, userRepo, &req, &res)
	uc.Execute()

	assert.ErrorIs(t, res.Err, add_admin.ErrNotGroupOwnerOrAdmin)
	assert.Len(t, group.Admins, 0)
}

func TestAddMemberAsAdminByNonMember(t *testing.T) {
	userRepo, groupRepo := setup(t)

	memberId := uuid.New()
	nonMemberId := uuid.New()
	groupId := uuid.New()
	ownerId := uuid.New()

	member := entity.NewUser(memberId, "member", "Member", "member@email.com", false)
	nonMember := entity.NewUser(memberId, "non_member2", "NonMember", "non-member@email.com", false)
	owner := entity.NewUser(ownerId, "owner", "Owner", "owner@email.com", false)
	group := entity.NewGroup(groupId, "my group", owner, entity_enums.GROUP_PUBLIC)
	group.AddMember(memberId, uuid.Nil, ownerId)

	userRepo.EXPECT().GetUserById(memberId).Return(member, nil)
	userRepo.EXPECT().GetUserById(nonMemberId).Return(nonMember, nil)
	userRepo.EXPECT().GetUserById(ownerId).Return(owner, nil)
	groupRepo.EXPECT().GetGroupById(groupId).Return(group, nil)

	req := add_admin.NewAddAdminUseCaseReq(memberId, groupId, nonMemberId)
	res := add_admin.NewAddAdminUseCaseRes()
	uc := add_admin.NewAddAdminUseCase(groupRepo, userRepo, &req, &res)
	uc.Execute()

	assert.ErrorIs(t, res.Err, add_admin.ErrNotGroupOwnerOrAdmin)
	assert.Len(t, group.Admins, 0)
}

func TestAddAdminAsAdmin(t *testing.T) {
	userRepo, groupRepo := setup(t)

	adminId := uuid.New()
	groupId := uuid.New()
	ownerId := uuid.New()

	admin := entity.NewUser(adminId, "admin", "Admin", "admin@email.com", false)
	owner := entity.NewUser(ownerId, "owner", "Owner", "owner@email.com", false)
	group := entity.NewGroup(groupId, "my group", owner, entity_enums.GROUP_PUBLIC)
	group.AddAdmin(adminId, ownerId)

	userRepo.EXPECT().GetUserById(adminId).Return(admin, nil)
	userRepo.EXPECT().GetUserById(ownerId).Return(owner, nil)
	groupRepo.EXPECT().GetGroupById(groupId).Return(group, nil)

	req := add_admin.NewAddAdminUseCaseReq(adminId, groupId, ownerId)
	res := add_admin.NewAddAdminUseCaseRes()
	uc := add_admin.NewAddAdminUseCase(groupRepo, userRepo, &req, &res)
	uc.Execute()

	assert.ErrorIs(t, res.Err, add_admin.ErrIsAlreadyAdmin)
	assert.Len(t, group.Admins, 1)
}

func TestAddSelfAsAdmin(t *testing.T) {
	userRepo, groupRepo := setup(t)

	memberId := uuid.New()
	groupId := uuid.New()
	ownerId := uuid.New()

	member := entity.NewUser(memberId, "member", "Member", "member@email.com", false)
	owner := entity.NewUser(ownerId, "owner", "Owner", "owner@email.com", false)
	group := entity.NewGroup(groupId, "my group", owner, entity_enums.GROUP_PUBLIC)
	group.AddMember(memberId, uuid.Nil, ownerId)

	userRepo.EXPECT().GetUserById(memberId).AnyTimes().Return(member, nil)
	userRepo.EXPECT().GetUserById(ownerId).Return(owner, nil)
	groupRepo.EXPECT().GetGroupById(groupId).Return(group, nil)

	req := add_admin.NewAddAdminUseCaseReq(memberId, groupId, memberId)
	res := add_admin.NewAddAdminUseCaseRes()
	uc := add_admin.NewAddAdminUseCase(groupRepo, userRepo, &req, &res)
	uc.Execute()

	assert.ErrorIs(t, res.Err, add_admin.ErrNotGroupOwnerOrAdmin)
	assert.Len(t, group.Admins, 0)
}
