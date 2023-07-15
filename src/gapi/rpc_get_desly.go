package gapi

import (
	"context"
	"database/sql"
	db "desly/db/sqlc"
	"desly/pb"
	"desly/val"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) GetDesly(ctx context.Context, req *pb.GetDeslyRequest) (*pb.GetDeslyResponse, error) {
	authPayload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	violations := validateGetDeslyRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	arg := db.GetDeslyParams{
		Desly: req.Desly,
		Owner: authPayload.Username,
	}

	desly, err := server.store.GetDesly(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "cannot find desly: %s", err)
		}

		return nil, status.Errorf(codes.Internal, "cannot get desly: %s", err)
	}

	rsp := &pb.GetDeslyResponse{
		Desly: convertDesly(&desly),
	}

	return rsp, nil
}

func validateGetDeslyRequest(req *pb.GetDeslyRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateString(req.GetDesly(), 6, 6); err != nil {
		violations = append(violations, fieldViolation("desly", err))
	}

	return
}
