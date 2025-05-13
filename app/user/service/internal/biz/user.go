package biz

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	// ErrUserNotFound 用户未找到错误
	ErrUserNotFound = errors.NotFound("USER_NOT_FOUND", "用户未找到")
	// ErrUserAlreadyExists 用户已存在错误
	ErrUserAlreadyExists = errors.Conflict("USER_ALREADY_EXISTS", "用户已存在")
	// ErrRoleNotFound 角色未找到错误
	ErrRoleNotFound = errors.NotFound("ROLE_NOT_FOUND", "角色未找到")
	// ErrPermissionNotFound 权限未找到错误
	ErrPermissionNotFound = errors.NotFound("PERMISSION_NOT_FOUND", "权限未找到")
	// ErrInvalidCredentials 凭证无效错误
	ErrInvalidCredentials = errors.Unauthorized("INVALID_CREDENTIALS", "用户名或密码错误")
)

// User 用户领域模型
type User struct {
	ID           string
	Username     string
	PasswordHash string
	Email        string
	DisplayName  string
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Role 角色领域模型
type Role struct {
	ID          string
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Permission 权限领域模型
type Permission struct {
	Name         string
	Description  string
	ResourceType string
	Action       string
}

// UserRepo 用户仓库接口
type UserRepo interface {
	Create(ctx context.Context, user *User) (*User, error)
	Update(ctx context.Context, user *User) (*User, error)
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*User, error)
	FindByUsername(ctx context.Context, username string) (*User, error)
	List(ctx context.Context, pageSize int32, pageToken string, filter string) ([]*User, string, int32, error)
	AssignRole(ctx context.Context, userID, roleID string) error
	RemoveRole(ctx context.Context, userID, roleID string) error
	ListRoles(ctx context.Context, userID string) ([]*Role, error)
}

// RoleRepo 角色仓库接口
type RoleRepo interface {
	Create(ctx context.Context, role *Role) (*Role, error)
	Update(ctx context.Context, role *Role) (*Role, error)
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*Role, error)
	List(ctx context.Context, pageSize int32, pageToken string) ([]*Role, string, int32, error)
	AddPermission(ctx context.Context, roleID, permissionName string) error
	RemovePermission(ctx context.Context, roleID, permissionName string) error
	ListPermissions(ctx context.Context, roleID string) ([]*Permission, error)
}

// PermissionRepo 权限仓库接口
type PermissionRepo interface {
	Create(ctx context.Context, permission *Permission) (*Permission, error)
	FindByName(ctx context.Context, name string) (*Permission, error)
	List(ctx context.Context, resourceType string) ([]*Permission, error)
}

// UserUsecase 用户用例
type UserUsecase struct {
	userRepo       UserRepo
	roleRepo       RoleRepo
	permissionRepo PermissionRepo
	log            *log.Helper
}

// NewUserUsecase 创建用户用例
func NewUserUsecase(ur UserRepo, rr RoleRepo, pr PermissionRepo, logger log.Logger) *UserUsecase {
	return &UserUsecase{
		userRepo:       ur,
		roleRepo:       rr,
		permissionRepo: pr,
		log:            log.NewHelper(logger),
	}
}

// CreateUser 创建用户
func (uc *UserUsecase) CreateUser(ctx context.Context, u *User, password string) (*User, error) {
	// 检查用户名是否已存在
	existingUser, err := uc.userRepo.FindByUsername(ctx, u.Username)
	if err == nil && existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	// 实际应用中需要对密码进行加密处理
	// 这里简化处理，实际项目中应使用bcrypt等算法
	u.PasswordHash = password + "_hashed" // 示例，实际应使用安全的哈希算法
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	u.IsActive = true

	return uc.userRepo.Create(ctx, u)
}

// GetUser 获取用户信息
func (uc *UserUsecase) GetUser(ctx context.Context, id string) (*User, error) {
	return uc.userRepo.FindByID(ctx, id)
}

// UpdateUser 更新用户信息
func (uc *UserUsecase) UpdateUser(ctx context.Context, u *User, password string) (*User, error) {
	// 检查用户是否存在
	existingUser, err := uc.userRepo.FindByID(ctx, u.ID)
	if err != nil {
		return nil, err
	}
	if existingUser == nil {
		return nil, ErrUserNotFound
	}

	// 如果提供了新密码，则更新密码
	if password != "" {
		// 实际应用中需要对密码进行加密处理
		u.PasswordHash = password + "_hashed" // 示例，实际应使用安全的哈希算法
	} else {
		u.PasswordHash = existingUser.PasswordHash
	}

	u.UpdatedAt = time.Now()
	return uc.userRepo.Update(ctx, u)
}

// DeleteUser 删除用户
func (uc *UserUsecase) DeleteUser(ctx context.Context, id string) error {
	// 检查用户是否存在
	existingUser, err := uc.userRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return ErrUserNotFound
	}

	return uc.userRepo.Delete(ctx, id)
}

