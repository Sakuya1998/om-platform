# 用户服务API优化示例代码

本文档提供了用户服务API优化的具体示例代码，作为优化方案的补充说明。

## 1. 统一错误码定义示例

```protobuf
// api/user/service/v1/error_codes.proto
syntax = "proto3";

package api.user.service.v1;

import "gnostic/openapi/v3/annotations.proto";

option go_package = "om-platform/api/user/service/v1;v1";
option java_multiple_files = true;
option java_package = "api.user.service.v1";

// 统一错误码定义
// 错误码按模块划分范围，便于管理和排查问题
enum ErrorCode {
  // 通用错误码 (0-999)
  ERROR_CODE_UNSPECIFIED = 0 [
    (gnostic.openapi.v3.property) = {description: "未指定错误"}
  ];
  ERROR_CODE_INTERNAL = 1 [
    (gnostic.openapi.v3.property) = {description: "内部服务错误"}
  ];
  ERROR_CODE_INVALID_ARGUMENT = 2 [
    (gnostic.openapi.v3.property) = {description: "无效参数"}
  ];
  ERROR_CODE_NOT_FOUND = 3 [
    (gnostic.openapi.v3.property) = {description: "资源不存在"}
  ];
  ERROR_CODE_PERMISSION_DENIED = 4 [
    (gnostic.openapi.v3.property) = {description: "权限不足"}
  ];
  ERROR_CODE_UNAUTHENTICATED = 5 [
    (gnostic.openapi.v3.property) = {description: "未认证"}
  ];
  ERROR_CODE_RESOURCE_EXHAUSTED = 6 [
    (gnostic.openapi.v3.property) = {description: "资源耗尽"}
  ];
  ERROR_CODE_FAILED_PRECONDITION = 7 [
    (gnostic.openapi.v3.property) = {description: "前置条件失败"}
  ];
  ERROR_CODE_ABORTED = 8 [
    (gnostic.openapi.v3.property) = {description: "操作中止"}
  ];
  ERROR_CODE_DEADLINE_EXCEEDED = 9 [
    (gnostic.openapi.v3.property) = {description: "超时"}
  ];
  ERROR_CODE_ALREADY_EXISTS = 10 [
    (gnostic.openapi.v3.property) = {description: "资源已存在"}
  ];
  
  // 用户相关错误码 (1000-1999)
  ERROR_CODE_USER_NOT_FOUND = 1000 [
    (gnostic.openapi.v3.property) = {description: "用户不存在"}
  ];
  ERROR_CODE_USER_ALREADY_EXISTS = 1001 [
    (gnostic.openapi.v3.property) = {description: "用户已存在"}
  ];
  ERROR_CODE_USER_DISABLED = 1002 [
    (gnostic.openapi.v3.property) = {description: "用户已禁用"}
  ];
  ERROR_CODE_USER_LOCKED = 1003 [
    (gnostic.openapi.v3.property) = {description: "用户已锁定"}
  ];
  ERROR_CODE_USER_DELETED = 1004 [
    (gnostic.openapi.v3.property) = {description: "用户已删除"}
  ];
  
  // 认证相关错误码 (2000-2999)
  ERROR_CODE_AUTHENTICATION_FAILED = 2000 [
    (gnostic.openapi.v3.property) = {description: "认证失败"}
  ];
  ERROR_CODE_TOKEN_EXPIRED = 2001 [
    (gnostic.openapi.v3.property) = {description: "令牌已过期"}
  ];
  ERROR_CODE_TOKEN_INVALID = 2002 [
    (gnostic.openapi.v3.property) = {description: "令牌无效"}
  ];
  ERROR_CODE_TOKEN_REVOKED = 2003 [
    (gnostic.openapi.v3.property) = {description: "令牌已撤销"}
  ];
  ERROR_CODE_SESSION_EXPIRED = 2004 [
    (gnostic.openapi.v3.property) = {description: "会话已过期"}
  ];
  ERROR_CODE_SESSION_INVALID = 2005 [
    (gnostic.openapi.v3.property) = {description: "会话无效"}
  ];
  ERROR_CODE_PASSWORD_EXPIRED = 2006 [
    (gnostic.openapi.v3.property) = {description: "密码已过期"}
  ];
  ERROR_CODE_PASSWORD_INCORRECT = 2007 [
    (gnostic.openapi.v3.property) = {description: "密码不正确"}
  ];
  
  // 权限相关错误码 (3000-3999)
  ERROR_CODE_PERMISSION_NOT_FOUND = 3000 [
    (gnostic.openapi.v3.property) = {description: "权限不存在"}
  ];
  ERROR_CODE_ROLE_NOT_FOUND = 3001 [
    (gnostic.openapi.v3.property) = {description: "角色不存在"}
  ];
  ERROR_CODE_INSUFFICIENT_PERMISSIONS = 3002 [
    (gnostic.openapi.v3.property) = {description: "权限不足"}
  ];
  
  // 组织相关错误码 (4000-4999)
  ERROR_CODE_ORGANIZATION_NOT_FOUND = 4000 [
    (gnostic.openapi.v3.property) = {description: "组织不存在"}
  ];
  ERROR_CODE_DEPARTMENT_NOT_FOUND = 4001 [
    (gnostic.openapi.v3.property) = {description: "部门不存在"}
  ];
  ERROR_CODE_USER_GROUP_NOT_FOUND = 4002 [
    (gnostic.openapi.v3.property) = {description: "用户组不存在"}
  ];
}

// 统一错误响应结构
message ErrorResponse {
  ErrorCode code = 1 [
    json_name = "code",
    (gnostic.openapi.v3.property) = {description: "错误码"}
  ];
  
  string message = 2 [
    json_name = "message",
    (gnostic.openapi.v3.property) = {description: "错误消息"}
  ];
  
  string request_id = 3 [
    json_name = "requestId",
    (gnostic.openapi.v3.property) = {description: "请求ID，用于跟踪和排查问题"}
  ];
  
  map<string, string> details = 4 [
    json_name = "details",
    (gnostic.openapi.v3.property) = {description: "错误详情，包含额外的上下文信息"}
  ];
}
```

