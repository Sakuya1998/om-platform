package data

import (
	"context"
	"fmt"
	"sync"
	"time"

	"om-platform/app/user/service/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
)

// 内存存储实现，实际项目中应替换为数据库实现
type userRepo struct {
	data *Data
	log  *log.Helper

	// 内存存储
	users           map[string]*biz.User
	roles           map[string]*biz.Role
	permissions     map[string]*biz.Permission
	userRoles       map[string][]string // userID -> roleIDs
	rolePermissions map[string][]string // roleID -> permissionNames

	mutex sync.RWMutex
}

// NewUserRepo 创建用户仓库
func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data:      data,
		log:       log.NewHelper(logger),
		users:     make(map[string]*biz.User),
		userRoles: make(map[string][]string),
	}
}

// Create 创建用户
func (r *userRepo) Create(ctx context.Context, user *biz.User) (*biz.User, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// 生成唯一ID
	if user.ID == "" {
		user.ID = uuid.New().String()
	}

	// 存储用户
	r.users[user.ID] = user
	r.log.WithContext(ctx).Infof("用户创建成功: %s", user.Username)

	return user, nil
}

// Update 更新用户
func (r *userRepo) Update(ctx context.Context, user *biz.User) (*biz.User, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// 检查用户是否存在
	_, ok := r.users[user.ID]
	if !ok {
		return nil, biz.ErrUserNotFound
	}

	// 更新用户
	r.users[user.ID] = user
	r.log.WithContext(ctx).Infof("用户更新成功: %s", user.Username)

	return user, nil
}

// Delete 删除用户
func (r *userRepo) Delete(ctx context.Context, id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// 检查用户是否存在
	_, ok := r.users[id]
	if !ok {
		return biz.ErrUserNotFound
	}

	// 删除用户
	delete(r.users, id)
	// 删除用户角色关联
	delete(r.userRoles, id)

	r.log.WithContext(ctx).Infof("用户删除成功: %s", id)
	return nil
}

// FindByID 根据ID查找用户
func (r *userRepo) FindByID(ctx context.Context, id string) (*biz.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	user, ok := r.users[id]
	if !ok {
		return nil, biz.ErrUserNotFound
	}

	return user, nil
}

// FindByUsername 根据用户名查找用户
func (r *userRepo) FindByUsername(ctx context.Context, username string) (*biz.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, user := range r.users {
		if user.Username == username {
			return user, nil
		}
	}

	return nil, biz.ErrUserNotFound
}

// List 列出用户
func (r *userRepo) List(ctx context.Context, pageSize int32, pageToken string, filter string) ([]*biz.User, string, int32, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	// 简单实现，实际项目中应该实现分页和过滤
	var users []*biz.User
	for _, user := range r.users {
		// 简单过滤实现
		if filter == "" || (filter != "" && (user.Username == filter || user.Email == filter)) {
			users = append(users, user)
		}
	}

	// 返回总数
	totalSize := int32(len(users))

	// 简单分页实现
	start := 0
	if pageToken != "" {
		// 实际项目中应该解析pageToken获取起始位置
		// 这里简化处理
		start = len(users) / 2
	}

	end := start + int(pageSize)
	if end > len(users) {
		end = len(users)
	}

	// 计算下一页的token
	nextPageToken := ""
	if end < len(users) {
		// 实际项目中应该生成有意义的pageToken
		nextPageToken = fmt.Sprintf("page_%d", end)
	}

	// 返回分页结果
	if start < len(users) {
		return users[start:end], nextPageToken, totalSize, nil
	}

	return []*biz.User{}, nextPageToken, totalSize, nil
}

