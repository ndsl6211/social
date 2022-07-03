package repository

import (
	"errors"

	"mashu.example/internal/adapter/datamapper/user_data_mapper"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"mashu.example/internal/entity"
	"mashu.example/internal/usecase/repository"
)

type userRepo struct {
	db *gorm.DB
}

func (ur *userRepo) GetUserById(userId uuid.UUID) (*entity.User, error) {
	// get user
	userData := &user_data_mapper.UserDataMapper{}
	if err := ur.db.
		Where("users.id = ?", userId).
		First(userData).Error; err != nil {
		return nil, err
	}

	// get follow relation
	if err := ur.db.
		Joins("JOIN users ON (follows.follower_id = users.id OR follows.user_id = users.id)").
		Where("users.id = ?", userId).
		First(&userData.Follows).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}
	user := userData.ToUser()

	return user, nil
}

func (ur *userRepo) GetUserByUserName(username string) (*entity.User, error) {
	userData := &user_data_mapper.UserDataMapper{}
	if err := ur.db.
		Where("users.name = ?", username).
		First(userData).Error; err != nil {
		return nil, err
	}

	// get follow relation
	if err := ur.db.
		Joins("JOIN users ON (follows.follower_id = users.id OR follows.user_id = users.id)").
		Where("users.name = ?", username).
		First(&userData.Follows).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}
	user := userData.ToUser()

	return user, nil
}

func (ur *userRepo) Save(user *entity.User) error {
	// build follow status
	var followDataMappers []*user_data_mapper.FollowDataMapper

	// store follow req and mark status as REQUESTED
	for _, followReq := range user.FollowRequests {
		followDataMappers = append(followDataMappers, user_data_mapper.NewFollowDataMapper(
			followReq.To,
			followReq.From,
			user_data_mapper.REQUESTED,
		))
	}

	// store followers and mark status as FOLLOWING
	for _, follower := range user.Followers {
		followDataMappers = append(followDataMappers, user_data_mapper.NewFollowDataMapper(
			user.ID,
			follower,
			user_data_mapper.FOLLOWING,
		))
	}

	// store followings and mark status as FOLLOWING
	for _, following := range user.Followings {
		followDataMappers = append(followDataMappers, user_data_mapper.NewFollowDataMapper(
			following,
			user.ID,
			user_data_mapper.FOLLOWING,
		))
	}

	// save user
	userDataMapper := user_data_mapper.NewUserDataMapper(user)
	ur.db.Transaction(func(tx *gorm.DB) error {
		if err := ur.db.Save(userDataMapper).Error; err != nil {
			return err
		}

		if len(followDataMappers) != 0 {
			if err := ur.db.Save(followDataMappers).Error; err != nil {
				return err
			}
		}

		return nil
	})

	return nil
}

func NewUserRepository(db *gorm.DB) repository.UserRepo {
	db.AutoMigrate(&user_data_mapper.UserDataMapper{}, &user_data_mapper.FollowDataMapper{})

	return &userRepo{db}
}