// ListUsers 列出用户
func (uc *UserUsecase) ListUsers(ctx context.Context, pageSize int32, pageToken string, filter string) ([]*User, string, int32, error) {
	return uc.userRepo.List(ctx, pageSize, pageToken, filter)
}

// AssignRoleToUser 为用户分配角色
func (uc *UserUsecase) AssignRoleToUser(ctx context.Context, userID, roleID string) error {
	// 检查用户是否存在
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	// 检查角色是否存在
	role, err := uc.roleRepo.FindByID(ctx, roleID)
	if err != nil {
		return err
	}
	if role == nil {
		return ErrRoleNotFound
	}

	return uc.userRepo.AssignRole(ctx, userID, roleID)
}

// RemoveRoleFromUser 从用户移除角色
func (uc *UserUsecase) RemoveRoleFromUser(ctx context.Context, userID, roleID string) error {
	// 检查用户是否存在
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	// 检查角色是否存在
	role, err := uc.roleRepo.FindByID(ctx, roleID)
	if err != nil {
		return err
	}
	if role == nil {
		return ErrRoleNotFound
	}

	return uc.userRepo.RemoveRole(ctx, userID, roleID)
}

// ListUserRoles 获取用户角色列表
func (uc *UserUsecase) ListUserRoles(ctx context.Context, userID string) ([]*Role, error) {
	// 检查用户是否存在
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	return uc.userRepo.ListRoles(ctx, userID)
}

// CreateRole 创建角色
func (uc *UserUsecase) CreateRole(ctx context.Context, r *Role) (*Role, error) {
	r.CreatedAt = time.Now()
	r.UpdatedAt = time.Now()
	return uc.roleRepo.Create(ctx, r)
}

// GetRole 获取角色信息
func (uc *UserUsecase) GetRole(ctx context.Context, id string) (*Role, error) {
	return uc.roleRepo.FindByID(ctx, id)
}

// UpdateRole 更新角色信息
func (uc *UserUsecase) UpdateRole(ctx context.Context, r *Role) (*Role, error) {
	// 检查角色是否存在
	existingRole, err := uc.roleRepo.FindByID(ctx, r.ID)
	if err != nil {
		return nil, err
	}
	if existingRole == nil {
		return nil, ErrRoleNotFound
	}

	r.UpdatedAt = time.Now()
	return uc.roleRepo.Update(ctx, r)
}

// DeleteRole 删除角色
func (uc *UserUsecase) DeleteRole(ctx context.Context, id string) error {
	// 检查角色是否存在
	existingRole, err := uc.roleRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existingRole == nil {
		return ErrRoleNotFound
	}

	return uc.roleRepo.Delete(ctx, id)
}

// ListRoles 列出角色
func (uc *UserUsecase) ListRoles(ctx context.Context, pageSize int32, pageToken string) ([]*Role, string, int32, error) {
	return uc.roleRepo.List(ctx, pageSize, pageToken)
}

// AddPermissionToRole 为角色添加权限
func (uc *UserUsecase) AddPermissionToRole(ctx context.Context, roleID, permissionName string) error {
	// 检查角色是否存在
	role, err := uc.roleRepo.FindByID(ctx, roleID)
	if err != nil {
		return err
	}
	if role == nil {
		return ErrRoleNotFound
	}

	// 检查权限是否存在
	permission, err := uc.permissionRepo.FindByName(ctx, permissionName)
	if err != nil {
		return err
	}
	if permission == nil {
		return ErrPermissionNotFound
	}

	return uc.roleRepo.AddPermission(ctx, roleID, permissionName)
}

// RemovePermissionFromRole 从角色移除权限
func (uc *UserUsecase) RemovePermissionFromRole(ctx context.Context, roleID, permissionName string) error {
	// 检查角色是否存在
	role, err := uc.roleRepo.FindByID(ctx, roleID)
	if err != nil {
		return err
	}
	if role == nil {
		return ErrRoleNotFound
	}

	// 检查权限是否存在
	permission, err := uc.permissionRepo.FindByName(ctx, permissionName)
	if err != nil {
		return err
	}
	if permission == nil {
		return ErrPermissionNotFound
	}

	return uc.roleRepo.RemovePermission(ctx, roleID, permissionName)
}

// ListRolePermissions 列出角色权限
func (uc *UserUsecase) ListRolePermissions(ctx context.Context, roleID string) ([]*Permission, error) {
	// 检查角色是否存在
	role, err := uc.roleRepo.FindByID(ctx, roleID)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, ErrRoleNotFound
	}

	return uc.roleRepo.ListPermissions(ctx, roleID)
}

// CreatePermission 创建权限
func (uc *UserUsecase) CreatePermission(ctx context.Context, p *Permission) (*Permission, error) {
	return uc.permissionRepo.Create(ctx, p)
}

// ListPermissions 列出权限
func (uc *UserUsecase) ListPermissions(ctx context.Context, resourceType string) ([]*Permission, error) {
	return uc.permissionRepo.List(ctx, resourceType)
}
