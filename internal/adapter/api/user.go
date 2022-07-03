package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"mashu.example/internal/usecase/user/login"
	"mashu.example/internal/usecase/user/register"
)

func registerUserApis(e *gin.Engine, h *restApiHandler) {
	user := e.Group("/user")
	{
		user.POST("/register", h.register)
		user.POST("/login", h.login)
	}
}

func (h *restApiHandler) register(ctx *gin.Context) {
	type registerPayload struct {
		Username    string `json:"username" binding:"required"`
		DisplayName string `json:"displayName" binding:"required"`
		Email       string `json:"email" binding:"email"`
	}
	p := &registerPayload{}
	if err := ctx.ShouldBindJSON(p); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, newRestErrResponse(err.Error()))
		return
	}
	req := register.NewRegisterUseCaseReq(p.Username, p.DisplayName, p.Email)
	res := register.NewRegisterUseCaseRes()
	uc := register.NewRegisterUseCase(h.userRepo, req, res)

	uc.Execute()

	if res.Err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, newRestErrResponse(res.Err.Error()))
		return
	}

	ctx.Status(http.StatusCreated)
}

func (h *restApiHandler) login(ctx *gin.Context) {
	req := &login.LoginUseCaseReq{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, newRestErrResponse(err.Error()))
		return
	}
	res := login.NewLoginUseCaseRes()
	uc := login.NewLoginUseCase(h.userRepo, req, res)
	uc.Execute()

	if res.Err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, newRestErrResponse(res.Err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, res)
}
