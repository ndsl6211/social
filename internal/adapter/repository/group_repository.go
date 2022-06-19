package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"mashu.example/internal/adapter/datamapper/group_data_mapper"
	"mashu.example/internal/entity"
	"mashu.example/internal/usecase/repository"
)

type groupRepo struct {
	db *gorm.DB
}

func (gr *groupRepo) GetGroupById(groupId uuid.UUID) (*entity.Group, error) {
	//get group
	groupData := &group_data_mapper.GroupDataMapper{}
	if err := gr.db.
		Where("groups.id = ?", groupId).
		First(groupData).Error; err != nil {
		return nil, err
	}

	group := groupData.ToGroup()

	return group, nil
}

func (gr *groupRepo) GetGroupByName(groupName string) (*entity.Group, error) {
	return nil, nil
}

func (gr *groupRepo) Save(group *entity.Group) error {
	return nil
}

func (gr *groupRepo) Delete(groupId uuid.UUID) error {
	return nil
}

func NewGroupRepository(db *gorm.DB) repository.GroupRepo {
	return &groupRepo{db}
}
