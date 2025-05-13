package service

import (
	"context"

	v1 "om-platform/api/user/service/v1"
	"om-platform/app/user/service/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// UserManagementService 实现用户管理服务接口
type UserManagementService struct {
	v1.UnimplementedUserManagementServiceServer

	uc  *biz.UserUsecase
	log *log.Helper
}

// NewUserManagementService 创建用户管理服务实例
func NewUserManagementService(uc *biz.UserUsecase, logger log.Logger) *UserManagementService {
	return &UserManagementService{
		uc:  uc,
		log: log.NewHelper(logger),
	}
}

// 将领域模型转换为API模型
func convertUser(u *biz.User) *v1.User {
	if u == nil {
		return nil
	}
	return &v1.User{
		Id:          u.ID,
		Username:    u.Username,
		Email:       u.Email,
		DisplayName: u.DisplayName,
		IsActive:    u.IsActive,
		CreatedAt:   timestamppb.New(u.CreatedAt),
		UpdatedAt:   timestamppb.New(u.UpdatedAt),
	}
}

// 将API模型转换为领域模型
func convertUserToBiz(u *v1.User) *biz.User {
	if u == nil {
		return nil
	}
	user := &biz.User{
		ID:          u.Id,
		Username:    u.Username,
		Email:       u.Email,
		DisplayName: u.DisplayName,
		IsActive:    u.IsActive,
	}
	if u.CreatedAt != nil {
		user.CreatedAt = u.CreatedAt.AsTime()
	}
	if u.UpdatedAt != nil {
		user.UpdatedAt = u.UpdatedAt.AsTime()
	}
	return user
}

// 将角色领域模型转换为API模型
func convertRole(r *biz.Role) *v1.Role {
	if r == nil {
		return nil
	}
	return &v1.Role{
		Id:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		CreatedAt:   timestamppb.New(r.CreatedAt),
		UpdatedAt:   timestamppb.New(r.UpdatedAt),
	}
}

// 将API角色模型转换为领域模型
func convertRoleToBiz(r *v1.Role) *biz.Role {
	if r == nil {
		return nil
	}
	role := &biz.Role{
		ID:          r.Id,
		Name:        r.Name,
		Description: r.Description,
	}
	if r.CreatedAt != nil {
		role.CreatedAt = r.CreatedAt.AsTime()
	}
	if r.UpdatedAt != nil {
		role.UpdatedAt = r.UpdatedAt.AsTime()
	}
	return role
}

// 将权限领域模型转换为API模型
func convertPermission(p *biz.Permission) *v1.Permission {
	if p == nil {
		return nil
	}
	return &v1.Permission{
		Name:         p.Name,
		Description:  p.Description,
		ResourceType: p.ResourceType,
		Action:       p.Action,
	}
}

// 将API权限模型转换为领域模型
func convertPermissionToBiz(p *v1.Permission) *biz.Permission {
	if p == nil {
		return nil
	}
	return &biz.Permission{
		Name:         p.Name,
		Description:  p.Description,
		ResourceType: p.ResourceType,
		Action:       p.Action,
	}
}

// CreateUser 创建用户
func (s *UserManagementService) CreateUser(ctx context.Context, req *v1.CreateUserRequest) (*v1.User, error) {
	s.log.WithContext(ctx).Infof("创建用户: %s", req.User.Username)

	user, err := s.uc.CreateUser(ctx, convertUserToBiz(req.User), req.Password)
	if err != nil {
		return nil, err
	}

	return convertUser(user), nil
}

// GetUser 获取用户信息
func (s *UserManagementService) GetUser(ctx context.Context, req *v1.GetUserRequest) (*v1.User, error) {
	s.log.WithContext(ctx).Infof("获取用户信息: %s", req.UserId)

	user, err := s.uc.GetUser(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	return convertUser(user), nil
}

// UpdateUser 更新用户信息
func (s *UserManagementService) UpdateUser(ctx context.Context, req *v1.UpdateUserRequest) (*v1.User, error) {
	s.log.WithContext(ctx).Infof("更新用户信息: %s", req.User.Id)

	user, err := s.uc.UpdateUser(ctx, convertUserToBiz(req.User), req.Password)
	if err != nil {
		return nil, err
	}

	return convertUser(user), nil
}

// DeleteUser 删除用户
func (s *UserManagementService) DeleteUser(ctx context.Context, req *v1.DeleteUserRequest) (*emptypb.Empty, error) {
	s.log.WithContext(ctx).Infof("删除用户: %s", req.UserId)

	err := s.uc.DeleteUser(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

// ListUsers 列出用户
func (s *UserManagementService) ListUsers(ctx context.Context, req *v1.ListUsersRequest) (*v1.ListUsersResponse, error) {
	s.log.WithContext(ctx).Info("列出用户")

	users, nextPageToken, totalSize, err := s.uc.ListUsers(ctx, req.PageSize, req.PageToken, req.Filter)
	if err != nil {
		return nil, err
	}

	var userList []*v1.User
	for _, user := range users {
		userList = append(userList, convertUser(user))
	}

	return &v1.ListUsersResponse{
		Users:         userList,
		NextPageToken: nextPageToken,
		TotalSize:     totalSize,
	}, nil
}

// AssignRoleToUser 为用户分配角色
func (s *UserManagementService) AssignRoleToUser(ctx context.Context, req *v1.AssignRoleToUserRequest) (*emptypb.Empty, error) {
	s.log.WithContext(ctx).Infof("为用户 %s 分配角色 %s", req.UserId, req.RoleId)

	err := s.uc.AssignRoleToUser(ctx, req.UserId, req.RoleId)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

// RemoveRoleFromUser 从用户移除角色
func (s *UserManagementService) RemoveRoleFromUser(ctx context.Context, req *v1.RemoveRoleFromUserRequest) (*emptypb.Empty, error) {
	s.log.WithContext(ctx).Infof("从用户 %s 移除角色 %s", req.UserId, req.RoleId)

	err := s.uc.RemoveRoleFromUser(ctx, req.UserId, req.RoleId)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

// ListUserRoles 获取用户角色列表
func (s *UserManagementService) ListUserRoles(ctx context.Context, req *v1.ListUserRolesRequest) (*v1.ListUserRolesResponse, error) {
	s.log.WithContext(ctx).Infof("获取用户 %s 的角色列表", req.UserId)

	roles, err := s.uc.ListUserRoles(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	var roleList []*v1.Role
	for _, role := range roles {
		roleList = append(roleList, convertRole(role))
	}

	return &v1.ListUserRolesResponse{
		Roles: roleList,
	}, nil
}

// CreateRole 创建角色
func (s *UserManagementService) CreateRole(ctx context.Context, req *v1.CreateRoleRequest) (*v1.Role, error) {
	s.log.WithContext(ctx).Infof("创建角色: %s", req.Role.Name)

	role, err := s.uc.CreateRole(ctx, convertRoleToBiz(req.Role))
	if err != nil {
		return nil, err
	}

	return convertRole(role), nil
}

// GetRole 获取角色信息
func (s *UserManagementService) GetRole(ctx context.Context, req *v1.GetRoleRequest) (*v1.Role, error) {
	s.log.WithContext(ctx).Infof("获取角色信息: %s", req.RoleId)

	role, err := s.uc.GetRole(ctx, req.RoleId)
	if err != nil {
		return nil, err
	}

	return convertRole(role), nil
}

// UpdateRole 更新角色信息
func (s *UserManagementService) UpdateRole(ctx context.Context, req *v1.UpdateRoleRequest) (*v1.Role, error) {
	s.log.WithContext(ctx).Infof("更新角色信息: %s", req.Role.Id)

	role, err := s.uc.UpdateRole(ctx, convertRoleToBiz(req.Role))
	if err != nil {
		return nil, err
	}

	return convertRole(role), nil
}

// DeleteRole 删除角色
func (s *UserManagementService) DeleteRole(ctx context.Context, req *v1.DeleteRoleRequest) (*emptypb.Empty, error) {
	s.log.WithContext(ctx).Infof("删除角色: %s", req.RoleId)

	err := s.uc.DeleteRole(ctx, req.RoleId)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

// ListRoles 列出角色
func (s *UserManagementService) ListRoles(ctx context.Context, req *v1.ListRolesRequest) (*v1.ListRolesResponse, error) {
	s.log.WithContext(ctx).Info("列出角色")

	roles, nextPageToken, totalSize, err := s.uc.ListRoles(ctx, req.PageSize, req.PageToken)
	if err != nil {
		return nil, err
	}

	var roleList []*v1.Role
	for _, role := range roles {
		roleList = append(roleList, convertRole(role))
	}

	return &v1.ListRolesResponse{
		Roles:         roleList,
		NextPageToken: nextPageToken,
		TotalSize:     totalSize,
	}, nil
}

// AddPermissionToRole 为角色添加权限
func (s *UserManagementService) AddPermissionToRole(ctx context.Context, req *v1.AddPermissionToRoleRequest) (*emptypb.Empty, error) {
	s.log.WithContext(ctx).Infof("为角色 %s 添加权限 %s", req.RoleId, req.PermissionName)

	err := s.uc.AddPermissionToRole(ctx, req.RoleId, req.PermissionName)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

// RemovePermissionFromRole 从角色移除权限
func (s *UserManagementService) RemovePermissionFromRole(ctx context.Context, req *v1.RemovePermissionFromRoleRequest) (*emptypb.Empty, error) {
	s.log.WithContext(ctx).Infof("从角色 %s 移除权限 %s", req.RoleId, req.PermissionName)

	err := s.uc.RemovePermissionFromRole(ctx, req.RoleId, req.PermissionName)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

// ListRolePermissions 列出角色权限
func (s *UserManagementService) ListRolePermissions(ctx context.Context, req *v1.ListRolePermissionsRequest) (*v1.ListRolePermissionsResponse, error) {
	s.log.WithContext(ctx).Infof("获取角色 %s 的权限列表", req.RoleId)

	permissions, err := s.uc.ListRolePermissions(ctx, req.RoleId)
	if err != nil {
		return nil, err
	}

	var permissionList []*v1.Permission
	for _, permission := range permissions {
		permissionList = append(permissionList, convertPermission(permission))
	}

	return &v1.ListRolePermissionsResponse{
		Permissions: permissionList,
	}, nil
}

// CreatePermission 创建权限
func (s *UserManagementService) CreatePermission(ctx context.Context, req *v1.CreatePermissionRequest) (*v1.Permission, error) {
	s.log.WithContext(ctx).Infof("创建权限: %s", req.Permission.Name)

	permission, err := s.uc.CreatePermission(ctx, convertPermissionToBiz(req.Permission))
	if err != nil {
		return nil, err
	}

	return convertPermission(permission), nil
}

// ListPermissions 列出权限
func (s *UserManagementService) ListPermissions(ctx context.Context, req *v1.ListPermissionsRequest) (*v1.ListPermissionsResponse, error) {
	s.log.WithContext(ctx).Info("列出权限")

	permissions, err := s.uc.ListPermissions(ctx, req.ResourceType)
	if err != nil {
		return nil, err
	}

	var permissionList []*v1.Permission
	for _, permission := range permissions {
		permissionList = append(permissionList, convertPermission(permission))
	}

	return &v1.ListPermissionsResponse{
		Permissions: permissionList,
	}, nil
}
