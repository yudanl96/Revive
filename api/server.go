package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	db "github.com/yudanl96/revive/db/sqlc"
	"github.com/yudanl96/revive/token"
	"github.com/yudanl96/revive/util"
)

// serve HTTP requests
type Server struct {
	store      db.Store
	router     *gin.Engine
	tokenMaker token.Maker
	config     util.Config
}

// creates a new HTTP server and setup routing
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("fail to create token maker: %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}
	router := gin.Default()

	router.POST("/users", server.createUser)
	router.GET("/users/:username", server.getUserByUsername)
	router.GET("/posts", server.listPosts)
	router.POST("/posts", server.createPost)

	server.router = router
	return server, nil
}

// start the http server on address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
