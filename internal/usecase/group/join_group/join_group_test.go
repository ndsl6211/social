package join_group_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"mashu.example/internal/entity"
	entity_enums "mashu.example/internal/entity/enums"
	"mashu.example/internal/usecase/group/join_group"
	"mashu.example/internal/usecase/repository/mock"
)

func TestJoinPublicGroup(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	owner := entity.NewUser(uuid.New(), "owner", "Owner", "owner@email.com", false)
	user := entity.NewUser(uuid.New(), "user", "User", "user@email.com", false)
	group := entity.NewGroup(uuid.New(), "group", owner, entity_enums.GROUP_PUBLIC)

	groupRepo := mock.NewMockGroupRepo(mockCtrl)
	userRepo := mock.NewMockUserRepo(mockCtrl)
	groupRepo.EXPECT().GetGroupById(group.ID).Return(group, nil)
	userRepo.EXPECT().GetUserById(user.ID).Return(user, nil)
	groupRepo.EXPECT().Save(gomock.AssignableToTypeOf(group)).Do(
		func(arg *entity.Group) { group = arg },
	)

	req := join_group.NewJoinGroupUseCaseReq(user.ID, group.ID)
	res := join_group.NewJoinGroupUseCaseRes()
	gc := join_group.NewJoinGroupUseCase(userRepo, groupRepo, &req, &res)

	gc.Execute()

	assert.Nil(t, res.Err)
	assert.Len(t, group.Members, 0)
	assert.Len(t, group.JoinRequests, 1)
	assert.Len(t, group.Admins, 0)
	assert.Equal(t, group.Owner, owner)
}
