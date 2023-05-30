package api

import (
	db "desly/db/sqlc"

	"github.com/gin-gonic/gin"
)

type Server struct {
	store db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/api/desly", server.createDesly)
	router.GET("/api/desly/:desly", server.getDesly)
	router.GET("/r/:desly", server.redirect)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
} 