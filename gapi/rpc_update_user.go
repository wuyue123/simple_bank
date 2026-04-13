package gapi

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	db "pxsemic.com/simplebank/db/sqlc"
	"pxsemic.com/simplebank/pb"
	"pxsemic.com/simplebank/util"
	"pxsemic.com/simplebank/val"
)

func (server *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	// TODO add authorization check
	authPayload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	violations := server.validateUpdateUserRequest(req)
	if len(violations) > 0 {
		return nil, invalidArgumentError(violations)
	}

	if authPayload.Username != req.GetUsername() {
		return nil, status.Errorf(codes.PermissionDenied, "cannot update other user's info")
	}

	arg := db.UpdateUserParams{
		Username: req.GetUsername(),
		FullName: pgtype.Text{
			String: req.GetFullName(),
			Valid:  req.FullName != nil,
		},
		Email: pgtype.Text{
			String: req.GetEmail(),
			Valid:  req.Email != nil,
		},
	}
	if req.Password != nil {
		hashedPassword, err := util.HashPassword(req.GetPassword())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to hash password: %v", err)
		}
		arg.HashedPassword = pgtype.Text{
			String: hashedPassword,
			Valid:  true,
		}
		arg.PasswordChangedAt = pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		}
	}
	user, err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "user name already exist: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to update user: %v", err)
	}

	rsp := &pb.UpdateUserResponse{
		User: convertUser(user),
	}
	return rsp, nil
}

func (server *Server) validateUpdateUserRequest(req *pb.UpdateUserRequest) []*errdetails.BadRequest_FieldViolation {
	var violations []*errdetails.BadRequest_FieldViolation

	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}
	if req.FullName != nil {
		if err := val.ValidateFullName(req.GetFullName()); err != nil {
			violations = append(violations, fieldViolation("full_name", err))
		}
	}
	if req.Email != nil {
		if err := val.ValidateEmail(req.GetEmail()); err != nil {
			violations = append(violations, fieldViolation("email", err))
		}
	}
	if req.Password != nil {
		if err := val.ValidatePassword(req.GetPassword()); err != nil {
			violations = append(violations, fieldViolation("password", err))
		}
	}

	return violations
}
