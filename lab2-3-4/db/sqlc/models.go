// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"database/sql"
	"time"
)

type AccessRequest struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id"`
	ResourceID int64     `json:"resource_id"`
	Status     string    `json:"status"`
	Reason     string    `json:"reason"`
	ExpiresAt  time.Time `json:"expires_at"`
	CreatedAt  time.Time `json:"created_at"`
}

type HotpCounter struct {
	UserID            int64        `json:"user_id"`
	Counter           int64        `json:"counter"`
	LastUsedTimestamp sql.NullTime `json:"last_used_timestamp"`
}

type Organization struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type Resource struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	OrgID     int64     `json:"org_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Role struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	OrgID     int64     `json:"org_id"`
	CreatedAt time.Time `json:"created_at"`
}

type RolePermission struct {
	ID         int64     `json:"id"`
	RoleID     int64     `json:"role_id"`
	ResourceID int64     `json:"resource_id"`
	CanRead    bool      `json:"can_read"`
	CanWrite   bool      `json:"can_write"`
	CanDelete  bool      `json:"can_delete"`
	CreatedAt  time.Time `json:"created_at"`
}

type User struct {
	ID           int64     `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	OtpSecret    string    `json:"otp_secret"`
	Is2faEnabled bool      `json:"is_2fa_enabled"`
	CreatedAt    time.Time `json:"created_at"`
}

type UserOrganization struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	OrgID     int64     `json:"org_id"`
	RoleID    int64     `json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
}
