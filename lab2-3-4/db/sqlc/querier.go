// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"context"
)

type Querier interface {
	AddUserToOrganization(ctx context.Context, arg AddUserToOrganizationParams) (UserOrganization, error)
	CleanupExpiredCounters(ctx context.Context) error
	CreateAccessRequest(ctx context.Context, arg CreateAccessRequestParams) (AccessRequest, error)
	CreateHotpCounter(ctx context.Context, arg CreateHotpCounterParams) error
	CreateInitialRoles(ctx context.Context, orgID int64) error
	CreateOrganization(ctx context.Context, name string) (Organization, error)
	CreateResource(ctx context.Context, arg CreateResourceParams) (Resource, error)
	CreateRole(ctx context.Context, arg CreateRoleParams) (Role, error)
	CreateRolePermission(ctx context.Context, arg CreateRolePermissionParams) (RolePermission, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteOrganization(ctx context.Context, id int64) error
	DeleteResource(ctx context.Context, id int64) error
	DeleteUser(ctx context.Context, id int64) error
	GetAccessRequest(ctx context.Context, id int64) (AccessRequest, error)
	GetActiveAccessRequest(ctx context.Context, arg GetActiveAccessRequestParams) (AccessRequest, error)
	GetAdminRole(ctx context.Context, orgID int64) (int64, error)
	GetCurrentCounter(ctx context.Context, userID int64) (int64, error)
	GetModeratorRole(ctx context.Context, orgID int64) (int64, error)
	GetOrganization(ctx context.Context, id int64) (Organization, error)
	GetResource(ctx context.Context, id int64) (Resource, error)
	GetRole(ctx context.Context, id int64) (Role, error)
	GetRolePermissions(ctx context.Context, arg GetRolePermissionsParams) (RolePermission, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserById(ctx context.Context, id int64) (User, error)
	GetUserOrganization(ctx context.Context, arg GetUserOrganizationParams) (UserOrganization, error)
	GetUserRole(ctx context.Context, arg GetUserRoleParams) (string, error)
	GetUserRoleId(ctx context.Context, orgID int64) (int64, error)
	IncreaseCounter(ctx context.Context, userID int64) (int64, error)
	ListActiveUserAccess(ctx context.Context, arg ListActiveUserAccessParams) ([]ListActiveUserAccessRow, error)
	ListOrganizationResources(ctx context.Context, orgID int64) ([]Resource, error)
	ListOrganizationRoles(ctx context.Context, orgID int64) ([]Role, error)
	ListPendingAccessRequests(ctx context.Context, orgID int64) ([]ListPendingAccessRequestsRow, error)
	ListUserAccessRequests(ctx context.Context, userID int64) ([]AccessRequest, error)
	ListUserOrganizations(ctx context.Context, userID int64) ([]UserOrganization, error)
	RemoveUserFromOrganization(ctx context.Context, arg RemoveUserFromOrganizationParams) error
	RevokeExpiredAccess(ctx context.Context) error
	UpdateAccessRequestStatus(ctx context.Context, arg UpdateAccessRequestStatusParams) error
	UpdateResource(ctx context.Context, arg UpdateResourceParams) (Resource, error)
	UpdateRole(ctx context.Context, arg UpdateRoleParams) (Role, error)
	UpdateRolePermissions(ctx context.Context, arg UpdateRolePermissionsParams) (RolePermission, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
	UpdateUserOrganizationRole(ctx context.Context, arg UpdateUserOrganizationRoleParams) error
	UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) error
}

var _ Querier = (*Queries)(nil)
