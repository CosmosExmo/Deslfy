package gapi

import (
	"context"
	db "desly/db/sqlc"
	"desly/pb"
	"desly/val"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) GetUserToken(ctx context.Context, req *pb.GetUserTokenRequest) (*pb.GetUserTokenResponse, error) {
	authPayload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	violations := validateGetUserTokenRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	arg := db.GetUserTokenParams{
		Owner: authPayload.Username,
		ID: req.GetId(),
	}
	token, err := server.store.GetUserToken(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot get user token: %s", err)
	}

	rsp := &pb.GetUserTokenResponse{
		UserToken: convertUserToken(&token),
	}

	return rsp, nil
}

func validateGetUserTokenRequest(req *pb.GetUserTokenRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUserTokenId(req.GetId()); err != nil {
		violations = append(violations, fieldViolation("id", err))
	}

	return
}
