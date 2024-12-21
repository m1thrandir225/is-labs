package dto

import (
	"database/sql/driver"
	"fmt"
)

type Role string

const (
	RoleAdmin     Role = "admin"
	RoleUser      Role = "user"
	RoleModerator Role = "moderator"
)

func (r *Role) Scan(value interface{}) error {
	if value == nil {
		*r = RoleUser // Default value
		return nil
	}

	str, ok := value.(string)

	if !ok {
		return fmt.Errorf("expected string for Role, got %T", value)
	}
	if str != "admin" && str != "user" && str != "moderator" {
		return fmt.Errorf("invalid role value: %q", str)
	}

	*r = Role(str)
	return nil
}

func (r *Role) ToString() string {
	return string(*r)
}

func (r Role) Value() (driver.Value, error) {
	return string(r), nil
}

func (r Role) IsValid() bool {
	switch r {
	case RoleAdmin, RoleUser, RoleModerator:
		return true
	default:
		return false
	}
}
