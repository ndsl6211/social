package main

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"mashu.example/internal/adapter"
	"mashu.example/internal/entity"
	"mashu.example/internal/usecase/user/follow_user"
	"mashu.example/pkg"
)

func useful(db *gorm.DB) {
	type User struct {
		ID        int
		Followers []*User `gorm:"many2many:user_followers"`
	}

	type UserFollowers struct {
		UserID     int `gorm:"primaryKey"`
		FollowerID int `gorm:"primaryKey"`

		Role int
	}

	serr := db.SetupJoinTable(&User{}, "Followers", &UserFollowers{})
	db.AutoMigrate(&User{}, &UserFollowers{})

	if serr != nil {
		fmt.Println("err when setup:", serr)
	}
}

func main() {
	db := pkg.NewSqliteGormClient()
	userRepo := adapter.NewUserRepository(db)
	postRepo := adapter.NewPostRepository(db)

	user1 := entity.NewUser(uuid.New(), "mashu6211", "Mashu", "mashu@email.com", false)
	user2 := entity.NewUser(uuid.New(), "moonnight612", "Winnie", "moonnight612@email.com", false)
	userRepo.Save(user1)
	userRepo.Save(user2)

	// follow user
	req := follow_user.NewFollowUserUseCaseReq(user1.ID.String(), user2.ID.String())
	res := follow_user.NewFollowUserUseCaseRes()
	uc := follow_user.NewFollowUserUseCase(userRepo, &req, &res)
	uc.Execute()
	if res.Err != nil {
		fmt.Println("failed to execute usecase")
		return
	}

	// create post
	postId := uuid.MustParse("11111111-0000-0000-0000-000000000000")
	post := entity.NewPost(postId, "title", "content", user1, true)
	postRepo.Save(post)
}
