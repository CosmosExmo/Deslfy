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

func (server *Server) DeleteUserToken(ctx context.Context, req *pb.DeleteUserTokenRequest) (*pb.DeleteUserTokenResponse, error) {
	authPayload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	violations := validateDeleteUserTokenRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	arg := db.DeleteUserTokenParams{
		Owner: authPayload.Username,
		ID: req.GetId(),
	}
	err = server.store.DeleteUserToken(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot delete user token: %s", err)
	}

	rsp := &pb.DeleteUserTokenResponse{
		IsDeleteSuccessful: true,
	}

	return rsp, nil
}

func validateDeleteUserTokenRequest(req *pb.DeleteUserTokenRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUserTokenId(req.GetId()); err != nil {
		violations = append(violations, fieldViolation("id", err))
	}

	return
}
