package api

import "github.com/gin-gonic/gin"

func registerPostApis(e *gin.Engine, h *restApiHandler) {
	post := e.Group("/post")
	{
		post.POST("", h.createPost)
		post.PUT("", h.editPost)
		post.DELETE("", h.deletePost)
	}
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
