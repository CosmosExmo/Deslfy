package gapi

import (
	"context"
	db "desly/db/sqlc"
	"desly/pb"
	"desly/token"
	"desly/val"
	"time"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUserToken(ctx context.Context, req *pb.CreateUserTokenRequest) (*pb.CreateUserTokenResponse, error) {
	authPayload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	violations := validateCreateUserTokenRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	expireAt := req.GetExpireAt().AsTime()
	duration := time.Until(expireAt)
	accessToken, _, err := server.tokenMaker.CreateToken(authPayload.Username, duration, token.AccessToken)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create access token: %s", err)
	}

	arg := db.CreateUserTokenParams{
		Owner:    authPayload.Username,
		Token:    accessToken,
		ExpireAt: expireAt,
	}

	token, err := server.store.CreateUserToken(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create user token: %s", err)
	}

	rsp := &pb.CreateUserTokenResponse{
		UserToken: convertUserToken(&token),
	}

	return rsp, nil
}

func validateCreateUserTokenRequest(req *pb.CreateUserTokenRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := req.GetExpireAt().CheckValid(); err != nil {
		violations = append(violations, fieldViolation("expire_at", err))
		//there is no point to check other validations if this true
		return
	}

	if err := val.ValidateExpireAt(req.GetExpireAt().AsTime()); err != nil {
		violations = append(violations, fieldViolation("expire_at", err))
	}

	return
}
