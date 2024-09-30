package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	db "github.com/yudanl96/revive/db/sqlc"
	"github.com/yudanl96/revive/util"
)

type listPostsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=20"`
}

func (server *Server) listPosts(ctx *gin.Context) {
	var request listPostsRequest
	if err := ctx.ShouldBindQuery(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListPostsParams{
		Limit:  request.PageSize,
		Offset: (request.PageID - 1) * request.PageSize,
	}

	posts, err := server.store.ListPosts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, posts)

}

type createPostsRequest struct {
	UserID      string `json:"user_id" binding:"required"`
	GenAIPrompt string `json:"prompt" binding:"required"` //	Description string `json:"descr" binding:"required,email"`
	Price       int32  `json:"price" binding:"required,numeric,gte=0"`
}

func (server *Server) createPost(ctx *gin.Context) {
	var request createPostsRequest
	if err := ctx.ShouldBindBodyWithJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	description, err := util.GenerateText(request.GenAIPrompt, 30)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreatePostParams{
		ID:          uuid.NewString(),
		UserID:      request.UserID,
		Price:       request.Price,
		Description: description,
	}

	fmt.Println(arg.Description)

	err = server.store.CreatePost(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, arg.ID)
}
