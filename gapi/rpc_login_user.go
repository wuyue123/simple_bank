package gapi

import (
	"context"
	"errors"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	db "pxsemic.com/simplebank/db/sqlc"
	"pxsemic.com/simplebank/pb"
	"pxsemic.com/simplebank/token"
	"pxsemic.com/simplebank/util"
	"pxsemic.com/simplebank/val"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	violations := server.validateLoginUserRequest(req)
	if len(violations) > 0 {
		return nil, invalidArgumentError(violations)
	}
	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "user name already exist: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}
	if err = util.CheckPassword(req.Password, user.HashedPassword); err != nil {
		return nil, status.Errorf(codes.NotFound, "invalid password")
	}
	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
		user.Username,
		util.BankerRole,
		server.config.AccessTokenDuration,
		token.TokenTypeAccessToken,
	)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create access token: %v", err)
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(
		user.Username,
		util.BankerRole,
		server.config.RefreshTokenDuration,
		token.TokenTypeRefreshToken,
	)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create refresh token: %v", err)
	}
	mtdt := server.extractMetadata(ctx)
	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		UserAgent:    mtdt.UserAgent,
		ClientIp:     mtdt.ClientIP,
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
		RefreshToken: refreshToken,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create session: %v", err)
	}

	rsp := &pb.LoginUserResponse{
		SessionId:             session.ID.String(),
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpiredAt),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiredAt),
		User:                  convertUser(user),
	}
	return rsp, nil
}

func (server *Server) validateLoginUserRequest(req *pb.LoginUserRequest) []*errdetails.BadRequest_FieldViolation {
	var violations []*errdetails.BadRequest_FieldViolation

	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}
	if err := val.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}
	return violations
}
