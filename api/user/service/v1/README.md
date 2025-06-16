# 用户服务API文档

## 概述

本目录包含用户管理平台的核心API定义，基于Protocol Buffers 3.0规范，提供完整的用户生命周期管理功能。

## 服务架构

### 核心服务

- **UserService**: 用户基础信息管理（增删改查、批量操作）
- **AccountService**: 账户状态管理（锁定、解锁、状态变更）
- **AuthenticationService**: 认证服务（登录、登出、令牌管理）
- **RoleService**: 角色权限管理
- **PermissionService**: 权限管理
- **OrganizationService**: 组织架构管理
- **DepartmentService**: 部门管理
- **TenantService**: 多租户管理
- **IdentityProviderService**: 身份提供商管理
- **AuditService**: 审计日志服务

## 设计规范

### 1. 包结构与Option规范

所有proto文件统一使用以下配置：

```protobuf
package api.user.service.v1;

// Go语言包路径配置
option go_package = "github.com/Sakuya1998/om-platform/api/user/service/v1;userv1";
// Java包路径配置
option java_package = "com.omplatform.api.user.service.v1";
// Java多文件生成配置
option java_multiple_files = true;
// C#命名空间配置
option csharp_namespace = "OmPlatform.Api.User.Service.V1";
// PHP命名空间配置
option php_namespace = "OmPlatform\\Api\\User\\Service\\V1";
// Ruby包配置
option ruby_package = "OmPlatform::Api::User::Service::V1";
```

### 2. 通用消息结构

#### 分页请求 (PagingRequest)

```protobuf
message PagingRequest {
  optional int32 page = 1;      // 页码，从1开始
  optional int32 page_size = 2; // 每页记录数，默认20，最大100
  optional int32 offset = 3;    // 偏移量（可选，与page互斥）
  optional int32 limit = 4;     // 限制数量（可选，与page_size互斥）
}
```

#### 分页响应 (PagingResponse)

```protobuf
message PagingResponse {
  int32 total_count = 1;        // 总记录数
  int32 page = 2;               // 当前页码
  int32 page_size = 3;          // 每页记录数
  int32 total_pages = 4;        // 总页数
  bool has_next = 5;            // 是否有下一页
  bool has_previous = 6;        // 是否有上一页
}
```

#### 错误响应 (ErrorResponse)

```protobuf
message ErrorResponse {
  ErrorCode code = 1;                    // 错误码
  string message = 2;                    // 错误消息
  string details = 3;                    // 详细信息
  map<string, string> i18n_messages = 4; // 国际化消息
  google.protobuf.Timestamp timestamp = 5; // 错误时间戳
  string trace_id = 6;                   // 链路追踪ID
}
```

#### 批量操作结果 (BatchOperationResult)

```protobuf
message BatchOperationResult {
  int32 total_count = 1;        // 总操作数
  int32 success_count = 2;      // 成功数
  int32 failed_count = 3;       // 失败数
  int32 skipped_count = 4;      // 跳过数
  repeated ErrorResponse errors = 5;     // 错误详情
  google.protobuf.Any original_data = 6; // 原始数据
  repeated string success_ids = 7;       // 成功处理的ID列表
}
```

### 3. 错误码分段规范

错误码采用分段设计，便于分类管理：

