package gapi

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	db "pxsemic.com/simplebank/db/sqlc"
	"pxsemic.com/simplebank/pb"
)

func convertUser(user db.User) *pb.User {
	return &pb.User{
		Username:          user.Username,
		Email:             user.Email,
		FullName:          user.FullName,
		CreatedAt:         timestamppb.New(user.CreatedAt),
		PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
	}
}
