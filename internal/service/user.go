package service

import (
	"context"

	v1 "om-platform/api/user/service/v1"
	"om-platform/internal/biz"

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

// 将业务实体转换为API响应
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

// CreateUser 创建用户
func (s *UserManagementService) CreateUser(ctx context.Context, req *v1.CreateUserRequest) (*v1.User, error) {
	s.log.WithContext(ctx).Infof("创建用户: %v", req.GetUser().GetUsername())

	// 转换请求为业务实体
	user := &biz.User{
		Username:    req.GetUser().GetUsername(),
		Email:       req.GetUser().GetEmail(),
		DisplayName: req.GetUser().GetDisplayName(),
		IsActive:    req.GetUser().GetIsActive(),
	}

	// 调用业务逻辑
	result, err := s.uc.CreateUser(ctx, user, req.GetPassword())
	if err != nil {
		return nil, err
	}

	// 转换业务实体为响应
	return convertUser(result), nil
}

// GetUser 获取用户信息
func (s *UserManagementService) GetUser(ctx context.Context, req *v1.GetUserRequest) (*v1.User, error) {
	s.log.WithContext(ctx).Infof("获取用户: %v", req.GetUserId())

	// 调用业务逻辑
	user, err := s.uc.GetUser(ctx, req.GetUserId())
	if err != nil {
		return nil, err
	}

	// 转换业务实体为响应
	return convertUser(user), nil
}

// UpdateUser 更新用户信息
func (s *UserManagementService) UpdateUser(ctx context.Context, req *v1.UpdateUserRequest) (*v1.User, error) {
	s.log.WithContext(ctx).Infof("更新用户: %v", req.GetUser().GetId())

	// 转换请求为业务实体
	user := &biz.User{
		ID:          req.GetUser().GetId(),
		Username:    req.GetUser().GetUsername(),
		Email:       req.GetUser().GetEmail(),
		DisplayName: req.GetUser().GetDisplayName(),
		IsActive:    req.GetUser().GetIsActive(),
	}

	// 调用业务逻辑
	result, err := s.uc.UpdateUser(ctx, user, req.GetPassword())
	if err != nil {
		return nil, err
	}

	// 转换业务实体为响应
	return convertUser(result), nil
}

// DeleteUser 删除用户
func (s *UserManagementService) DeleteUser(ctx context.Context, req *v1.DeleteUserRequest) (*emptypb.Empty, error) {
	s.log.WithContext(ctx).Infof("删除用户: %v", req.GetUserId())

	// 调用业务逻辑
	err := s.uc.DeleteUser(ctx, req.GetUserId())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

// ListUsers 列出用户
func (s *UserManagementService) ListUsers(ctx context.Context, req *v1.ListUsersRequest) (*v1.ListUsersResponse, error) {
	s.log.WithContext(ctx).Info("列出用户")

	// 调用业务逻辑
	users, nextPageToken, totalSize, err := s.uc.ListUsers(ctx, req.GetPageSize(), req.GetPageToken(), req.GetFilter())
	if err != nil {
		return nil, err
	}

	// 转换业务实体为响应
	response := &v1.ListUsersResponse{
		NextPageToken: nextPageToken,
		TotalSize:     totalSize,
		Users:         make([]*v1.User, 0, len(users)),
	}

	for _, user := range users {
		response.Users = append(response.Users, convertUser(user))
	}

	return response, nil
}

// AssignRoleToUser 为用户分配角色
func (s *UserManagementService) AssignRoleToUser(ctx context.Context, req *v1.AssignRoleToUserRequest) (*emptypb.Empty, error) {
	s.log.WithContext(ctx).Infof("为用户 %s 分配角色 %s", req.GetUserId(), req.GetRoleId())

	// 调用业务逻辑
	err := s.uc.AssignRoleToUser(ctx, req.GetUserId(), req.GetRoleId())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

// RemoveRoleFromUser 从用户移除角色
func (s *UserManagementService) RemoveRoleFromUser(ctx context.Context, req *v1.RemoveRoleFromUserRequest) (*emptypb.Empty, error) {
	s.log.WithContext(ctx).Infof("从用户 %s 移除角色 %s", req.GetUserId(), req.GetRoleId())

	// 调用业务逻辑
	err := s.uc.RemoveRoleFromUser(ctx, req.GetUserId(), req.GetRoleId())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

// ListUserRoles 获取用户角色列表
func (s *UserManagementService) ListUserRoles(ctx context.Context, req *v1.ListUserRolesRequest) (*v1.ListUserRolesResponse, error) {
	s.log.WithContext(ctx).Infof("获取用户 %s 的角色列表", req.GetUserId())

	// 调用业务逻辑
	roles, err := s.uc.ListUserRoles(ctx, req.GetUserId())
	if err != nil {
		return nil, err
	}

	// 转换业务实体为响应
	response := &v1.ListUserRolesResponse{
		Roles: make([]*v1.Role, 0, len(roles)),
	}

	for _, role := range roles {
		response.Roles = append(response.Roles, convertRole(role))
	}

	return response, nil
}

// CreateRole 创建角色
func (s *UserManagementService) CreateRole(ctx context.Context, req *v1.CreateRoleRequest) (*v1.Role, error) {
	s.log.WithContext(ctx).Infof("创建角色: %v", req.GetRole().GetName())

	// 转换请求为业务实体
	role := &biz.Role{
		Name:        req.GetRole().GetName(),
		Description: req.GetRole().GetDescription(),
	}

	// 调用业务逻辑
	result, err := s.uc.CreateRole(ctx, role)
	if err != nil {
		return nil, err
	}

	// 转换业务实体为响应
	return convertRole(result), nil
}

// GetRole 获取角色信息
func (s *UserManagementService) GetRole(ctx context.Context, req *v1.GetRoleRequest) (*v1.Role, error) {
	s.log.WithContext(ctx).Infof("获取角色: %v", req.GetRoleId())

	// 调用业务逻辑
	role, err := s.uc.GetRole(ctx, req.GetRoleId())
	if err != nil {
		return nil, err
	}

	// 转换业务实体为响应
	return convertRole(role), nil
}

// UpdateRole 更新角色信息
func (s *UserManagementService) UpdateRole(ctx context.Context, req *v1.UpdateRoleRequest) (*v1.Role, error) {
	s.log.WithContext(ctx).Infof("更新角色: %v", req.GetRole().GetId())

	// 转换请求为业务实体
	role := &biz.Role{
		ID:          req.GetRole().GetId(),
		Name:        req.GetRole().GetName(),
		Description: req.GetRole().GetDescription(),
	}

	// 调用业务逻辑
	result, err := s.uc.UpdateRole(ctx, role)
	if err != nil {
		return nil, err
	}

	// 转换业务实体为响应
	return convertRole(result), nil
}

// DeleteRole 删除角色
func (s *UserManagementService) DeleteRole(ctx context.Context, req *v1.DeleteRoleRequest) (*emptypb.Empty, error) {
	s.log.WithContext(ctx).Infof("删除角色: %v", req.GetRoleId())

	// 调用业务逻辑
	err := s.uc.DeleteRole(ctx, req.GetRoleId())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

// ListRoles 列出角色
func (s *UserManagementService) ListRoles(ctx context.Context, req *v1.ListRolesRequest) (*v1.ListRolesResponse, error) {
	s.log.WithContext(ctx).Info("列出角色")

	// 调用业务逻辑
	roles, nextPageToken, totalSize, err := s.uc.ListRoles(ctx, req.GetPageSize(), req.GetPageToken())
	if err != nil {
		return nil, err
	}

	// 转换业务实体为响应
	response := &v1.ListRolesResponse{
		NextPageToken: nextPageToken,
		TotalSize:     totalSize,
		Roles:         make([]*v1.Role, 0, len(roles)),
	}

	for _, role := range roles {
		response.Roles = append(response.Roles, convertRole(role))
	}

	return response, nil
}

// AddPermissionToRole 为角色添加权限
func (s *UserManagementService) AddPermissionToRole(ctx context.Context, req *v1.AddPermissionToRoleRequest) (*emptypb.Empty, error) {
	s.log.WithContext(ctx).Infof("为角色 %s 添加权限 %s", req.GetRoleId(), req.GetPermissionName())

	// 调用业务逻辑
	err := s.uc.AddPermissionToRole(ctx, req.GetRoleId(), req.GetPermissionName())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

// RemovePermissionFromRole 从角色移除权限
func (s *UserManagementService) RemovePermissionFromRole(ctx context.Context, req *v1.RemovePermissionFromRoleRequest) (*emptypb.Empty, error) {
	s.log.WithContext(ctx).Infof("从角色 %s 移除权限 %s", req.GetRoleId(), req.GetPermissionName())

	// 调用业务逻辑
	err := s.uc.RemovePermissionFromRole(ctx, req.GetRoleId(), req.GetPermissionName())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

// ListRolePermissions 列出角色权限
func (s *UserManagementService) ListRolePermissions(ctx context.Context, req *v1.ListRolePermissionsRequest) (*v1.ListRolePermissionsResponse, error) {
	s.log.WithContext(ctx).Infof("列出角色 %s 的权限", req.GetRoleId())

	// 调用业务逻辑
	permissions, err := s.uc.ListRolePermissions(ctx, req.GetRoleId())
	if err != nil {
		return nil, err
	}

	// 转换业务实体为响应
	response := &v1.ListRolePermissionsResponse{
		Permissions: make([]*v1.Permission, 0, len(permissions)),
	}

	for _, permission := range permissions {
		response.Permissions = append(response.Permissions, convertPermission(permission))
	}

	return response, nil
}

// CreatePermission 创建权限定义
func (s *UserManagementService) CreatePermission(ctx context.Context, req *v1.CreatePermissionRequest) (*v1.Permission, error) {
	s.log.WithContext(ctx).Infof("创建权限: %v", req.GetPermission().GetName())

	// 转换请求为业务实体
	permission := &biz.Permission{
		Name:         req.GetPermission().GetName(),
		Description:  req.GetPermission().GetDescription(),
		ResourceType: req.GetPermission().GetResourceType(),
		Action:       req.GetPermission().GetAction(),
	}

	// 调用业务逻辑
	result, err := s.uc.CreatePermission(ctx, permission)
	if err != nil {
		return nil, err
	}

	// 转换业务实体为响应
	return convertPermission(result), nil
}

// ListPermissions 列出所有权限定义
func (s *UserManagementService) ListPermissions(ctx context.Context, req *v1.ListPermissionsRequest) (*v1.ListPermissionsResponse, error) {
	s.log.WithContext(ctx).Info("列出权限")

	// 调用业务逻辑
	permissions, err := s.uc.ListPermissions(ctx, req.GetResourceType())
	if err != nil {
		return nil, err
	}

	// 转换业务实体为响应
	response := &v1.ListPermissionsResponse{
		Permissions: make([]*v1.Permission, 0, len(permissions)),
	}

	for _, permission := range permissions {
		response.Permissions = append(response.Permissions, convertPermission(permission))
	}

	return response, nil
}

// 服务层的Provider集合已在service.go中定义
func init() {
	// 无需重复定义ProviderSet
}