## 2. 重构后的认证服务示例

```protobuf
// api/user/service/v1/authentication.proto
syntax = "proto3";

package api.user.service.v1;

import "google/protobuf/empty.proto";
import "google/protobuf/timestamp.proto";
import "gnostic/openapi/v3/annotations.proto";

import "api/user/service/v1/error_codes.proto";
import "api/user/service/v1/common.proto";

option go_package = "om-platform/api/user/service/v1;v1";
option java_multiple_files = true;
option java_package = "api.user.service.v1";

// 认证服务
// 提供用户认证、令牌管理等核心功能
// 错误处理：所有接口在发生错误时将返回ErrorResponse结构，包含统一的错误码、错误消息和详细信息
// 性能说明：登录接口已实现限流保护，避免暴力破解
// 安全说明：所有接口已实现防重放攻击和CSRF保护
service AuthenticationService {
  // 用户登录
  // 支持多种认证方式，包括密码、LDAP、SSO和OAuth2
  // 示例请求:
  // {
  //   "username": "john.doe",
  //   "password": "secure_password",
  //   "auth_type": "PASSWORD",
  //   "device_info": {
  //     "device_type": "WEB",
  //     "ip_address": "192.168.1.1",
  //     "user_agent": "Mozilla/5.0..."
  //   }
  // }
  //
  // 示例响应:
  // {
  //   "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  //   "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  //   "expires_in": 3600,
  //   "token_type": "Bearer",
  //   "user_info": {
  //     "id": 123,
  //     "username": "john.doe",
  //     "display_name": "John Doe"
  //   }
  // }
  //
  // 可能的错误码:
  // - ERROR_CODE_AUTHENTICATION_FAILED: 认证失败
  // - ERROR_CODE_USER_LOCKED: 用户已锁定
  // - ERROR_CODE_USER_DISABLED: 用户已禁用
  // - ERROR_CODE_PASSWORD_EXPIRED: 密码已过期
  // - ERROR_CODE_INVALID_ARGUMENT: 无效参数
  rpc Login (LoginRequest) returns (LoginResponse) {
    option (rate_limit) = {
      requests_per_minute: 10
      burst: 5
    };
  }
  
  // 用户登出
  // 终止当前会话并使相关令牌失效
  rpc Logout (LogoutRequest) returns (google.protobuf.Empty) {}
  
  // 刷新令牌
  // 使用刷新令牌获取新的访问令牌，延长会话有效期
  rpc RefreshToken (RefreshTokenRequest) returns (RefreshTokenResponse) {}
  
  // 验证令牌
  // 检查令牌有效性并返回关联的用户信息和权限
  rpc ValidateToken (ValidateTokenRequest) returns (ValidateTokenResponse) {
    option (cache) = {
      cacheable: true
      ttl_seconds: 60
      cache_keys: ["token"]
    };
  }
  
  // 修改密码
  // 用户主动修改自己的密码，需要提供旧密码验证
  rpc ChangePassword (ChangePasswordRequest) returns (google.protobuf.Empty) {}
  
  // 重置密码
  // 管理员或通过重置流程重置用户密码，无需旧密码
  rpc ResetPassword (ResetPasswordRequest) returns (ResetPasswordResponse) {}
}

// 认证类型
enum AuthType {
  AUTH_TYPE_UNSPECIFIED = 0 [
    (gnostic.openapi.v3.property) = {description: "未指定"}
  ];
  AUTH_TYPE_PASSWORD = 1 [
    (gnostic.openapi.v3.property) = {description: "密码认证"}
  ];
  AUTH_TYPE_LDAP = 2 [
    (gnostic.openapi.v3.property) = {description: "LDAP认证"}
  ];
  AUTH_TYPE_SSO = 3 [
    (gnostic.openapi.v3.property) = {description: "单点登录"}
  ];
  AUTH_TYPE_OAUTH2 = 4 [
    (gnostic.openapi.v3.property) = {description: "OAuth2认证"}
  ];
  AUTH_TYPE_MFA = 5 [
    (gnostic.openapi.v3.property) = {description: "多因素认证"}
  ];
}

// 设备类型
enum DeviceType {
  DEVICE_TYPE_UNSPECIFIED = 0 [
    (gnostic.openapi.v3.property) = {description: "未指定"}
  ];
  DEVICE_TYPE_WEB = 1 [
    (gnostic.openapi.v3.property) = {description: "Web浏览器"}
  ];
  DEVICE_TYPE_MOBILE = 2 [
    (gnostic.openapi.v3.property) = {description: "移动设备"}
  ];
  DEVICE_TYPE_DESKTOP = 3 [
    (gnostic.openapi.v3.property) = {description: "桌面应用"}
  ];
  DEVICE_TYPE_API = 4 [
    (gnostic.openapi.v3.property) = {description: "API调用"}
  ];
}

// 设备信息
message DeviceInfo {
  DeviceType device_type = 1 [
    json_name = "deviceType",
    (gnostic.openapi.v3.property) = {description: "设备类型"}
  ];
  
  string ip_address = 2 [
    json_name = "ipAddress",
    (gnostic.openapi.v3.property) = {description: "IP地址"}
  ];
  
  string user_agent = 3 [
    json_name = "userAgent",
    (gnostic.openapi.v3.property) = {description: "用户代理"}
  ];
  
  string device_id = 4 [
    json_name = "deviceId",
    (gnostic.openapi.v3.property) = {description: "设备ID"}
  ];
  
  map<string, string> extra_info = 5 [
    json_name = "extraInfo",
    (gnostic.openapi.v3.property) = {description: "额外设备信息"}
  ];
}

// 登录请求
message LoginRequest {
  string username = 1 [
    json_name = "username",
    (gnostic.openapi.v3.property) = {description: "用户名"}
  ];
  
  string password = 2 [
    json_name = "password",
    (gnostic.openapi.v3.property) = {description: "密码"}
  ];
  
  AuthType auth_type = 3 [
    json_name = "authType",
    (gnostic.openapi.v3.property) = {description: "认证类型"}
  ];
  
  DeviceInfo device_info = 4 [
    json_name = "deviceInfo",
    (gnostic.openapi.v3.property) = {description: "设备信息"}
  ];
  
  uint32 tenant_id = 5 [
    json_name = "tenantId",
    (gnostic.openapi.v3.property) = {description: "租户ID"}
  ];
  
  map<string, string> auth_params = 6 [
    json_name = "authParams",
    (gnostic.openapi.v3.property) = {description: "认证参数，用于不同认证类型的额外参数"}
  ];
}

// 简化的用户信息
message UserInfo {
  uint32 id = 1 [
    json_name = "id",
    (gnostic.openapi.v3.property) = {description: "用户ID"}
  ];
  
  string username = 2 [
    json_name = "username",
    (gnostic.openapi.v3.property) = {description: "用户名"}
  ];
  
  string display_name = 3 [
    json_name = "displayName",
    (gnostic.openapi.v3.property) = {description: "显示名称"}
  ];
  
  string email = 4 [
    json_name = "email",
    (gnostic.openapi.v3.property) = {description: "电子邮箱"}
  ];
  
  string avatar = 5 [
    json_name = "avatar",
    (gnostic.openapi.v3.property) = {description: "头像URL"}
  ];
  
  UserAuthority authority = 6 [
    json_name = "authority",
    (gnostic.openapi.v3.property) = {description: "用户权限级别"}
  ];
  
  uint32 tenant_id = 7 [
    json_name = "tenantId",
    (gnostic.openapi.v3.property) = {description: "租户ID"}
  ];
}

// 登录响应
message LoginResponse {
  string access_token = 1 [
    json_name = "accessToken",
    (gnostic.openapi.v3.property) = {description: "访问令牌"}
  ];
  
  string refresh_token = 2 [
    json_name = "refreshToken",
    (gnostic.openapi.v3.property) = {description: "刷新令牌"}
  ];
  
  uint32 expires_in = 3 [
    json_name = "expiresIn",
    (gnostic.openapi.v3.property) = {description: "访问令牌过期时间(秒)"}
  ];
  
  string token_type = 4 [
    json_name = "tokenType",
    (gnostic.openapi.v3.property) = {description: "令牌类型"}
  ];
  
  UserInfo user_info = 5 [
    json_name = "userInfo",
    (gnostic.openapi.v3.property) = {description: "用户信息"}
  ];
  
  string session_id = 6 [
    json_name = "sessionId",
    (gnostic.openapi.v3.property) = {description: "会话ID"}
  ];
}

// 其他请求和响应消息定义...
```

