// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1

package db

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	CreateDesly(ctx context.Context, arg CreateDeslyParams) (Desly, error)
	CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	CreateUserToken(ctx context.Context, arg CreateUserTokenParams) (UserToken, error)
	CreateVerifyEmail(ctx context.Context, arg CreateVerifyEmailParams) (VerifyEmail, error)
	DeleteUserToken(ctx context.Context, arg DeleteUserTokenParams) error
	GetDesly(ctx context.Context, arg GetDeslyParams) (Desly, error)
	GetRedirectByDesly(ctx context.Context, desly string) (string, error)
	GetSession(ctx context.Context, id uuid.UUID) (Session, error)
	GetUser(ctx context.Context, username string) (User, error)
	GetUserToken(ctx context.Context, arg GetUserTokenParams) (UserToken, error)
	GetUserTokens(ctx context.Context, owner string) ([]UserToken, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
	UpdateVerifyEmail(ctx context.Context, arg UpdateVerifyEmailParams) (VerifyEmail, error)
}

var _ Querier = (*Queries)(nil)
