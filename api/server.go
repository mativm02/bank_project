package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/mativm02/bank_system/db/sqlc"
)

// Server servers HTTP requests to the API.
type Server struct {
	store  db.Store    // It will allow us to access the database.
	router *gin.Engine // It will allow us to access the router.
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)
	server.router = router
	return server
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

// Start starts the HTTP server.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