## 3. 批量操作接口示例

```protobuf
// api/user/service/v1/user_batch.proto
syntax = "proto3";

package api.user.service.v1;

import "google/protobuf/empty.proto";
import "gnostic/openapi/v3/annotations.proto";

import "api/user/service/v1/user.proto";
import "api/user/service/v1/error_codes.proto";

option go_package = "om-platform/api/user/service/v1;v1";
option java_multiple_files = true;
option java_package = "api.user.service.v1";

// 用户批量操作服务
// 提供高效的批量用户管理功能
// 错误处理：所有接口在发生错误时将返回ErrorResponse结构，包含统一的错误码、错误消息和详细信息
// 性能说明：批量接口支持大量数据处理，建议单次请求不超过1000条记录
service UserBatchService {
  // 批量创建用户
  // 一次创建多个用户，提高创建效率
  // 权限要求：USER_CREATE 或 ADMIN
  rpc BatchCreateUsers (BatchCreateUsersRequest) returns (BatchCreateUsersResponse) {}
  
  // 批量更新用户
  // 一次更新多个用户信息
  // 权限要求：USER_UPDATE 或 ADMIN
  rpc BatchUpdateUsers (BatchUpdateUsersRequest) returns (BatchUpdateUsersResponse) {}
  
  // 批量删除用户
  // 一次删除多个用户
  // 权限要求：USER_DELETE 或 ADMIN
  rpc BatchDeleteUsers (BatchDeleteUsersRequest) returns (BatchDeleteUsersResponse) {}
  
  // 批量导入用户
  // 从文件导入用户数据
  // 支持CSV、Excel格式
  // 权限要求：USER_IMPORT 或 ADMIN
  rpc ImportUsers (stream ImportUsersRequest) returns (ImportUsersResponse) {}
  
  // 批量导出用户
  // 将用户数据导出为文件
  // 支持CSV、Excel、JSON格式
  // 权限要求：USER_EXPORT 或 ADMIN
  rpc ExportUsers (ExportUsersRequest) returns (stream ExportUsersResponse) {}
}

// 批量创建用户请求
message BatchCreateUsersRequest {
  repeated CreateUserRequest users = 1 [
    json_name = "users",
    (gnostic.openapi.v3.property) = {description: "用户创建请求列表"}
  ];
  
  bool skip_on_error = 2 [
    json_name = "skipOnError",
    (gnostic.openapi.v3.property) = {description: "遇到错误是否跳过继续处理"}
  ];
}

// 批量创建用户响应
message BatchCreateUsersResponse {
  message Result {
    uint32 index = 1; // 请求中的索引位置
    oneof result {
      uint32 user_id = 2; // 成功时返回用户ID
      ErrorResponse error = 3; // 失败时返回错误信息
    }
  }
  
  repeated Result results = 1 [
    json_name = "results",
    (gnostic.openapi.v3.property) = {description: "创建结果列表"}
  ];
  
  uint32 success_count = 2 [
    json_name = "successCount",
    (gnostic.openapi.v3.property) = {description: "成功创建的用户数量"}
  ];
  
  uint32 failure_count = 3 [
    json_name = "failureCount",
    (gnostic.openapi.v3.property) = {description: "创建失败的用户数量"}
  ];
}

// 批量更新用户请求
message BatchUpdateUsersRequest {
  repeated UpdateUserRequest users = 1 [
    json_name = "users",
    (gnostic.openapi.v3.property) = {description: "用户更新请求列表"}
  ];
  
  bool skip_on_error = 2 [
    json_name = "skipOnError",
    (gnostic.openapi.v3.property) = {description: "遇到错误是否跳过继续处理"}
  ];
}

// 批量更新用户响应
message BatchUpdateUsersResponse {
  message Result {
    uint32 index = 1; // 请求中的索引位置
    uint32 user_id = 2; // 用户ID
    oneof result {
      bool success = 3; // 成功标志
      ErrorResponse error = 4; // 失败时返回错误信息
    }
  }
  
  repeated Result results = 1 [
    json_name = "results",
    (gnostic.openapi.v3.property) = {description: "更新结果列表"}
  ];
  
  uint32 success_count = 2 [
    json_name = "successCount",
    (gnostic.openapi.v3.property) = {description: "成功更新的用户数量"}
  ];
  
  uint32 failure_count = 3 [
    json_name = "failureCount",
    (gnostic.openapi.v3.property) = {description: "更新失败的用户数量"}
  ];
}

// 其他批量操作请求和响应消息定义...
```

