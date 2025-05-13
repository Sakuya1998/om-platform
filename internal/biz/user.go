package biz

import (
	"context"
	"errors"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
	"github.com/google/wire"
)

// 错误定义
var (
	ErrUserNotFound            = errors.New("user not found")
	ErrRoleNotFound            = errors.New("role not found")
	ErrPermissionNotFound      = errors.New("permission not found")
	ErrUserAlreadyExists       = errors.New("user already exists")
	ErrRoleAlreadyExists       = errors.New("role already exists")
	ErrPermissionAlreadyExists = errors.New("permission already exists")
	ErrInvalidCredentials      = errors.New("invalid credentials")
)

// User 用户实体
type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"password_hash,omitempty"`
	Email        string    `json:"email"`
	DisplayName  string    `json:"display_name"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Role 角色实体
type Role struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	Permissions []*Permission `json:"permissions,omitempty"`
}

// Permission 权限实体
type Permission struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	ResourceType string `json:"resource_type"`
	Action       string `json:"action"`
}

// UserRoleRelation 用户角色关联
type UserRoleRelation struct {
	UserID string `json:"user_id"`
	RoleID string `json:"role_id"`
}

// RolePermissionRelation 角色权限关联
type RolePermissionRelation struct {
	RoleID         string `json:"role_id"`
	PermissionName string `json:"permission_name"`
}

// UserRepo 用户仓储接口
type UserRepo interface {
	Create(ctx context.Context, user *User) (*User, error)
	Get(ctx context.Context, id string) (*User, error)
	Update(ctx context.Context, user *User) (*User, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, pageSize int32, pageToken string, filter string) ([]*User, string, int32, error)
	FindByUsername(ctx context.Context, username string) (*User, error)
	VerifyPassword(ctx context.Context, userID string, password string) (bool, error)
	UpdatePassword(ctx context.Context, userID string, newPasswordHash string) error
}

// RoleRepo 角色仓储接口
type RoleRepo interface {
	Create(ctx context.Context, role *Role) (*Role, error)
	Get(ctx context.Context, id string) (*Role, error)
	Update(ctx context.Context, role *Role) (*Role, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, pageSize int32, pageToken string) ([]*Role, string, int32, error)
	FindByName(ctx context.Context, name string) (*Role, error)
}

// PermissionRepo 权限仓储接口
type PermissionRepo interface {
	Create(ctx context.Context, permission *Permission) (*Permission, error)
	Get(ctx context.Context, name string) (*Permission, error)
	List(ctx context.Context, resourceType string) ([]*Permission, error)
}

// UserRoleRepo 用户角色关联仓储接口
type UserRoleRepo interface {
	AssignRoleToUser(ctx context.Context, userID string, roleID string) error
	RemoveRoleFromUser(ctx context.Context, userID string, roleID string) error
	ListUserRoles(ctx context.Context, userID string) ([]*Role, error)
	ListRoleUsers(ctx context.Context, roleID string) ([]*User, error)
	HasRole(ctx context.Context, userID string, roleID string) (bool, error)
}

// RolePermissionRepo 角色权限关联仓储接口
type RolePermissionRepo interface {
	AddPermissionToRole(ctx context.Context, roleID string, permissionName string) error
	RemovePermissionFromRole(ctx context.Context, roleID string, permissionName string) error
	ListRolePermissions(ctx context.Context, roleID string) ([]*Permission, error)
	HasPermission(ctx context.Context, roleID string, permissionName string) (bool, error)
}

// UserUsecase 用户业务逻辑
type UserUsecase struct {
	repo           UserRepo
	roleRepo       RoleRepo
	permissionRepo PermissionRepo
	userRoleRepo   UserRoleRepo
	rolePermRepo   RolePermissionRepo
	log            *log.Helper
}

// NewUserUsecase 创建用户业务逻辑实例
func NewUserUsecase(
	repo UserRepo,
	roleRepo RoleRepo,
	permissionRepo PermissionRepo,
	userRoleRepo UserRoleRepo,
	rolePermRepo RolePermissionRepo,
	logger log.Logger) *UserUsecase {
	return &UserUsecase{
		repo:           repo,
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
		userRoleRepo:   userRoleRepo,
		rolePermRepo:   rolePermRepo,
		log:            log.NewHelper(logger),
	}
}

