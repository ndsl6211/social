package main

import (
	"fmt"

	"gorm.io/gorm"
	"mashu.example/internal/adapter"
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
	db := pkg.NewMemoryGormClient()
	adapter.NewUserRepository(db)
	adapter.NewPostRepository(db)

	// user1 := entity.NewUser(uuid.New(), "mashu6211", "Mashu", "mashu@email.com", false)
	// user2 := entity.NewUser(uuid.New(), "moonnight612", "Winnie", "moonnight612@email.com", false)
	// userRepo.Save(user1)
	// userRepo.Save(user2)

	// req := usecase.NewFollowUserUseCaseReq(user1.ID.String(), user2.ID.String())
	// res := usecase.NewFollowUserUseCaseRes()
	// uc := usecase.NewFollowUserUseCase(userRepo, &req, &res)

	// uc.Execute()
	// if res.Err != nil {
	// 	fmt.Println("failed to execute usecase")
	// 	return
	// }

	// fmt.Println("----result----")
	// user1res, err := userRepo.GetUserById(user2.ID.String())
	// if err != nil {
	// 	fmt.Printf("fail to get result user %s\n", err.Error())
	// }

	// fmt.Printf("%+v\n", user1res)
}
