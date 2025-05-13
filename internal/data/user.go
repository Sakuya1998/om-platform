package data

import (
	"context"
	"fmt"
	"time"

	"om-platform/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// 内存存储实现 (实际项目中应替换为数据库实现)
type userRepo struct {
	data *Data
	log  *log.Helper

	// 模拟数据存储
	users           map[string]*biz.User       // 用户ID -> 用户
	usersByUsername map[string]*biz.User       // 用户名 -> 用户
	roles           map[string]*biz.Role       // 角色ID -> 角色
	rolesByName     map[string]*biz.Role       // 角色名 -> 角色
	permissions     map[string]*biz.Permission // 权限名 -> 权限
	userRoles       map[string]map[string]bool // 用户ID -> 角色ID -> 是否有该角色
	rolePermissions map[string]map[string]bool // 角色ID -> 权限名 -> 是否有该权限
}

// NewUserRepo 创建用户仓储实例
func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data:            data,
		log:             log.NewHelper(logger),
		users:           make(map[string]*biz.User),
		usersByUsername: make(map[string]*biz.User),
	}
}

// NewRoleRepo 创建角色仓储实例
func NewRoleRepo(data *Data, logger log.Logger) biz.RoleRepo {
	return &roleRepo{
		data:        data,
		log:         log.NewHelper(logger),
		roles:       make(map[string]*biz.Role),
		rolesByName: make(map[string]*biz.Role),
	}
}

// NewPermissionRepo 创建权限仓储实例
func NewPermissionRepo(data *Data, logger log.Logger) biz.PermissionRepo {
	return &permissionRepo{
		data:        data,
		log:         log.NewHelper(logger),
		permissions: make(map[string]*biz.Permission),
	}
}

// NewUserRoleRepo 创建用户角色关联仓储实例
func NewUserRoleRepo(data *Data, logger log.Logger) biz.UserRoleRepo {
	return &userRoleRepo{
		data:      data,
		log:       log.NewHelper(logger),
		userRoles: make(map[string]map[string]bool),
	}
}

// NewRolePermissionRepo 创建角色权限关联仓储实例
func NewRolePermissionRepo(data *Data, logger log.Logger) biz.RolePermissionRepo {
	return &rolePermissionRepo{
		data:            data,
		log:             log.NewHelper(logger),
		rolePermissions: make(map[string]map[string]bool),
	}
}

// 用户仓储实现
func (r *userRepo) Create(ctx context.Context, user *biz.User) (*biz.User, error) {
	r.log.WithContext(ctx).Infof("创建用户: %v", user.Username)

	// 检查用户名是否已存在
	if _, ok := r.usersByUsername[user.Username]; ok {
		return nil, biz.ErrUserAlreadyExists
	}

	// 生成用户ID
	if user.ID == "" {
		user.ID = uuid.New().String()
	}

	// 设置创建和更新时间
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	// 存储用户
	r.users[user.ID] = user
	r.usersByUsername[user.Username] = user

	return user, nil
}

func (r *userRepo) Get(ctx context.Context, id string) (*biz.User, error) {
	r.log.WithContext(ctx).Infof("获取用户: %v", id)

	user, ok := r.users[id]
	if !ok {
		return nil, biz.ErrUserNotFound
	}

	return user, nil
}

func (r *userRepo) Update(ctx context.Context, user *biz.User) (*biz.User, error) {
	r.log.WithContext(ctx).Infof("更新用户: %v", user.ID)

	// 检查用户是否存在
	oldUser, ok := r.users[user.ID]
	if !ok {
		return nil, biz.ErrUserNotFound
	}

	// 如果用户名发生变化，需要更新usersByUsername映射
	if oldUser.Username != user.Username {
		// 检查新用户名是否已被占用
		if _, exists := r.usersByUsername[user.Username]; exists {
			return nil, biz.ErrUserAlreadyExists
		}

		// 删除旧映射，添加新映射
		delete(r.usersByUsername, oldUser.Username)
		r.usersByUsername[user.Username] = user
	}

	// 更新时间
	user.UpdatedAt = time.Now()

	// 更新用户
	r.users[user.ID] = user

	return user, nil
}