// CreateUser 创建用户
func (uc *UserUsecase) CreateUser(ctx context.Context, u *User, password string) (*User, error) {
	uc.log.WithContext(ctx).Infof("创建用户: %s", u.Username)

	// 检查用户名是否已存在
	existingUser, err := uc.repo.FindByUsername(ctx, u.Username)
	if err == nil && existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	// 生成用户ID
	u.ID = uuid.New().String()
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	u.IsActive = true

	// 密码哈希处理在数据层实现
	return uc.repo.Create(ctx, u)
}

// GetUser 获取用户信息
func (uc *UserUsecase) GetUser(ctx context.Context, id string) (*User, error) {
	uc.log.WithContext(ctx).Infof("获取用户: %s", id)
	return uc.repo.Get(ctx, id)
}

// UpdateUser 更新用户信息
func (uc *UserUsecase) UpdateUser(ctx context.Context, u *User, password string) (*User, error) {
	uc.log.WithContext(ctx).Infof("更新用户: %s", u.ID)

	// 检查用户是否存在
	existingUser, err := uc.repo.Get(ctx, u.ID)
	if err != nil {
		return nil, err
	}

	// 更新用户信息
	u.UpdatedAt = time.Now()

	// 如果提供了新密码，则更新密码
	if password != "" {
		if err := uc.repo.UpdatePassword(ctx, u.ID, password); err != nil {
			return nil, err
		}
	}

	return uc.repo.Update(ctx, u)
}

// DeleteUser 删除用户
func (uc *UserUsecase) DeleteUser(ctx context.Context, id string) error {
	uc.log.WithContext(ctx).Infof("删除用户: %s", id)
	return uc.repo.Delete(ctx, id)
}

// ListUsers 列出用户
func (uc *UserUsecase) ListUsers(ctx context.Context, pageSize int32, pageToken string, filter string) ([]*User, string, int32, error) {
	uc.log.WithContext(ctx).Info("列出用户")
	return uc.repo.List(ctx, pageSize, pageToken, filter)
}

// CreateRole 创建角色
func (uc *UserUsecase) CreateRole(ctx context.Context, r *Role) (*Role, error) {
	uc.log.WithContext(ctx).Infof("创建角色: %s", r.Name)

	// 检查角色名是否已存在
	existingRole, err := uc.roleRepo.FindByName(ctx, r.Name)
	if err == nil && existingRole != nil {
		return nil, ErrRoleAlreadyExists
	}

	// 生成角色ID
	r.ID = uuid.New().String()
	r.CreatedAt = time.Now()
	r.UpdatedAt = time.Now()

	return uc.roleRepo.Create(ctx, r)
}

// GetRole 获取角色信息
func (uc *UserUsecase) GetRole(ctx context.Context, id string) (*Role, error) {
	uc.log.WithContext(ctx).Infof("获取角色: %s", id)
	return uc.roleRepo.Get(ctx, id)
}

// UpdateRole 更新角色信息
func (uc *UserUsecase) UpdateRole(ctx context.Context, r *Role) (*Role, error) {
	uc.log.WithContext(ctx).Infof("更新角色: %s", r.ID)

	// 检查角色是否存在
	existingRole, err := uc.roleRepo.Get(ctx, r.ID)
	if err != nil {
		return nil, err
	}

	// 更新角色信息
	r.UpdatedAt = time.Now()

	return uc.roleRepo.Update(ctx, r)
}

// DeleteRole 删除角色
func (uc *UserUsecase) DeleteRole(ctx context.Context, id string) error {
	uc.log.WithContext(ctx).Infof("删除角色: %s", id)
	return uc.roleRepo.Delete(ctx, id)
}

// ListRoles 列出角色
func (uc *UserUsecase) ListRoles(ctx context.Context, pageSize int32, pageToken string) ([]*Role, string, int32, error) {
	uc.log.WithContext(ctx).Info("列出角色")
	return uc.roleRepo.List(ctx, pageSize, pageToken)
}