## 4. 流式API示例

```protobuf
// api/user/service/v1/user_stream.proto
syntax = "proto3";

package api.user.service.v1;

import "google/protobuf/empty.proto";
import "google/protobuf/timestamp.proto";
import "gnostic/openapi/v3/annotations.proto";

import "api/user/service/v1/user.proto";
import "api/user/service/v1/error_codes.proto";

option go_package = "om-platform/api/user/service/v1;v1";
option java_multiple_files = true;
option java_package = "api.user.service.v1";

// 用户流式服务
// 提供流式API，适用于大数据量的用户管理场景
// 错误处理：流式接口会在流中返回错误信息，客户端需要妥善处理
// 性能说明：流式接口适合处理大量数据，无需一次性加载全部数据到内存
service UserStreamService {
  // 流式获取用户列表
  // 适用于需要处理大量用户数据的场景
  // 权限要求：USER_READ 或 ADMIN
  rpc StreamUsers (StreamUsersRequest) returns (stream User) {}
  
  // 流式导出用户数据
  // 将用户数据以流式方式导出
  // 权限要求：USER_EXPORT 或 ADMIN
  rpc StreamExportUsers (StreamExportUsersRequest) returns (stream StreamExportUsersResponse) {}
  
  // 流式导入用户数据
  // 以流式方式导入用户数据
  // 权限要求：USER_IMPORT 或 ADMIN
  rpc StreamImportUsers (stream StreamImportUsersRequest) returns (StreamImportUsersResponse) {}
  
  // 流式用户操作日志
  // 实时获取用户操作日志
  // 权限要求：USER_AUDIT 或 ADMIN
  rpc StreamUserAuditLogs (StreamUserAuditLogsRequest) returns (stream UserAuditLog) {}
}

// 流式获取用户请求
message StreamUsersRequest {
  // 过滤条件
  message Filter {
    repeated uint32 user_ids = 1; // 用户ID列表
    repeated uint32 department_ids = 2; // 部门ID列表
    repeated uint32 role_ids = 3; // 角色ID列表
    string username_pattern = 4; // 用户名模式
    string email_pattern = 5; // 邮箱模式
    google.protobuf.Timestamp created_after = 6; // 创建时间起始
    google.protobuf.Timestamp created_before = 7; // 创建时间截止
  }
  
  Filter filter = 1; // 过滤条件
  uint32 batch_size = 2; // 批次大小
  string sort_by = 3; // 排序字段
  bool ascending = 4; // 是否升序
}

// 流式导出用户数据请求
message StreamExportUsersRequest {
  StreamUsersRequest.Filter filter = 1; // 过滤条件
  string export_format = 2; // 导出格式(JSON, CSV, EXCEL)
  repeated string fields = 3; // 要导出的字段
}

// 流式导出用户数据响应
message StreamExportUsersResponse {
  oneof data {
    bytes chunk = 1; // 数据块
    ExportMetadata metadata = 2; // 元数据
  }
  
  message ExportMetadata {
    string format = 1; // 格式
    uint32 total_records = 2; // 总记录数
    uint32 total_chunks = 3; // 总块数
    string file_name = 4; // 文件名
  }
}

// 流式导入用户数据请求
message StreamImportUsersRequest {
  oneof data {
    bytes chunk = 1; // 数据块
    ImportMetadata metadata = 2; // 元数据
  }
  
  message ImportMetadata {
    string format = 1; // 格式
    uint32 total_records = 2; // 总记录数
    uint32 total_chunks = 3; // 总块数
    string file_name = 4; // 文件名
    bool skip_on_error = 5; // 遇到错误是否跳过
    bool update_existing = 6; // 是否更新已存在的用户
  }
}

// 流式导入用户数据响应
message StreamImportUsersResponse {
  uint32 total_processed = 1; // 总处理记录数
  uint32 success_count = 2; // 成功记录数
  uint32 failure_count = 3; // 失败记录数
  repeated ImportError errors = 4; // 导入错误
  
  message ImportError {
    uint32 row = 1; // 行号
    string field = 2; // 字段
    string message = 3; // 错误消息
  }
}

// 用户审计日志
message UserAuditLog {
  string id = 1; // 日志ID
  uint32 user_id = 2; // 用户ID
  string action = 3; // 操作
  string resource_type = 4; // 资源类型
  string resource_id = 5; // 资源ID
  google.protobuf.Timestamp timestamp = 6; // 时间戳
  string ip_address = 7; // IP地址
  string user_agent = 8; // 用户代理
  map<string, string> details = 9; // 详细信息
}

// 流式用户操作日志请求
message StreamUserAuditLogsRequest {
  uint32 user_id = 1; // 用户ID
  google.protobuf.Timestamp start_time = 2; // 开始时间
  google.protobuf.Timestamp end_time = 3; // 结束时间
  repeated string actions = 4; // 操作类型
  bool real_time = 5; // 是否实时订阅新日志
}
```

