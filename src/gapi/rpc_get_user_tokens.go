package gapi

import (
	"context"
	"desly/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) GetUserTokens(ctx context.Context, req *pb.GetUserTokensRequest) (*pb.GetUserTokensResponse, error) {
	authPayload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	tokens, err := server.store.GetUserTokens(ctx, authPayload.Username)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot get user tokens: %s", err)
	}

	rsp := &pb.GetUserTokensResponse{
		UserTokens: convertUserTokens(&tokens),
	}

	return rsp, nil
}