// AssignRole 为用户分配角色
func (r *userRepo) AssignRole(ctx context.Context, userID, roleID string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// 检查用户是否存在
	_, ok := r.users[userID]
	if !ok {
		return biz.ErrUserNotFound
	}

	// 检查角色是否已分配
	roles, ok := r.userRoles[userID]
	if !ok {
		r.userRoles[userID] = []string{roleID}
		return nil
	}

	// 检查是否已经分配了该角色
	for _, id := range roles {
		if id == roleID {
			// 已经分配了该角色，无需重复分配
			return nil
		}
	}

	// 分配角色
	r.userRoles[userID] = append(r.userRoles[userID], roleID)
	r.log.WithContext(ctx).Infof("用户 %s 分配角色 %s 成功", userID, roleID)

	return nil
}

// RemoveRole 从用户移除角色
func (r *userRepo) RemoveRole(ctx context.Context, userID, roleID string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// 检查用户是否存在
	_, ok := r.users[userID]
	if !ok {
		return biz.ErrUserNotFound
	}

	// 检查用户是否有角色
	roles, ok := r.userRoles[userID]
	if !ok {
		// 用户没有角色，无需移除
		return nil
	}

	// 移除角色
	var newRoles []string
	for _, id := range roles {
		if id != roleID {
			newRoles = append(newRoles, id)
		}
	}

	r.userRoles[userID] = newRoles
	r.log.WithContext(ctx).Infof("用户 %s 移除角色 %s 成功", userID, roleID)

	return nil
}

// ListRoles 获取用户角色列表
func (r *userRepo) ListRoles(ctx context.Context, userID string) ([]*biz.Role, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	// 检查用户是否存在
	_, ok := r.users[userID]
	if !ok {
		return nil, biz.ErrUserNotFound
	}

	// 获取用户角色ID列表
	roleIDs, ok := r.userRoles[userID]
	if !ok {
		// 用户没有角色
		return []*biz.Role{}, nil
	}

	// 获取角色详情
	var roles []*biz.Role
	for _, roleID := range roleIDs {
		role, ok := r.data.roles[roleID]
		if ok {
			roles = append(roles, role)
		}
	}

	return roles, nil
}

// roleRepo 角色仓库实现
type roleRepo struct {
	data *Data
	log  *log.Helper
}