```protobuf
enum ErrorCode {
  // 通用错误码 (0-999)
  ERROR_CODE_UNSPECIFIED = 0;
  ERROR_CODE_INTERNAL_ERROR = 1;
  ERROR_CODE_INVALID_ARGUMENT = 2;
  ERROR_CODE_PERMISSION_DENIED = 3;
  ERROR_CODE_NOT_FOUND = 4;
  ERROR_CODE_ALREADY_EXISTS = 5;
  ERROR_CODE_RESOURCE_EXHAUSTED = 6;
  ERROR_CODE_FAILED_PRECONDITION = 7;
  ERROR_CODE_ABORTED = 8;
  ERROR_CODE_OUT_OF_RANGE = 9;
  ERROR_CODE_UNIMPLEMENTED = 10;
  ERROR_CODE_UNAVAILABLE = 11;
  ERROR_CODE_DATA_LOSS = 12;
  ERROR_CODE_UNAUTHENTICATED = 13;
  
  // 用户相关错误码 (1000-1999)
  ERROR_CODE_USER_NOT_FOUND = 1000;
  ERROR_CODE_USER_ALREADY_EXISTS = 1001;
  ERROR_CODE_USER_INVALID_USERNAME = 1002;
  ERROR_CODE_USER_INVALID_EMAIL = 1003;
  ERROR_CODE_USER_INVALID_PHONE = 1004;
  ERROR_CODE_USER_INVALID_PASSWORD = 1005;
  
  // 认证相关错误码 (2000-2999)
  ERROR_CODE_AUTH_INVALID_CREDENTIALS = 2000;
  ERROR_CODE_AUTH_TOKEN_EXPIRED = 2001;
  ERROR_CODE_AUTH_TOKEN_INVALID = 2002;
  ERROR_CODE_AUTH_CAPTCHA_REQUIRED = 2003;
  ERROR_CODE_AUTH_CAPTCHA_INVALID = 2004;
  ERROR_CODE_AUTH_MFA_REQUIRED = 2005;
  ERROR_CODE_AUTH_MFA_INVALID = 2006;
  
  // 权限相关错误码 (3000-3999)
  ERROR_CODE_PERMISSION_INSUFFICIENT = 3000;
  ERROR_CODE_PERMISSION_ROLE_NOT_FOUND = 3001;
  ERROR_CODE_PERMISSION_INVALID_SCOPE = 3002;
  
  // 组织相关错误码 (4000-4999)
  ERROR_CODE_ORG_NOT_FOUND = 4000;
  ERROR_CODE_ORG_INVALID_HIERARCHY = 4001;
  ERROR_CODE_DEPT_NOT_FOUND = 4002;
  
  // 租户相关错误码 (5000-5999)
  ERROR_CODE_TENANT_NOT_FOUND = 5000;
  ERROR_CODE_TENANT_QUOTA_EXCEEDED = 5001;
  ERROR_CODE_TENANT_SUSPENDED = 5002;
  
  // 会话相关错误码 (6000-6999)
  ERROR_CODE_SESSION_EXPIRED = 6000;
  ERROR_CODE_SESSION_INVALID = 6001;
  ERROR_CODE_SESSION_CONCURRENT_LIMIT = 6002;
  
  // 身份提供商相关错误码 (7000-7999)
  ERROR_CODE_IDP_NOT_FOUND = 7000;
  ERROR_CODE_IDP_CONFIG_INVALID = 7001;
  ERROR_CODE_IDP_CONNECTION_FAILED = 7002;
  
  // 限流与熔断相关错误码 (8000-8999)
  ERROR_CODE_RATE_LIMIT_EXCEEDED = 8000;
  ERROR_CODE_CIRCUIT_BREAKER_OPEN = 8001;
  ERROR_CODE_QUOTA_EXCEEDED = 8002;
}
```

### 4. 通用枚举定义

#### 用户账户状态

```protobuf
enum UserAccountStatus {
  USER_ACCOUNT_STATUS_UNSPECIFIED = 0; // 未指定
  USER_ACCOUNT_STATUS_ACTIVE = 1;      // 激活
  USER_ACCOUNT_STATUS_INACTIVE = 2;    // 未激活
  USER_ACCOUNT_STATUS_LOCKED = 3;      // 锁定
  USER_ACCOUNT_STATUS_SUSPENDED = 4;   // 暂停
  USER_ACCOUNT_STATUS_DELETED = 5;     // 已删除
  USER_ACCOUNT_STATUS_PENDING = 6;     // 待审核
}
```

#### 认证类型

```protobuf
enum AuthenticationType {
  AUTHENTICATION_TYPE_UNSPECIFIED = 0; // 未指定
  AUTHENTICATION_TYPE_PASSWORD = 1;    // 用户名密码认证
  AUTHENTICATION_TYPE_SMS = 2;         // 手机短信验证码
  AUTHENTICATION_TYPE_EMAIL = 3;       // 邮箱验证码
  AUTHENTICATION_TYPE_LDAP = 4;        // LDAP认证
  AUTHENTICATION_TYPE_OAUTH2 = 5;      // OAuth2认证
  AUTHENTICATION_TYPE_SAML = 6;        // SAML认证
  AUTHENTICATION_TYPE_MFA = 7;         // 多因素认证
  AUTHENTICATION_TYPE_BIOMETRIC = 8;   // 生物识别认证
  AUTHENTICATION_TYPE_API_KEY = 9;     // API密钥认证
}
```

