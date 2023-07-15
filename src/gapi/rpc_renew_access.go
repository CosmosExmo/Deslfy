package gapi

import (
	"context"
	"database/sql"
	"desly/pb"
	"desly/token"
	"desly/util"
	"desly/val"
	"time"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) RenewAccess(ctx context.Context, req *pb.RenewAccessRequest) (*pb.RenewAccessResponse, error) {
	violations := validateRenewAccessRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	refreshPayload, err := server.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "cannot verify refresh token: %s", err)
	}

	session, err := server.store.GetSession(ctx, refreshPayload.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "cannot find users session: %s", err)
		}
		return nil, status.Errorf(codes.Internal, "cannot users session: %s", err)
	}

	if session.RefreshToken != req.RefreshToken {
		err := util.ErrorMismatchSessionToken
		return nil, status.Errorf(codes.Unauthenticated, "cannot authenticate token: %s", err)
	}

	if session.IsBlocked {
		err := util.ErrorBlockedSession
		return nil, status.Errorf(codes.Unauthenticated, "blocked session: %s", err)
	}

	if session.Username != refreshPayload.Username {
		err := util.ErrorIncorrectUser
		return nil, status.Errorf(codes.Unauthenticated, "unauthorized access: %s", err)
	}

	if time.Now().After(session.ExpiresAt) {
		err := util.ErrorExpiredSession
		return nil, status.Errorf(codes.Unauthenticated, "unauthorized access: %s", err)
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(refreshPayload.Username, server.config.AccessTokenDuration, token.AccessToken)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create access token: %s", err)
	}

	rsp := &pb.RenewAccessResponse{
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpiredAt),
	}

	return rsp, nil
}

func validateRenewAccessRequest(req *pb.RenewAccessRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateString(req.GetRefreshToken(), 100, 600); err != nil {
		violations = append(violations, fieldViolation("refresh_token", err))
	}

	return
}