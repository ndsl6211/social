package main

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"mashu.example/internal/adapter"
	"mashu.example/internal/entity"
	"mashu.example/internal/entity/enums/post_permission"
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

	// create users
	user1 := entity.NewUser(uuid.New(), "mashu6211", "Mashu", "mashu@email.com", false)
	user2 := entity.NewUser(uuid.New(), "moonnight612", "Winnie", "moonnight612@email.com", false)
	// user3 := entity.NewUser(uuid.New(), "moonnight612", "Winnie", "moonnight612@email.com", false)
	userRepo.Save(user1)
	userRepo.Save(user2)
	// if err := userRepo.Save(user3); err != nil {
	// 	fmt.Println("ERROR!")
	// 	return
	// }

	// have user1 follow user2
	req := follow_user.NewFollowUserUseCaseReq(user1.ID, user2.ID)
	res := follow_user.NewFollowUserUseCaseRes()
	uc := follow_user.NewFollowUserUseCase(userRepo, &req, &res)
	uc.Execute()
	if res.Err != nil {
		fmt.Println(res.Err.Error())
		return
	}

	// create post with comment
	postId := uuid.MustParse("11111111-0000-0000-0000-000000000000")
	post := entity.NewPost(postId, "My First Post", "My first content", user1, post_permission.PUBLIC)
	post.Comments = append(post.Comments, &entity.Comment{
		ID:        uuid.New(),
		Owner:     user1,
		Post:      post,
		Content:   "my first comment",
		CreatedAt: time.Now(),
	})
	post.Comments = append(post.Comments, &entity.Comment{
		ID:        uuid.New(),
		Owner:     user1,
		Post:      post,
		Content:   "my second comment",
		CreatedAt: time.Now(),
	})

	postRepo.Save(post)
}
