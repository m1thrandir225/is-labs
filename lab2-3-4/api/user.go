package api

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"m1thrandir225/lab-2-3-4/auth"
	db "m1thrandir225/lab-2-3-4/db/sqlc"
	"m1thrandir225/lab-2-3-4/util"
	"net/http"
)

type LoginRequest struct {
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type RegisterRequest struct {
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type LoginResponse struct {
	LimitedAccessToken string `json:"limited_access_token"`
	OTP                string `json:"otp"`
	Email              string `json:"email"`
	Message            string `json:"message"`
}

type VerifyOTPRequest struct {
	OTP   string `form:"otp" json:"otp" binding:"required"`
	Email string `form:"email" json:"email" binding:"required"`
}

type VerifyOTPResponse struct {
	User         db.User `json:"user"`
	AccessToken  string  `json:"access_token"`
	RefreshToken string  `json:"refresh_token"`
}

func (server *Server) Login(ctx *gin.Context) {
	var req LoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	//1. Get user from db
	user, err := server.store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	//2. Check password hash
	err = util.VerifyPassword(user.PasswordHash, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	//3. Sync Counter
	server.authenticationManager.SyncCounterWithDatabase(user.Email, user.Counter)

	token, err := server.tokenMaker.GenerateToken(user.Email, false, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	otp, err := server.authenticationManager.GenerateAndIncrement(user.Email, []byte(user.SecretKey), 6)

	response := LoginResponse{
		OTP:                otp,
		Email:              user.Email,
		LimitedAccessToken: token,
		Message:            "Login successful, 2FA required",
	}

	ctx.JSON(http.StatusOK, response)
}

func (server *Server) Register(ctx *gin.Context) {
	var req RegisterRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	passwordHash, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	hotpKey := auth.GenerateHOTPKey()

	arg := db.CreateUserParams{
		Email:        req.Email,
		PasswordHash: passwordHash,
		SecretKey:    string(hotpKey),
		Counter:      int64(0),
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	//Sync counter with authentication manager
	server.authenticationManager.SyncCounterWithDatabase(user.Email, user.Counter)

	otp, err := server.authenticationManager.GenerateAndIncrement(user.Email, hotpKey, 6)

	token, err := server.tokenMaker.GenerateToken(user.Email, false, server.config.AccessTokenDuration)

	response := LoginResponse{
		OTP:                otp,
		Email:              user.Email,
		LimitedAccessToken: token,
		Message:            "Login successful, 2FA required",
	}

	ctx.JSON(http.StatusOK, response)
}

func (server *Server) VerifyOTP(ctx *gin.Context) {
	var req VerifyOTPRequest

	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	verified, newCounter := server.authenticationManager.VerifyAndIncrementCounter(
		req.Email,
		[]byte(user.SecretKey),
		req.OTP,
		6,
		5,
	)
	if !verified {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	arg := db.UpdateUserCounterParams{
		ID:      user.ID,
		Counter: int64(newCounter),
	}

	err = server.store.UpdateUserCounter(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	accessToken, err := server.tokenMaker.GenerateToken(req.Email, true, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, err := server.tokenMaker.GenerateToken(req.Email, true, server.config.RefreshTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := VerifyOTPResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	ctx.JSON(http.StatusOK, response)
}
