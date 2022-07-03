package api

import (
	"github.com/gin-gonic/gin"
	"mashu.example/internal/usecase/repository"
)

func newRestErrResponse(message string) map[string]string {
	return map[string]string{"error": message}
}

type restApiHandler struct {
	userRepo  repository.UserRepo
	postRepo  repository.PostRepo
	groupRepo repository.GroupRepo
}

func RegisterRestfulApis(
	e *gin.Engine,
	userRepo repository.UserRepo,
	postRepo repository.PostRepo,
	groupRepo repository.GroupRepo,
) {
	h := newRestApiHandler(userRepo, postRepo, groupRepo)

	registerGroupApis(e, h)
	registerPostApis(e, h)
	registerUserApis(e, h)
}

func newRestApiHandler(
	userRepo repository.UserRepo,
	postRepo repository.PostRepo,
	groupRepo repository.GroupRepo,
) *restApiHandler {
	return &restApiHandler{userRepo, postRepo, groupRepo}
}
