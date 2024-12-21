// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"database/sql"

	"m1thrandir225/lab-2-3-4/dto"
)

type HotpCounter struct {
	UserID            int64        `json:"user_id"`
	Counter           int64        `json:"counter"`
	LastUsedTimestamp sql.NullTime `json:"last_used_timestamp"`
}

type User struct {
	ID           int64    `json:"id"`
	Email        string   `json:"email"`
	PasswordHash string   `json:"password_hash"`
	OtpSecret    string   `json:"otp_secret"`
	Is2faEnabled bool     `json:"is_2fa_enabled"`
	Role         dto.Role `json:"role"`
}
