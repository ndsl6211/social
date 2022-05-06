package group_data_mapper

import (
	"time"

	"github.com/google/uuid"
	"mashu.example/internal/adapter/datamapper/user_data_mapper"
	"mashu.example/internal/entity"
	entity_enums "mashu.example/internal/entity/enums"
)

type GroupDataMapper struct {
	ID   uuid.UUID
	Name string

	OwnerId uuid.UUID
	Owner   *user_data_mapper.UserDataMapper

	Permission entity_enums.GroupPermission
	Admins     []*user_data_mapper.UserDataMapper
	CreatedAt  time.Time
	Members    []JoinDataMapper
}

func (GroupDataMapper) TableName() string {
	return "groups"
}

func (g GroupDataMapper) ToGroup() *entity.Group {
	return &entity.Group{
		ID:         g.ID,
		Name:       g.Name,
		Owner:      g.Owner.ToUser(),
		Permission: g.Permission,
		CreatedAt:  g.CreatedAt,
	}
}

func NewGroupDataMapper(group *entity.Group) *GroupDataMapper {
	return &GroupDataMapper{
		ID:         group.ID,
		OwnerId:    group.Owner.ID,
		Name:       group.Name,
		Permission: group.Permission,
		CreatedAt:  group.CreatedAt,
	}
}

type JoinStatus string

const (
	REQUESTED JoinStatus = "REQUESTED"
	JOINONG   JoinStatus = "JOINING"
)

type JoinDataMapper struct {
	Group  uuid.UUID
	User   uuid.UUID
	Status JoinStatus
}

func (JoinDataMapper) TableName() string {
	return "joins"
}

func NewJoinDataMapper(
	groupId uuid.UUID,
	userId uuid.UUID,
	status JoinStatus,
) *JoinDataMapper {
	return &JoinDataMapper{groupId, userId, status}
}