## 5. API限流与熔断定义示例

```protobuf
// api/user/service/v1/api_control.proto
syntax = "proto3";

package api.user.service.v1;

import "google/protobuf/descriptor.proto";

option go_package = "om-platform/api/user/service/v1;v1";
option java_multiple_files = true;
option java_package = "api.user.service.v1";

// API控制选项
extend google.protobuf.MethodOptions {
  // 限流配置
  RateLimitOptions rate_limit = 50001;
  
  // 熔断配置
  CircuitBreakerOptions circuit_breaker = 50002;
  
  // 缓存配置
  CacheOptions cache = 50003;
  
  // 超时配置
  TimeoutOptions timeout = 50004;
  
  // 重试配置
  RetryOptions retry = 50005;
  
  // 追踪配置
  TracingOptions tracing = 50006;
}

// 限流配置
message RateLimitOptions {
  // 每分钟请求数限制
  uint32 requests_per_minute = 1;
  
  // 突发请求数限制
  uint32 burst = 2;
  
  // 限流策略
  enum Strategy {
    STRATEGY_UNSPECIFIED = 0; // 未指定
    STRATEGY_REJECT = 1; // 拒绝请求
    STRATEGY_WAIT = 2; // 等待
    STRATEGY_ADAPTIVE = 3; // 自适应
  }
  Strategy strategy = 3;
  
  // 限流级别
  enum Level {
    LEVEL_UNSPECIFIED = 0; // 未指定
    LEVEL_GLOBAL = 1; // 全局
    LEVEL_IP = 2; // 按IP
    LEVEL_USER = 3; // 按用户
    LEVEL_TENANT = 4; // 按租户
  }
  Level level = 4;
}

// 熔断配置
message CircuitBreakerOptions {
  // 错误阈值百分比
  uint32 error_threshold_percentage = 1;
  
  // 最小请求数
  uint32 min_request_amount = 2;
  
  // 熔断后休眠时间(毫秒)
  uint32 sleep_window_ms = 3;
  
  // 熔断策略
  enum Strategy {
    STRATEGY_UNSPECIFIED = 0; // 未指定
    STRATEGY_FAIL_FAST = 1; // 快速失败
    STRATEGY_FALLBACK = 2; // 降级
  }
  Strategy strategy = 4;
}

// 缓存配置
message CacheOptions {
  // 是否可缓存
  bool cacheable = 1;
  
  // 缓存有效期(秒)
  uint32 ttl_seconds = 2;
  
  // 缓存键字段
  repeated string cache_keys = 3;
  
  // 缓存策略
  enum Strategy {
    STRATEGY_UNSPECIFIED = 0; // 未指定
    STRATEGY_LOCAL = 1; // 本地缓存
    STRATEGY_DISTRIBUTED = 2; // 分布式缓存
    STRATEGY_LAYERED = 3; // 多级缓存
  }
  Strategy strategy = 4;
}

// 超时配置
message TimeoutOptions {
  // 超时时间(毫秒)
  uint32 timeout_ms = 1;
}

// 重试配置
message RetryOptions {
  // 最大重试次数
  uint32 max_retries = 1;
  
  // 初始重试延迟(毫秒)
  uint32 initial_delay_ms = 2;
  
  // 最大重试延迟(毫秒)
  uint32 max_delay_ms = 3;
  
  // 重试策略
  enum Strategy {
    STRATEGY_UNSPECIFIED = 0; // 未指定
    STRATEGY_CONSTANT = 1; // 固定间隔
    STRATEGY_LINEAR = 2; // 线性增长
    STRATEGY_EXPONENTIAL = 3; // 指数增长
  }
  Strategy strategy = 4;
}

// 追踪配置
message TracingOptions {
  // 是否启用追踪
  bool enabled = 1;
  
  // 采样率(0.0-1.0)
  float sampling_rate = 2;
  
  // 追踪标签
  map<string, string> tags = 3;
}
```