// CreatePermission 创建权限
func (uc *UserUsecase) CreatePermission(ctx context.Context, p *Permission) (*Permission, error) {
	uc.log.WithContext(ctx).Infof("创建权限: %s", p.Name)

	// 检查权限是否已存在
	existingPerm, err := uc.permissionRepo.Get(ctx, p.Name)
	if err == nil && existingPerm != nil {
		return nil, ErrPermissionAlreadyExists
	}

	return uc.permissionRepo.Create(ctx, p)
}

// ListPermissions 列出权限
func (uc *UserUsecase) ListPermissions(ctx context.Context, resourceType string) ([]*Permission, error) {
	uc.log.WithContext(ctx).Info("列出权限")
	return uc.permissionRepo.List(ctx, resourceType)
}

// AssignRoleToUser 为用户分配角色
func (uc *UserUsecase) AssignRoleToUser(ctx context.Context, userID string, roleID string) error {
	uc.log.WithContext(ctx).Infof("为用户 %s 分配角色 %s", userID, roleID)

	// 检查用户是否存在
	_, err := uc.repo.Get(ctx, userID)
	if err != nil {
		return err
	}

	// 检查角色是否存在
	_, err = uc.roleRepo.Get(ctx, roleID)
	if err != nil {
		return err
	}

	// 检查是否已分配
	hasRole, err := uc.userRoleRepo.HasRole(ctx, userID, roleID)
	if err != nil {
		return err
	}
	if hasRole {
		// 已分配，视为成功
		return nil
	}

	return uc.userRoleRepo.AssignRoleToUser(ctx, userID, roleID)
}

// RemoveRoleFromUser 从用户移除角色
func (uc *UserUsecase) RemoveRoleFromUser(ctx context.Context, userID string, roleID string) error {
	uc.log.WithContext(ctx).Infof("从用户 %s 移除角色 %s", userID, roleID)
	return uc.userRoleRepo.RemoveRoleFromUser(ctx, userID, roleID)
}

// ListUserRoles 获取用户角色列表
func (uc *UserUsecase) ListUserRoles(ctx context.Context, userID string) ([]*Role, error) {
	uc.log.WithContext(ctx).Infof("获取用户 %s 的角色列表", userID)
	return uc.userRoleRepo.ListUserRoles(ctx, userID)
}

// AddPermissionToRole 为角色添加权限
func (uc *UserUsecase) AddPermissionToRole(ctx context.Context, roleID string, permissionName string) error {
	uc.log.WithContext(ctx).Infof("为角色 %s 添加权限 %s", roleID, permissionName)

	// 检查角色是否存在
	_, err := uc.roleRepo.Get(ctx, roleID)
	if err != nil {
		return err
	}

	// 检查权限是否存在
	_, err = uc.permissionRepo.Get(ctx, permissionName)
	if err != nil {
		return err
	}

	// 检查是否已添加
	hasPermission, err := uc.rolePermRepo.HasPermission(ctx, roleID, permissionName)
	if err != nil {
		return err
	}
	if hasPermission {
		// 已添加，视为成功
		return nil
	}

	return uc.rolePermRepo.AddPermissionToRole(ctx, roleID, permissionName)
}

// RemovePermissionFromRole 从角色移除权限
func (uc *UserUsecase) RemovePermissionFromRole(ctx context.Context, roleID string, permissionName string) error {
	uc.log.WithContext(ctx).Infof("从角色 %s 移除权限 %s", roleID, permissionName)
	return uc.rolePermRepo.RemovePermissionFromRole(ctx, roleID, permissionName)
}

// ListRolePermissions 列出角色权限
func (uc *UserUsecase) ListRolePermissions(ctx context.Context, roleID string) ([]*Permission, error) {
	uc.log.WithContext(ctx).Infof("列出角色 %s 的权限", roleID)
	return uc.rolePermRepo.ListRolePermissions(ctx, roleID)
}

// 业务逻辑层的Provider集合更新
func init() {
	// 更新ProviderSet，添加用户管理相关的业务逻辑
	ProviderSet = wire.NewSet(
		NewGreeterUsecase,
		NewUserUsecase,
	)
}
