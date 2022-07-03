package api

import "github.com/gin-gonic/gin"

func registerPostApis(e *gin.Engine, h *restApiHandler) {
	e.POST("/post", h.createPost)
	e.PUT("/post", h.editPost)
	e.DELETE("/post", h.deletePost)
}

func (h *restApiHandler) createPost(ctx *gin.Context) {
	ctx.JSON(201, "success create")
}

func (h *restApiHandler) editPost(ctx *gin.Context) {
	ctx.JSON(201, "success edit")
}

func (h *restApiHandler) deletePost(ctx *gin.Context) {
	ctx.JSON(201, "success delete")
}
