package gapi

import (
	"fmt"

	db "github.com/mativm02/bank_system/db/sqlc"
	"github.com/mativm02/bank_system/pb"
	"github.com/mativm02/bank_system/token"
	"github.com/mativm02/bank_system/util"
	"github.com/mativm02/bank_system/worker"
)

// Server servers HTTP requests to the API.
type Server struct {
	pb.UnimplementedSimpleBankServer
	store           db.Store    // It will allow us to access the database.
	tokenMaker      token.Maker // It will allow us to access the token maker.
	config          util.Config // It will allow us to access the configuration.
	taskDistributor worker.TaskDistributor
}

func NewServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		store:           store,
		tokenMaker:      tokenMaker,
		config:          config,
		taskDistributor: taskDistributor,
	}

	return server, nil
}

// Start starts the HTTP server.
func (server *Server) Start(address string) error {
	return nil
}
