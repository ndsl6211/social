package main

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"mashu.example/internal/adapter/chatbot/discord"
	adapter_repository "mashu.example/internal/adapter/repository"
	"mashu.example/internal/entity"
	entity_enums "mashu.example/internal/entity/enums"
	"mashu.example/internal/usecase/repository"
	"mashu.example/internal/usecase/user/follow_user"
	"mashu.example/pkg"
)

// create users and have user1 follow user2
func createUsers(userRepo repository.UserRepo) {
	userId1 := uuid.MustParse("089b5667-85fe-4e9c-8990-1be35ca6f082")
	userId2 := uuid.MustParse("e7b81c43-d9b6-4f0c-b349-88e321115cc5")
	user1 = entity.NewUser(userId1, "mashu6211", "Mashu", "mashu@email.com", false)
	user2 = entity.NewUser(userId2, "moonnight612", "Winnie", "moonnight612@email.com", true)
	userRepo.Save(user1)
	logrus.Infof("[DATA PRELOAD] user %s created", user1.UserName)
	userRepo.Save(user2)
	logrus.Infof("[DATA PRELOAD] user %s created", user2.UserName)

	// have user1 follow user2
	req := follow_user.NewFollowUserUseCaseReq(user1.ID, user2.ID)
	res := follow_user.NewFollowUserUseCaseRes()
	uc := follow_user.NewFollowUserUseCase(userRepo, &req, &res)
	uc.Execute()
	if res.Err != nil {
		fmt.Println(res.Err.Error())
		return
	}
	logrus.Infof("[DATA PRELOAD] have %s to follow %s", user1.UserName, user2.UserName)
}

func createPost() {
	// create post with comment
	postId := uuid.MustParse("11111111-0000-0000-0000-000000000000")
	post := entity.NewPost(postId, "My First Post", "My first content", user1, nil, entity_enums.POST_PUBLIC)
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

var (
	user1 *entity.User
	user2 *entity.User
)

var (
	userRepo  repository.UserRepo
	postRepo  repository.PostRepo
	groupRepo repository.GroupRepo
	chatRepo  repository.ChatRepo
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)

	sqlite := pkg.NewSqliteGormClient()
	userRepo = adapter_repository.NewUserRepository(sqlite)
	postRepo = adapter_repository.NewPostRepository(sqlite)
	groupRepo = adapter_repository.NewGroupRepository(sqlite)
	chatRepo = adapter_repository.NewMemChatRepository()

	// redis := pkg.NewRedisClient()
	// chatRepo := adapter_repository.NewRedisChatRepository(redis)

	createUsers(userRepo)

	// // start restful api
	// engine := pkg.NewGinEngine()
	// api.RegisterWebsocketApi(engine, userRepo, chatRepo)
	// api.RegisterRestfulApis(engine, userRepo, postRepo, groupRepo)
	// engine.Run(":11000")

	// start DiscordBot
	dcRedis := pkg.NewDiscordBotUserSessionRedisClient() // violate CA, fix in the
	dcBot, err := discord.NewDiscordBot(userRepo, postRepo, groupRepo, dcRedis)
	if err != nil {
		logrus.Error("failed to create discord bot")
		return
	}
	dcBot.RegisterDiscordBotCommandHandler()
	dcBot.Start()
}
