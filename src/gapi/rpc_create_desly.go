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

func (server *Server) CreateDesly(ctx context.Context, req *pb.CreateDeslyRequest) (*pb.CreateDeslyResponse, error) {
	authPayload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	violations := validateCreateDeslyRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	arg := db.CreateDeslyParams{
		Redirect: req.GetRedirect(),
		Owner:    authPayload.Username,
	}

	desly, err := server.store.CreateDesly(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create desly: %s", err)
	}

	rsp := &pb.CreateDeslyResponse{
		Desly: convertDesly(&desly),
	}

	return rsp, nil
}

func validateCreateDeslyRequest(req *pb.CreateDeslyRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateRedirectUrl(req.GetRedirect()); err != nil {
		violations = append(violations, fieldViolation("redirect", err))
	}

	return
}