func (r *userRepo) Delete(ctx context.Context, id string) error {
	r.log.WithContext(ctx).Infof("删除用户: %v", id)

	// 检查用户是否存在
	user, ok := r.users[id]
	if !ok {
		return biz.ErrUserNotFound
	}

	// 删除用户
	delete(r.users, id)
	delete(r.usersByUsername, user.Username)

	return nil
}

func (r *userRepo) List(ctx context.Context, pageSize int32, pageToken string, filter string) ([]*biz.User, string, int32, error) {
	r.log.WithContext(ctx).Info("列出用户")

	// 简单实现，不考虑分页和过滤
	users := make([]*biz.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}

	return users, "", int32(len(users)), nil
}

func (r *userRepo) FindByUsername(ctx context.Context, username string) (*biz.User, error) {
	r.log.WithContext(ctx).Infof("通过用户名查找用户: %v", username)

	user, ok := r.usersByUsername[username]
	if !ok {
		return nil, biz.ErrUserNotFound
	}

	return user, nil
}

func (r *userRepo) VerifyPassword(ctx context.Context, userID string, password string) (bool, error) {
	r.log.WithContext(ctx).Infof("验证用户密码: %v", userID)

	// 获取用户
	user, ok := r.users[userID]
	if !ok {
		return false, biz.ErrUserNotFound
	}

	// 验证密码
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	return err == nil, nil
}

func (r *userRepo) UpdatePassword(ctx context.Context, userID string, newPassword string) error {
	r.log.WithContext(ctx).Infof("更新用户密码: %v", userID)

	// 获取用户
	user, ok := r.users[userID]
	if !ok {
		return biz.ErrUserNotFound
	}

	// 生成密码哈希
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("生成密码哈希失败: %w", err)
	}

	// 更新密码哈希
	user.PasswordHash = string(hashedPassword)
	user.UpdatedAt = time.Now()

	return nil
}

// 角色仓储实现
type roleRepo struct {
	data *Data
	log  *log.Helper

	// 模拟数据存储
	roles       map[string]*biz.Role // 角色ID -> 角色
	rolesByName map[string]*biz.Role // 角色名 -> 角色
}

func (r *roleRepo) Create(ctx context.Context, role *biz.Role) (*biz.Role, error) {
	r.log.WithContext(ctx).Infof("创建角色: %v", role.Name)

	// 检查角色名是否已存在
	if _, ok := r.rolesByName[role.Name]; ok {
		return nil, biz.ErrRoleAlreadyExists
	}

	// 生成角色ID
	if role.ID == "" {
		role.ID = uuid.New().String()
	}

	// 设置创建和更新时间
	now := time.Now()
	role.CreatedAt = now
	role.UpdatedAt = now

	// 存储角色
	r.roles[role.ID] = role
	r.rolesByName[role.Name] = role

	return role, nil
}

func (r *roleRepo) Get(ctx context.Context, id string) (*biz.Role, error) {
	r.log.WithContext(ctx).Infof("获取角色: %v", id)

	role, ok := r.roles[id]
	if !ok {
		return nil, biz.ErrRoleNotFound
	}

	return role, nil
}

func (r *roleRepo) Update(ctx context.Context, role *biz.Role) (*biz.Role, error) {
	r.log.WithContext(ctx).Infof("更新角色: %v", role.ID)

	// 检查角色是否存在
	oldRole, ok := r.roles[role.ID]
	if !ok {
		return nil, biz.ErrRoleNotFound
	}

	// 如果角色名发生变化，需要更新rolesByName映射
	if oldRole.Name != role.Name {
		// 检查新角色名是否已被占用
		if _, exists := r.rolesByName[role.Name]; exists {
			return nil, biz.ErrRoleAlreadyExists
		}

		// 删除旧映射，添加新映射
		delete(r.rolesByName, oldRole.Name)
		r.rolesByName[role.Name] = role
	}

	// 更新时间
	role.UpdatedAt = time.Now()

	// 更新角色
	r.roles[role.ID] = role

	return role, nil
}

