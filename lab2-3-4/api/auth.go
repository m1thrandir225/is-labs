package api

import (
	"database/sql"
	"errors"
	"log"
	"m1thrandir225/lab-2-3-4/auth"
	db "m1thrandir225/lab-2-3-4/db/sqlc"
	"m1thrandir225/lab-2-3-4/mail"
	"m1thrandir225/lab-2-3-4/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

type refreshTokenRequest struct {
	AccessToken  string `form:"access_token" json:"access_token" binding:"required"`
	RefreshToken string `form:"refresh_token" json:"refresh_token" binding:"required"`
}
type loginRequest struct {
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type registerRequest struct {
	Email        string `form:"email" json:"email" binding:"required"`
	Password     string `form:"password" json:"password" binding:"required"`
	TwoFAEnabled bool   `form:"2fa_enabled" json:"2fa_enabled" binding:"required"`
}

type verifyOTPRequest struct {
	UserID  int64  `form:"user_id" json:"user_id" binding:"required"`
	OTPCode string `form:"otp_code" json:"otp_code" binding:"required"`
}

type tokenPairResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (server *Server) login(ctx *gin.Context) {
	var req loginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
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

	err = util.VerifyPassword(user.PasswordHash, req.Password)
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

	ctx.JSON(http.StatusOK, tokenPairResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func (server *Server) register(ctx *gin.Context) {
	var req registerRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
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
		Email:        req.Email,
		PasswordHash: hashedPassword,
		OtpSecret:    otpSecret,
		Is2faEnabled: req.TwoFAEnabled,
	}

	user, err := server.store.CreateUser(ctx, args)
	if err != nil {
		if util.IsDuplicateKeyError(err) {
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

func (server *Server) verifyOTP(ctx *gin.Context) {
	var req verifyOTPRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUserById(ctx, req.UserID)
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
		req.OTPCode,
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

	ctx.JSON(http.StatusOK, tokenPairResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})

}

func (server *Server) refreshToken(ctx *gin.Context) {
	var req refreshTokenRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	_, err := server.tokenMaker.ValidateToken(req.AccessToken)
	if err == nil {
		ctx.Status(http.StatusForbidden)
		return
	}

	claims, err := server.tokenMaker.ValidateToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	newToken, err := server.tokenMaker.GenerateToken(claims.Email, server.config.AccessTokenDuration)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access_token": newToken,
	})
}
