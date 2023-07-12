package gapi

import (
	db "desly/db/sqlc"
	"desly/pb"
	"desly/token"
	"desly/util"
	"desly/worker"
	"fmt"
)

type Server struct {
	pb.UnimplementedDeslfyServer
	config          util.Config
	store           db.Store
	tokenMaker      token.Maker
	taskDistributor worker.TaskDistributor
}

func NewServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{config: config, store: store, tokenMaker: tokenMaker, taskDistributor: taskDistributor}

	return server, nil
}