func (r *roleRepo) Delete(ctx context.Context, id string) error {
	r.log.WithContext(ctx).Infof("删除角色: %v", id)

	// 检查角色是否存在
	role, ok := r.roles[id]
	if !ok {
		return biz.ErrRoleNotFound
	}

	// 删除角色
	delete(r.roles, id)
	delete(r.rolesByName, role.Name)

	return nil
}

func (r *roleRepo) List(ctx context.Context, pageSize int32, pageToken string) ([]*biz.Role, string, int32, error) {
	r.log.WithContext(ctx).Info("列出角色")

	// 简单实现，不考虑分页
	roles := make([]*biz.Role, 0, len(r.roles))
	for _, role := range r.roles {
		roles = append(roles, role)
	}

	return roles, "", int32(len(roles)), nil
}

func (r *roleRepo) FindByName(ctx context.Context, name string) (*biz.Role, error) {
	r.log.WithContext(ctx).Infof("通过名称查找角色: %v", name)

	role, ok := r.rolesByName[name]
	if !ok {
		return nil, biz.ErrRoleNotFound
	}

	return role, nil
}

// 权限仓储实现
type permissionRepo struct {
	data *Data
	log  *log.Helper

	// 模拟数据存储
	permissions map[string]*biz.Permission // 权限名 -> 权限
}

func (r *permissionRepo) Create(ctx context.Context, permission *biz.Permission) (*biz.Permission, error) {
	r.log.WithContext(ctx).Infof("创建权限: %v", permission.Name)

	// 检查权限是否已存在
	if _, ok := r.permissions[permission.Name]; ok {
		return nil, biz.ErrPermissionAlreadyExists
	}

	// 存储权限
	r.permissions[permission.Name] = permission

	return permission, nil
}

func (r *permissionRepo) Get(ctx context.Context, name string) (*biz.Permission, error) {
	r.log.WithContext(ctx).Infof("获取权限: %v", name)

	permission, ok := r.permissions[name]
	if !ok {
		return nil, biz.ErrPermissionNotFound
	}

	return permission, nil
}

func (r *permissionRepo) List(ctx context.Context, resourceType string) ([]*biz.Permission, error) {
	r.log.WithContext(ctx).Info("列出权限")

	permissions := make([]*biz.Permission, 0)
	for _, permission := range r.permissions {
		// 如果指定了资源类型，则过滤
		if resourceType != "" && permission.ResourceType != resourceType {
			continue
		}
		permissions = append(permissions, permission)
	}

	return permissions, nil
}

// 用户角色关联仓储实现
type userRoleRepo struct {
	data *Data
	log  *log.Helper

	// 模拟数据存储
	userRoles map[string]map[string]bool // 用户ID -> 角色ID -> 是否有该角色
}

func (r *userRoleRepo) AssignRoleToUser(ctx context.Context, userID string, roleID string) error {
	r.log.WithContext(ctx).Infof("为用户 %s 分配角色 %s", userID, roleID)

	// 初始化用户的角色映射
	if _, ok := r.userRoles[userID]; !ok {
		r.userRoles[userID] = make(map[string]bool)
	}

	// 分配角色
	r.userRoles[userID][roleID] = true

	return nil
}

func (r *userRoleRepo) RemoveRoleFromUser(ctx context.Context, userID string, roleID string) error {
	r.log.WithContext(ctx).Infof("从用户 %s 移除角色 %s", userID, roleID)

	// 检查用户是否有角色映射
	userRoles, ok := r.userRoles[userID]
	if !ok {
		return nil // 用户没有任何角色，视为成功
	}

	// 移除角色
	delete(userRoles, roleID)

	return nil
}

func (r *userRoleRepo) ListUserRoles(ctx context.Context, userID string) ([]*biz.Role, error) {
	r.log.WithContext(ctx).Infof("获取用户 %s 的角色列表", userID)

	// 获取用户的角色映射
	userRoles, ok := r.userRoles[userID]
	if !ok {
		return []*biz.Role{}, nil // 用户没有任何角色，返回空列表
	}

	// 简化实现，避免循环依赖
	// 在实际项目中，应该使用数据库查询来获取角色信息
	roles := make([]*biz.Role, 0, len(userRoles))
	for roleID := range userRoles {
		// 创建一个简单的角色对象，实际项目中应从数据库获取完整信息
		roles = append(roles, &biz.Role{
			ID:   roleID,
			Name: "Role-" + roleID,
		})
	}

	return roles, nil
}

