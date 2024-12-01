package api

import (
	"fmt"

	db "github.com/daniel-adam-ce/go-bank/db/sqlc"
	"github.com/daniel-adam-ce/go-bank/token"
	"github.com/daniel-adam-ce/go-bank/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {

	// this can be interchanged with NewJWTMaker
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("ccanot create token maker: %v", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	server.setupRouter()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	return server, nil
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func (server *Server) setupRouter() {

	router := gin.Default()
	// users routes
	router.POST("/users", server.createUser)
	// this is not RESTful
	router.POST("/users/login", server.loginUser)

	// session routes
	// this is not restful
	router.POST("/tokens/refresh", server.renewAccessToken)

	// add auth middleware
	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	// accounts routes
	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.GET("/accounts", server.listAccount)

	// tranfers routes
	authRoutes.POST("/transfers", server.createTransfer)

	server.router = router
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