### 5. 自定义选项扩展

#### 限流选项 (RateLimitOptions)

```protobuf
message RateLimitOptions {
  int32 requests_per_minute = 1;  // 每分钟请求数限制
  int32 requests_per_hour = 2;    // 每小时请求数限制
  int32 requests_per_day = 3;     // 每天请求数限制
  int32 burst_size = 4;           // 突发请求大小
  string key_pattern = 5;         // 限流键模式
}

extend google.protobuf.MethodOptions {
  RateLimitOptions rate_limit = 50001; // 字段号：50001
}
```

#### 缓存选项 (CacheOptions)

```protobuf
message CacheOptions {
  int32 ttl_seconds = 1;          // 缓存生存时间（秒）
  string cache_key_pattern = 2;   // 缓存键模式
  bool enable_cache = 3;          // 是否启用缓存
  repeated string invalidate_tags = 4; // 缓存失效标签
}

extend google.protobuf.MethodOptions {
  CacheOptions cache = 50002; // 字段号：50002
}
```

#### 熔断选项 (CircuitBreakerOptions)

```protobuf
message CircuitBreakerOptions {
  int32 failure_threshold = 1;    // 失败阈值
  int32 timeout_seconds = 2;      // 超时时间（秒）
  int32 recovery_timeout = 3;     // 恢复超时时间（秒）
  bool enable_circuit_breaker = 4; // 是否启用熔断
}

extend google.protobuf.MethodOptions {
  CircuitBreakerOptions circuit_breaker = 50003; // 字段号：50003
}
```

## 使用示例

### 1. 用户列表查询

```json
{
  "paging": {
    "page": 1,
    "pageSize": 20
  },
  "sortOptions": [
    {
      "field": "created_at",
      "direction": "SORT_DIRECTION_DESC"
    }
  ],
  "filter": {
    "status": ["USER_ACCOUNT_STATUS_ACTIVE"],
    "createdAtStart": "2024-01-01T00:00:00Z",
    "createdAtEnd": "2024-12-31T23:59:59Z"
  }
}
```

### 2. 批量用户创建

```json
{
  "users": [
    {
      "username": "user1",
      "email": "user1@example.com",
      "displayName": "用户一"
    },
    {
      "username": "user2",
      "email": "user2@example.com",
      "displayName": "用户二"
    }
  ],
  "skipDuplicates": true
}
```

### 3. 用户登录

```json
{
  "username": "admin@example.com",
  "password": "password123",
  "authType": "AUTHENTICATION_TYPE_PASSWORD",
  "deviceInfo": {
    "deviceType": "DEVICE_TYPE_BROWSER",
    "deviceName": "Chrome 120.0",
    "ipAddress": "192.168.1.100",
    "userAgent": "Mozilla/5.0..."
  },
  "rememberMe": true
}
```

## 兼容性说明

### 向前兼容性

- 所有新增字段使用`optional`或`repeated`修饰符
- 枚举值只能追加，不能删除或修改现有值
- 消息字段号不能重复使用
- 服务方法只能新增，不能删除或修改签名

### 国际化支持

- 错误消息支持多语言，通过`i18n_messages`字段提供
- 枚举值描述支持国际化
- 字段描述支持多语言注释

### 扩展性设计

- 预留扩展字段`extra_attributes`用于存储JSON格式的自定义属性
- 支持标签系统用于灵活分类和搜索
- 审计信息统一管理，支持完整的操作追踪

## 开发指南

### 代码生成

```bash
# Go代码生成
protoc --go_out=. --go-grpc_out=. *.proto

# Java代码生成
protoc --java_out=. *.proto

# TypeScript代码生成
protoc --ts_out=. *.proto
```

### 测试建议

1. **单元测试**: 针对每个服务方法编写单元测试
2. **集成测试**: 测试服务间的交互和数据一致性
3. **性能测试**: 验证分页、批量操作的性能表现
4. **兼容性测试**: 确保API版本间的兼容性

### 监控指标

- API调用次数和响应时间
- 错误率和错误类型分布
- 限流和熔断触发情况
- 缓存命中率
- 用户活跃度和登录成功率

## 更新日志

### v1.0.0 (2024-01-01)
- 初始版本发布
- 完整的用户管理API
- 支持多租户架构
- 集成认证和权限管理

---

**注意**: 本文档随API版本更新，请关注版本变更说明。如有疑问，请联系开发团队。