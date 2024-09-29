package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/yudanl96/revive/db/sqlc"
)

// serve HTTP requests
type Server struct {
	store  *db.Store
	router *gin.Engine
}

// creates a new HTTP server and setup routing
func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/users", server.createUser)
	router.GET("/users/:username", server.getUserByUsername)
	router.GET("/posts", server.listPosts)

	server.router = router
	return server
}

// start the http server on address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
