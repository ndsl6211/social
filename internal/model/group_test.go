package model_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"mashu.example/internal/entity"
	entity_enums "mashu.example/internal/entity/enums"
)

func TestCreateGroup(t *testing.T) {
	user := entity.NewUser(uuid.New(), "user", "User", "user@email.com", true)

	groupId := uuid.New()
	groupName := "new group"

	group := entity.NewGroup(groupId, groupName, user, entity_enums.GROUP_PUBLIC)

	assert.Len(t, group.Members, 0)
	assert.Len(t, group.Admins, 0)
	assert.Equal(t, groupId, group.ID)
	assert.Equal(t, groupName, group.Name)
	assert.Equal(t, user.ID, group.Owner.ID)
	assert.Equal(t, entity_enums.GROUP_PUBLIC, group.Permission)
}

func TestAddMemberToGroup(t *testing.T) {
	user := entity.NewUser(uuid.New(), "user", "User", "user@email.com", true)
	member1 := entity.NewUser(uuid.New(), "member1", "Member1", "member1@email.com", true)
	member2 := entity.NewUser(uuid.New(), "member2", "Member2", "member2@email.com", true)
	member3 := entity.NewUser(uuid.New(), "member3", "Member3", "member3@email.com", true)

	groupId := uuid.New()
	groupName := "new group"
	group := entity.NewGroup(groupId, groupName, user, entity_enums.GROUP_PUBLIC)

	group.AddMember(member1.ID, member1.ID, uuid.Nil)
	group.AddMember(member2.ID, member2.ID, uuid.Nil)
	group.AddMember(member3.ID, member3.ID, uuid.Nil)

	assert.Len(t, group.Members, 3)
}

func TestPromoteMemberAsAdmin(t *testing.T) {
	user := entity.NewUser(uuid.New(), "user", "User", "user@email.com", true)
	member1 := entity.NewUser(uuid.New(), "member1", "Member1", "member1@email.com", true)
	member2 := entity.NewUser(uuid.New(), "member2", "Member2", "member2@email.com", true)
	member3 := entity.NewUser(uuid.New(), "member3", "Member3", "member3@email.com", true)

	groupId := uuid.New()
	groupName := "new group"
	group := entity.NewGroup(groupId, groupName, user, entity_enums.GROUP_PUBLIC)

	// the user should be the member first
	group.AddMember(member1.ID, member1.ID, uuid.Nil)
	group.AddMember(member2.ID, member2.ID, uuid.Nil)
	group.AddMember(member3.ID, member3.ID, uuid.Nil)

	// and then be promoted as admin
	group.AddAdmin(member1.ID, user.ID)
	group.AddAdmin(member2.ID, user.ID)

	assert.Len(t, group.Members, 3)
	assert.Len(t, group.Admins, 2)
}

func TestRemoveMemberFromGroup(t *testing.T) {
	user := entity.NewUser(uuid.New(), "user", "User", "user@email.com", true)
	member1 := entity.NewUser(uuid.New(), "member1", "Member1", "member1@email.com", true)
	member2 := entity.NewUser(uuid.New(), "member2", "Member2", "member2@email.com", true)
	member3 := entity.NewUser(uuid.New(), "member3", "Member3", "member3@email.com", true)

	groupId := uuid.New()
	groupName := "new group"
	group := entity.NewGroup(groupId, groupName, user, entity_enums.GROUP_PUBLIC)

	group.AddMember(member1.ID, member1.ID, uuid.Nil)
	group.AddMember(member2.ID, member2.ID, uuid.Nil)
	group.AddMember(member3.ID, member3.ID, uuid.Nil)

	group.RemoveMember(member2.ID)

	assert.Len(t, group.Members, 2)
	assert.Equal(t, member1.ID, group.Members[0].UserId)
	assert.Equal(t, member3.ID, group.Members[1].UserId)
}

func TestRemoveAdminFromGroup(t *testing.T) {
	user := entity.NewUser(uuid.New(), "user", "User", "user@email.com", true)
	member1 := entity.NewUser(uuid.New(), "member1", "Member1", "member1@email.com", true)
	member2 := entity.NewUser(uuid.New(), "member2", "Member2", "member2@email.com", true)
	member3 := entity.NewUser(uuid.New(), "member3", "Member3", "member3@email.com", true)

	groupId := uuid.New()
	groupName := "new group"
	group := entity.NewGroup(groupId, groupName, user, entity_enums.GROUP_PUBLIC)

	group.AddMember(member1.ID, member1.ID, uuid.Nil)
	group.AddMember(member2.ID, member2.ID, uuid.Nil)
	group.AddMember(member3.ID, member3.ID, uuid.Nil)

	group.AddAdmin(member3.ID, user.ID)
	group.AddAdmin(member1.ID, user.ID)
	group.AddAdmin(member2.ID, user.ID)

	group.RemoveAdmin(member1.ID)

	assert.Len(t, group.Admins, 2)
	assert.Equal(t, member3.ID, group.Admins[0].UserId)
	assert.Equal(t, member2.ID, group.Admins[1].UserId)
}
