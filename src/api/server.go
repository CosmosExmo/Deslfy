package api

import (
	db "desly/db/sqlc"
	"desly/token"
	"desly/util"
	"fmt"

	"github.com/gin-gonic/gin"
)

type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{config: config, store: store, tokenMaker: tokenMaker}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.GET("/r/:desly", server.redirect)

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.POST("/desly", server.createDesly)
	authRoutes.GET("/desly/:desly", server.getDesly)

	authRoutes.POST("/token", server.createUserToken)
	authRoutes.POST("/token/delete", server.deleteUserToken)
	authRoutes.GET("/token/:id", server.getUserToken)
	authRoutes.GET("/token", server.getUserTokens)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
