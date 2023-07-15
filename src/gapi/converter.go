package gapi

import (
	db "desly/db/sqlc"
	"desly/pb"
	"fmt"

	"google.golang.org/protobuf/types/known/timestamppb"
)

const deslyRedirectUrl = "deslfy.com/r/"

func convertUser(user *db.User) *pb.User {
	return &pb.User{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
		CreatedAt:         timestamppb.New(user.CreatedAt),
	}
}

func convertUserToken(userToken *db.UserToken) *pb.UserToken {
	return &pb.UserToken{
		Id:        userToken.ID,
		Owner:     userToken.Owner,
		Token:     userToken.Token,
		ExpireAt:  timestamppb.New(userToken.ExpireAt),
		CreatedAt: timestamppb.New(userToken.CreatedAt),
	}
}

func convertUserTokens(userTokens *[]db.UserToken) []*pb.UserToken {
	destinationList := make([]*pb.UserToken, len(*userTokens))
	for i, source := range *userTokens {
		destinationList[i] = &pb.UserToken{
			Id:        source.ID,
			Owner:     source.Owner,
			Token:     source.Token,
			ExpireAt:  timestamppb.New(source.ExpireAt),
			CreatedAt: timestamppb.New(source.CreatedAt),
		}
	}
	return destinationList
}

func convertDesly(desly *db.Desly) *pb.Desly {
	return &pb.Desly{
		Id:        desly.ID,
		Redirect:  desly.Redirect,
		Desly:     desly.Desly,
		DeslyUrl:  fmt.Sprintf("%s%s", deslyRedirectUrl, desly.Desly),
		Clicked:   desly.Clicked,
		CreatedAt: timestamppb.New(desly.CreatedAt),
		Owner:     desly.Owner,
	}
}
