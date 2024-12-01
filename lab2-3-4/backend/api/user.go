package api

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"m1thrandir225/lab-2-3-4/auth"
	db "m1thrandir225/lab-2-3-4/db/sqlc"
	"m1thrandir225/lab-2-3-4/mail"
	"m1thrandir225/lab-2-3-4/util"
	"net/http"
	"strings"
)

type LoginRequest struct {
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type RegisterRequest struct {
	Email        string `form:"email" json:"email" binding:"required"`
	Password     string `form:"password" json:"password" binding:"required"`
	TwoFAEnabled bool   `form:"2fa_enabled" json:"2fa_enabled" binding:"required"`
}

type VerifyOTPRequest struct {
	UserID  int64  `form:"user_id" json:"user_id" binding:"required"`
	OTPCode string `form:"otp_code" json:"otp_code" binding:"required"`
}

type TokenPairResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (server *Server) Login(ctx *gin.Context) {
	var loginRequest LoginRequest

	if err := ctx.ShouldBindJSON(&loginRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUserByEmail(ctx, loginRequest.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = util.VerifyPassword(user.PasswordHash, loginRequest.Password)
	if err != nil {
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}

	if user.Is2faEnabled {
		currentCounterInt, err := server.store.GetCurrentCounter(ctx, user.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		currentCounter := uint64(currentCounterInt)

		otpCode, err := auth.GenerateHOTP(user.OtpSecret, currentCounter)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		otpContent := mail.GenerateOtpEmail(otpCode)

		err = server.mailService.SendMail(
			"theteam@sebastijanzindl.me",
			user.Email,
			"Your OTP Verification Code",
			otpContent,
		)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message":      "2FA required",
			"requires_2fa": true,
			"user_id":      user.ID,
			"otp_code":     otpCode,
		})
		return
	}
	accessToken, err := server.tokenMaker.GenerateToken(user.Email, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, err := server.tokenMaker.GenerateToken(user.Email, server.config.RefreshTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, TokenPairResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func (server *Server) Register(ctx *gin.Context) {
	var registerRequest RegisterRequest

	if err := ctx.ShouldBindJSON(&registerRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(registerRequest.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	otpSecret, err := auth.GenerateOTPSecret()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	args := db.CreateUserParams{
		Email:        registerRequest.Email,
		PasswordHash: hashedPassword,
		OtpSecret:    otpSecret,
		Is2faEnabled: registerRequest.TwoFAEnabled,
	}

	user, err := server.store.CreateUser(ctx, args)
	if err != nil {
		if isDuplicateKeyError(err) {
			ctx.JSON(http.StatusConflict, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateHotpCounterParams{
		UserID:  user.ID,
		Counter: 0,
	}

	err = server.store.CreateHotpCounter(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	welcomeContent := mail.GenerateWelcomeEmail(user.Email)

	err = server.mailService.SendMail(
		"theteam@sebastijanzindl.me",
		user.Email,
		"Welcome to the team",
		welcomeContent,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"user_id":    user.ID,
		"otp_secret": otpSecret,
		"email":      user.Email,
		"message":    "Registration successful. Please save your OTP secret securely.",
	})

}

func (server *Server) VerifyOTP(ctx *gin.Context) {
	var verifyOTPRequest VerifyOTPRequest

	if err := ctx.ShouldBindJSON(&verifyOTPRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUserById(ctx, verifyOTPRequest.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	isValid, err := server.otpService.GenerateAndVerifyOTP(
		ctx,
		user.ID,
		verifyOTPRequest.OTPCode,
		user.OtpSecret,
	)

	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if !isValid {
		log.Println(err)
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, err := server.tokenMaker.GenerateToken(user.Email, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, err := server.tokenMaker.GenerateToken(user.Email, server.config.RefreshTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, TokenPairResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})

}

func isDuplicateKeyError(err error) bool {
	// This will depend on your specific database and error handling
	// For SQLite, you might check for a unique constraint violation
	return strings.Contains(err.Error(), "UNIQUE constraint failed") ||
		strings.Contains(err.Error(), "unique constraint") ||
		strings.Contains(err.Error(), "duplicate key")
}
