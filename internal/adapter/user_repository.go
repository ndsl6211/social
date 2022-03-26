package adapter

import (
	"errors"

	"gorm.io/gorm"
	"mashu.example/internal/adapter/datamapper"
	"mashu.example/internal/entity"
	"mashu.example/internal/usecase/repository"
)

type userRepo struct {
	db *gorm.DB
}

func (ur *userRepo) GetUserById(userId string) (*entity.User, error) {
	// get user
	userData := &datamapper.UserDataMapper{}
	if err := ur.db.
		Where("users.id = ?", userId).
		First(&userData).Error; err != nil {
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

func (ur *userRepo) Save(user *entity.User) error {
	// save user
	userDataMapper := datamapper.NewUserDataMapper(user)
	if err := ur.db.Save(userDataMapper).Error; err != nil {
		return err
	}

	// build follow status
	var followDataMappers []*datamapper.FollowDataMapper

	// store follow req and mark status as REQUESTED
	for _, followReq := range user.FollowRequests {
		followDataMappers = append(followDataMappers, datamapper.NewFollowDataMapper(
			followReq.To,
			followReq.From,
			datamapper.REQUESTED,
		))
	}

	// store followers and mark status as FOLLOWING
	for _, follower := range user.Followers {
		followDataMappers = append(followDataMappers, datamapper.NewFollowDataMapper(
			user.ID,
			follower,
			datamapper.FOLLOWING,
		))
	}

	// store followings and mark status as FOLLOWING
	for _, following := range user.Followings {
		followDataMappers = append(followDataMappers, datamapper.NewFollowDataMapper(
			following,
			user.ID,
			datamapper.FOLLOWING,
		))
	}

	if err := ur.db.Save(followDataMappers).Error; err != nil {
		return err
	}

	return nil
}

func NewUserRepository(db *gorm.DB) repository.UserRepo {
	db.AutoMigrate(&datamapper.UserDataMapper{}, &datamapper.FollowDataMapper{})

	return &userRepo{db}
}
