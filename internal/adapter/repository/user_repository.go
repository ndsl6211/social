package repository

import (
	"mashu.example/internal/adapter/repository/entity"
	"mashu.example/internal/model"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"mashu.example/internal/usecase/repository"
)

type userRepo struct {
	db *gorm.DB
}

func (ur *userRepo) GetUserById(userId uuid.UUID) (*model.User, error) {
	// get user
	user := &entity.UserEntity{}
	if err := ur.db.
		Where("users.id = ?", userId).
		First(user).Error; err != nil {
		return nil, err
	}

	// get follow relation
	if err := ur.db.
		Preload("Follows").
		Where("users.id = ?", userId).
		First(user).Error; err != nil {
		logrus.Errorf("failed to get follow relation by user id %s: %v", userId, err)
		return nil, err
	}

	return user.ToUser(), nil
}

func (ur *userRepo) GetUserByUserName(userName string) (*model.User, error) {
	// get user
	user := &entity.UserEntity{}
	if err := ur.db.
		Where("users.name = ?", userName).
		First(user).Error; err != nil {
		logrus.Errorf("failed to get user by user name %s: %v", userName, err)
		return nil, err
	}

	// get follow relation
	if err := ur.db.
		Preload("Follows").
		Where("users.name = ?", userName).
		First(user).Error; err != nil {
		logrus.Errorf("failed to get follow relation by user name %s: %v", userName, err)
		return nil, err
	}

	return user.ToUser(), nil
}

func (ur *userRepo) GetUserByDiscordUserId(discordUserId string) (*model.User, error) {
	// get user
	user := &entity.UserEntity{}
	if err := ur.db.
		Where("users.discord_user_id = ?", discordUserId).
		First(user).Error; err != nil {
		logrus.Errorf("failed to get user by discord user id %s: %v", discordUserId, err)
		return nil, err
	}

	// get follow relation
	if err := ur.db.
		Preload("Follows").
		Where("users.discord_user_id = ?", discordUserId).
		First(user).Error; err != nil {
		logrus.Errorf("failed to get follow relation by discord user id %s: %v", discordUserId, err)
		return nil, err
	}

	return user.ToUser(), nil
}

func (ur *userRepo) Save(user *model.User) error {
	// build follow status
	var userFollowsEntity []*entity.UserFollowsEntity

	// store follow req and mark status as REQUESTED
	for _, followReq := range user.FollowRequests {
		userFollowsEntity = append(userFollowsEntity, &entity.UserFollowsEntity{
			User:     followReq.From,
			Follower: followReq.To,
			Status:   entity.USER_FOLLOWS_REQUESTED,
		})
	}

	// store followers and mark status as FOLLOWING
	for _, follower := range user.Followers {
		userFollowsEntity = append(userFollowsEntity, &entity.UserFollowsEntity{
			User:     user.ID,
			Follower: follower,
			Status:   entity.USER_FOLLOWS_FOLLOWING,
		})
	}

	// store followings and mark status as FOLLOWING
	for _, following := range user.Followings {
		userFollowsEntity = append(userFollowsEntity, entity.NewUserFollowsEntity(
			following,
			user.ID,
			entity.USER_FOLLOWS_FOLLOWING,
		))
	}

	// save user
	userEntity := entity.NewUserEntity(user)
	if err := ur.db.Transaction(func(tx *gorm.DB) error {
		if err := ur.db.Save(userEntity).Error; err != nil {
			logrus.Error("failed to save UserEntity:", err)
			return err
		}

		if len(userFollowsEntity) != 0 {
			if err := ur.db.Save(userFollowsEntity).Error; err != nil {
				logrus.Error("failed to save UserFollowsEntity:", err)
				return err
			}
		}

		return nil
	}); err != nil {
		logrus.Error("failed to execute DB transaction:", err)
		return err
	}

	return nil
}

func NewUserRepository(db *gorm.DB) repository.UserRepo {
	db.AutoMigrate(&entity.UserEntity{}, &entity.UserFollowsEntity{})

	return &userRepo{db}
}
