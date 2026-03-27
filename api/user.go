package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "pxsemic.com/simplebank/db/sqlc"
	"pxsemic.com/simplebank/util"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,min=1,max=6"`
	Password string `json:"password" binding:"required,min=6,max=32"`
	FullName string `json:"full_name" binding:"required,min=1,max=32"`
	Email    string `json:"email" binding:"required,email"`
}

type createUserResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	CreateAt          time.Time `json:"create_at"`
	PasswordChangedAt time.Time `json:"password_create_at"`
}

// create user
func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	arg := db.CreateUserParams{
		Username:       req.Username,
		FullName:       req.FullName,
		HashedPassword: hashedPassword,
		Email:          req.Email,
	}
	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case db.UniqueViolation, db.ForeignKeyViolation:
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := createUserResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		CreateAt:          user.CreatedAt,
		PasswordChangedAt: user.PasswordChangedAt,
	}
	ctx.JSON(http.StatusOK, rsp)
}