// NewRoleRepo 创建角色仓库
func NewRoleRepo(data *Data, logger log.Logger) biz.RoleRepo {
	return &roleRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// Create 创建角色
func (r *roleRepo) Create(ctx context.Context, role *biz.Role) (*biz.Role, error) {
	r.data.mutex.Lock()
	defer r.data.mutex.Unlock()

	// 生成唯一ID
	if role.ID == "" {
		role.ID = uuid.New().String()
	}

	// 存储角色
	if r.data.roles == nil {
		r.data.roles = make(map[string]*biz.Role)
	}
	r.data.roles[role.ID] = role
	r.log.WithContext(ctx).Infof("角色创建成功: %s", role.Name)

	return role, nil
}

// Update 更新角色
func (r *roleRepo) Update(ctx context.Context, role *biz.Role) (*biz.Role, error) {
	r.data.mutex.Lock()
	defer r.data.mutex.Unlock()

	// 检查角色是否存在
	_, ok := r.data.roles[role.ID]
	if !ok {
		return nil, biz.ErrRoleNotFound
	}

	// 更新角色
	r.data.roles[role.ID] = role
	r.log.WithContext(ctx).Infof("角色更新成功: %s", role.Name)

	return role, nil
}

// Delete 删除角色
func (r *roleRepo) Delete(ctx context.Context, id string) error {
	r.data.mutex.Lock()
	defer r.data.mutex.Unlock()

	// 检查角色是否存在
	_, ok := r.data.roles[id]
	if !ok {
		return biz.ErrRoleNotFound
	}

	// 删除角色
	delete(r.data.roles, id)
	// 删除角色权限关联
	delete(r.data.rolePermissions, id)

	// 从所有用户中移除该角色
	for userID, roleIDs := range r.data.userRoles {
		var newRoleIDs []string
		for _, roleID := range roleIDs {
			if roleID != id {
				newRoleIDs = append(newRoleIDs, roleID)
			}
		}
		r.data.userRoles[userID] = newRoleIDs
	}

	r.log.WithContext(ctx).Infof("角色删除成功: %s", id)
	return nil
}

// FindByID 根据ID查找角色
func (r *roleRepo) FindByID(ctx context.Context, id string) (*biz.Role, error) {
	r.data.mutex.RLock()
	defer r.data.mutex.RUnlock()

	role, ok := r.data.roles[id]
	if !ok {
		return nil, biz.ErrRoleNotFound
	}

	return role, nil
}

// List 列出角色
func (r *roleRepo) List(ctx context.Context, pageSize int32, pageToken string) ([]*biz.Role, string, int32, error) {
	r.data.mutex.RLock()
	defer r.data.mutex.RUnlock()

	// 简单实现，实际项目中应该实现分页
	var roles []*biz.Role
	for _, role := range r.data.roles {
		roles = append(roles, role)
	}

	// 返回总数
	totalSize := int32(len(roles))

	// 简单分页实现
	start := 0
	if pageToken != "" {
		// 实际项目中应该解析pageToken获取起始位置
		// 这里简化处理
		start = len(roles) / 2
	}

	end := start + int(pageSize)
	if end > len(roles) {
		end = len(roles)
	}

	// 计算下一页的token
	nextPageToken := ""
	if end < len(roles) {
		// 实际项目中应该生成有意义的pageToken
		nextPageToken = fmt.Sprintf("page_%d", end)
	}

	// 返回分页结果
	if start < len(roles) {
		return roles[start:end], nextPageToken, totalSize, nil
	}

	return []*biz.Role{}, nextPageToken, totalSize, nil
}

// AddPermission 为角色添加权限
func (r *roleRepo) AddPermission(ctx context.Context, roleID, permissionName string) error {
	r.data.mutex.Lock()
	defer r.data.mutex.Unlock()

	// 检查角色是否存在
	_, ok := r.data.roles[roleID]
	if !ok {
		return biz.ErrRoleNotFound
	}

	// 初始化角色权限映射
	if r.data.rolePermissions == nil {
		r.data.rolePermissions = make(map[string][]string)
	}

	// 检查权限是否已分配
	permissions, ok := r.data.rolePermissions[roleID]
	if !ok {
		r.data.rolePermissions[roleID] = []string{permissionName}
		return nil
	}

	// 检查是否已经分配了该权限
	for _, name := range permissions {
		if name == permissionName {
			// 已经分配了该权限，无需重复分配
			return nil
		}
	}

	// 分配权限
	r.data.rolePermissions[roleID] = append(r.data.rolePermissions[roleID], permissionName)
	r.log.WithContext(ctx).Infof("角色 %s 添加权限 %s 成功", roleID, permissionName)

	return nil
}

// RemovePermission 从角色移除权限
func (r *roleRepo) RemovePermission(ctx context.Context, roleID, permissionName string) error {
	r.data.mutex.Lock()
	defer r.data.mutex.Unlock()

	// 检查角色是否存在
	_, ok := r.data.roles[roleID]
	if !ok {
		return biz.ErrRoleNotFound
	}

	// 检查角色是否有权限
	permissions, ok := r.data.rolePermissions[roleID]
	if !ok {
		// 角色没有权限，无需移除
		return nil
	}

	// 移除权限
	var newPermissions []string
	for _, name := range permissions {
		if name != permissionName {
			newPermissions = append(newPermissions, name)
		}
	}

	r.data.rolePermissions[roleID] = newPermissions
	r.log.WithContext(ctx).Infof("角色 %s 移除权限 %s 成功", roleID, permissionName)

	return nil
}

// ListPermissions 列出角色权限
func (r *roleRepo) ListPermissions(ctx context.Context, roleID string) ([]*biz.Permission, error) {
	r.data.mutex.RLock()
	defer r.data.mutex.RUnlock()

	// 检查角色是否存在
	_, ok := r.data.roles[roleID]
	if !ok {
		return nil, biz.ErrRoleNotFound
	}

	// 获取角色权限名称列表
	permissionNames, ok := r.data.rolePermissions[roleID]
	if !ok {
		// 角色没有权限
		return []*biz.Permission{}, nil
	}

	// 获取权限详情
	var permissions []*biz.Permission
	for _, name := range permissionNames {
		permission, ok := r.data.permissions[name]
		if ok {
			permissions = append(permissions, permission)
		}
	}

	return permissions, nil
}

// permissionRepo 权限仓库实现
type permissionRepo struct {
	data *Data
	log  *log.Helper
}

// NewPermissionRepo 创建权限仓库
func NewPermissionRepo(data *Data, logger log.Logger) biz.PermissionRepo {
	return &permissionRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// Create 创建权限
func (r *permissionRepo) Create(ctx context.Context, permission *biz.Permission) (*biz.Permission, error) {
	r.data.mutex.Lock()
	defer r.data.mutex.Unlock()

	// 初始化权限映射
	if r.data.permissions == nil {
		r.data.permissions = make(map[string]*biz.Permission)
	}

	// 检查权限是否已存在
	_, ok := r.data.permissions[permission.Name]
	if ok {
		return nil, fmt.Errorf("权限已存在: %s", permission.Name)
	}

	// 存储权限
	r.data.permissions[permission.Name] = permission
	r.log.WithContext(ctx).Infof("权限创建成功: %s", permission.Name)

	return permission, nil
}

// FindByName 根据名称查找权限
func (r *permissionRepo) FindByName(ctx context.Context, name string) (*biz.Permission, error) {
	r.data.mutex.RLock()
	defer r.data.mutex.RUnlock()

	permission, ok := r.data.permissions[name]
	if !ok {
		return nil, biz.ErrPermissionNotFound
	}

	return permission, nil
}

// List 列出权限
func (r *permissionRepo) List(ctx context.Context, resourceType string) ([]*biz.Permission, error) {
	r.data.mutex.RLock()
	defer r.data.mutex.RUnlock()

	var permissions []*biz.Permission
	for _, permission := range r.data.permissions {
		// 如果指定了资源类型，则过滤
		if resourceType == "" || permission.ResourceType == resourceType {
			permissions = append(permissions, permission)
		}
	}

	return permissions, nil
}

// 更新Data结构，添加用户管理相关字段
func init() {
	// 预设一些权限
	defaultPermissions := []*biz.Permission{
		{Name: "user:read", Description: "查看用户信息", ResourceType: "user", Action: "read"},
		{Name: "user:write", Description: "修改用户信息", ResourceType: "user", Action: "write"},
		{Name: "user:delete", Description: "删除用户", ResourceType: "user", Action: "delete"},
		{Name: "role:read", Description: "查看角色信息", ResourceType: "role", Action: "read"},
		{Name: "role:write", Description: "修改角色信息", ResourceType: "role", Action: "write"},
		{Name: "role:delete", Description: "删除角色", ResourceType: "role", Action: "delete"},
	}

	// 预设一些角色
	defaultRoles := []*biz.Role{
		{
			ID:          "admin",
			Name:        "管理员",
			Description: "系统管理员，拥有所有权限",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "user",
			Name:        "普通用户",
			Description: "普通用户，拥有基本权限",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	// 预设角色权限关系
	defaultRolePermissions := map[string][]string{
		"admin": {"user:read", "user:write", "user:delete", "role:read", "role:write", "role:delete"},
		"user":  {"user:read"},
	}

	// 预设一个管理员用户
	defaultUsers := []*biz.User{
		{
			ID:           "admin",
			Username:     "admin",
			PasswordHash: "admin_password_hashed",
			Email:        "admin@example.com",
			DisplayName:  "系统管理员",
			IsActive:     true,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
	}

	// 预设用户角色关系
	defaultUserRoles := map[string][]string{
		"admin": {"admin"},
	}
}