func (r *userRoleRepo) ListRoleUsers(ctx context.Context, roleID string) ([]*biz.User, error) {
	r.log.WithContext(ctx).Infof("获取角色 %s 的用户列表", roleID)

	// 简化实现，避免循环依赖
	// 在实际项目中，应该使用数据库查询来获取用户信息
	users := make([]*biz.User, 0)
	for userID, roles := range r.userRoles {
		if roles[roleID] {
			// 创建一个简单的用户对象，实际项目中应从数据库获取完整信息
			users = append(users, &biz.User{
				ID:       userID,
				Username: "User-" + userID,
			})
		}
	}

	return users, nil
}

func (r *userRoleRepo) HasRole(ctx context.Context, userID string, roleID string) (bool, error) {
	r.log.WithContext(ctx).Infof("检查用户 %s 是否拥有角色 %s", userID, roleID)

	// 获取用户的角色映射
	userRoles, ok := r.userRoles[userID]
	if !ok {
		return false, nil // 用户没有任何角色
	}

	return userRoles[roleID], nil
}

// 角色权限关联仓储实现
type rolePermissionRepo struct {
	data *Data
	log  *log.Helper

	// 模拟数据存储
	rolePermissions map[string]map[string]bool // 角色ID -> 权限名 -> 是否有该权限
}

func (r *rolePermissionRepo) AddPermissionToRole(ctx context.Context, roleID string, permissionName string) error {
	r.log.WithContext(ctx).Infof("为角色 %s 添加权限 %s", roleID, permissionName)

	// 初始化角色的权限映射
	if _, ok := r.rolePermissions[roleID]; !ok {
		r.rolePermissions[roleID] = make(map[string]bool)
	}

	// 添加权限
	r.rolePermissions[roleID][permissionName] = true

	return nil
}

func (r *rolePermissionRepo) RemovePermissionFromRole(ctx context.Context, roleID string, permissionName string) error {
	r.log.WithContext(ctx).Infof("从角色 %s 移除权限 %s", roleID, permissionName)

	// 检查角色是否有权限映射
	rolePermissions, ok := r.rolePermissions[roleID]
	if !ok {
		return nil // 角色没有任何权限，视为成功
	}

	// 移除权限
	delete(rolePermissions, permissionName)

	return nil
}

func (r *rolePermissionRepo) ListRolePermissions(ctx context.Context, roleID string) ([]*biz.Permission, error) {
	r.log.WithContext(ctx).Infof("获取角色 %s 的权限列表", roleID)

	// 获取角色的权限映射
	rolePermissions, ok := r.rolePermissions[roleID]
	if !ok {
		return []*biz.Permission{}, nil // 角色没有任何权限，返回空列表
	}

	// 简化实现，避免循环依赖
	// 在实际项目中，应该使用数据库查询来获取权限信息
	permissions := make([]*biz.Permission, 0, len(rolePermissions))
	for permName := range rolePermissions {
		// 创建一个简单的权限对象，实际项目中应从数据库获取完整信息
		permissions = append(permissions, &biz.Permission{
			Name:        permName,
			Description: "Permission-" + permName,
		})
	}

	return permissions, nil
}

func (r *rolePermissionRepo) HasPermission(ctx context.Context, roleID string, permissionName string) (bool, error) {
	r.log.WithContext(ctx).Infof("检查角色 %s 是否拥有权限 %s", roleID, permissionName)

	// 获取角色的权限映射
	rolePermissions, ok := r.rolePermissions[roleID]
	if !ok {
		return false, nil // 角色没有任何权限
	}

	return rolePermissions[permissionName], nil
}

// 初始化Data结构
func init() {
	// 已在文件顶部定义了ProviderSet
}