## 6. 字段筛选示例

```protobuf
// 使用FieldMask实现字段筛选的请求示例
message GetUserRequest {
  uint32 id = 1 [
    json_name = "id",
    (gnostic.openapi.v3.property) = {description: "用户ID"}
  ];
  
  google.protobuf.FieldMask field_mask = 2 [
    json_name = "fieldMask",
    (gnostic.openapi.v3.property) = {description: "指定需要返回的字段，格式如：username,email,roles"}
  ];
}

// 使用示例：
// 请求：
// {
//   "id": 123,
//   "fieldMask": {
//     "paths": ["username", "email", "roles"]
//   }
// }
//
// 响应将只包含指定的字段：
// {
//   "username": "john.doe",
//   "email": "john.doe@example.com",
//   "roles": [1, 2, 3]
// }
```

## 7. 扩展字段示例

```protobuf
// 在关键消息结构中添加扩展字段
message User {
  // 现有字段...
  optional uint32 id = 1;
  optional string username = 2;
  optional string email = 3;
  // ...
  
  // 扩展字段，用于未来功能扩展
  google.protobuf.Struct extensions = 1000 [
    json_name = "extensions",
    (gnostic.openapi.v3.property) = {description: "扩展字段，用于存储自定义属性"}
  ];
}

// 使用示例：
// {
//   "id": 123,
//   "username": "john.doe",
//   "email": "john.doe@example.com",
//   "extensions": {
//     "employee_id": "EMP001",
//     "skills": ["Go", "Kubernetes", "gRPC"],
//     "preferences": {
//       "theme": "dark",
//       "notifications": true
//     }
//   }
// }
```

这些示例代码展示了API优化方案中提出的各项建议的具体实现方式，可以作为开发团队实施优化的参考。