# 智能运维平台用户服务 API

## 概述

用户服务 API 是智能运维平台的核心用户管理服务，采用现代化的 API 设计模式，提供统一的错误处理、性能优化和扩展性支持。服务实现了 HTTP/gRPC 双协议支持，提供完整的 RESTful API 映射，便于不同场景下的集成和调用。

## 主要特性

### 1. 统一错误处理
- 引入统一的 `ErrorResponse` 结构
- 标准化错误码枚举 `ErrorCode`
- 支持多语言错误消息
- 详细的错误上下文信息

### 2. 性能优化
- 内置限流配置 (`RateLimitOptions`)，基于每秒请求数
- 精细化的突发流量控制 (`burst` 参数)
- 基于资源路径的缓存策略 (`key_pattern`)
- 熔断器模式 (`CircuitBreakerOptions`)
- 字段筛选和分页优化
- 针对读写操作的差异化限流策略

### 3. API 设计
- RESTful 风格的 HTTP 映射（所有RPC方法均支持HTTP调用）
- 统一的请求/响应结构
- 支持批量操作（统一的`/batch`端点）
- 流式数据处理
- 扩展字段支持
- 标准化的资源路径（如`/v1/users/{user_id}/roles`）
- 符合HTTP语义的方法映射（GET/POST/PUT/DELETE）

### 4. 功能增强
- 完整的审计日志
- 多租户支持
- 身份提供商集成
- 高级权限管理
- 数据导入导出

## HTTP API 端点映射

所有 gRPC 服务方法都已配置了对应的 HTTP 端点，支持 RESTful 风格的 API 调用：

### 用户管理 API
- `GET /v1/users` - 获取用户列表
- `GET /v1/users/{user_id}` - 获取单个用户
- `POST /v1/users` - 创建用户
- `PUT /v1/users/{user_id}` - 更新用户
- `DELETE /v1/users/{user_id}` - 删除用户
- `POST /v1/users/batch` - 批量创建用户
- `PUT /v1/users/batch` - 批量更新用户
- `DELETE /v1/users/batch` - 批量删除用户

### 认证管理 API
- `POST /v1/auth/login` - 用户登录
- `POST /v1/auth/logout` - 用户登出
- `POST /v1/auth/refresh` - 刷新令牌
- `GET /v1/auth/sessions/{session_id}` - 获取会话信息
- `DELETE /v1/auth/sessions/{session_id}` - 删除会话

### 权限管理 API
- `GET /v1/roles` - 获取角色列表
- `GET /v1/roles/{role_id}` - 获取角色详情
- `POST /v1/roles` - 创建角色
- `PUT /v1/roles/{role_id}` - 更新角色
- `DELETE /v1/roles/{role_id}` - 删除角色
- `GET /v1/users/{user_id}/roles` - 获取用户角色
- `PUT /v1/users/{user_id}/roles` - 分配用户角色
- `DELETE /v1/users/{user_id}/roles` - 撤销用户角色

### 权限检查 API
- `GET /v1/permissions` - 获取权限列表
- `POST /v1/permissions/check` - 检查权限
- `POST /v1/permissions/batch-check` - 批量检查权限
- `GET /v1/users/{user_id}/permissions` - 获取用户权限

## 服务列表

### 核心服务

1. **UserService** (`user.proto`)
   - 用户 CRUD 操作
   - 批量用户管理
   - 用户角色和权限
   - 流式数据处理

2. **AuthenticationService** (`authentication.proto`)
   - 用户认证
   - 会话管理
   - 双因素认证
   - 登录历史

3. **AccountService** (`account.proto`)
   - 账户状态管理
   - 安全设置
   - 密码策略
   - 第三方账号绑定

### 组织架构服务

4. **OrganizationService** (`organization.proto`)
   - 组织管理
   - 组织结构
   - 成员管理
   - 批量操作

5. **DepartmentService** (`department.proto`)
   - 部门管理
   - 部门层级
   - 成员分配
   - 数据导入导出

### 权限管理服务

6. **RoleService** (`role.proto`)
   - 角色定义
   - 角色继承
   - 权限分配
   - 批量管理

7. **PermissionService** (`permission.proto`)
   - 权限定义
   - 权限树管理
   - 资源权限
   - 权限检查

### 多租户和集成服务

8. **TenantService** (`tenant.proto`)
   - 租户管理
   - 租户配置
   - 资源隔离
   - 统计分析

9. **IdentityProviderService** (`identity_provider.proto`)
   - 身份提供商管理
   - LDAP/SAML/OIDC 集成
   - 用户同步
   - 连接测试

### 通用组件

10. **Common** (`common.proto`)
    - 通用枚举和结构
    - 分页和排序
    - 审计信息
    - 批量操作结果

11. **Error Codes** (`error_codes.proto`)
    - 统一错误码
    - 错误响应结构
    - 扩展选项定义

## API 设计特性

### 1. 错误处理

```protobuf
message CreateUserResponse {
  oneof result {
    User user = 1;
    ErrorResponse error = 2;
  }
}
```

### 2. 分页查询

```protobuf
message ListUsersRequest {
  PagingRequest paging = 1;
  repeated SortOption sort_options = 2;
  repeated string fields = 3; // 字段筛选
}

message ListUsersResponse {
  oneof result {
    UsersData data = 1;
    ErrorResponse error = 2;
  }
}

message UsersData {
  repeated User users = 1;
  PaginatedResponse pagination = 2;
}
```

### 3. 批量操作

```protobuf
// 批量创建用户
rpc BatchCreateUsers(BatchCreateUsersRequest) returns (BatchCreateUsersResponse);

// 批量更新用户
rpc BatchUpdateUsers(BatchUpdateUsersRequest) returns (BatchUpdateUsersResponse);

// 批量删除用户
rpc BatchDeleteUsers(BatchDeleteUsersRequest) returns (BatchDeleteUsersResponse);
```

### 4. 字段筛选和更新掩码

**字段筛选:**
```protobuf
message GetUserRequest {
  string user_id = 1;
  repeated string fields = 2; // 只返回指定字段
}
```

**更新掩码:**
```protobuf
message UpdateUserRequest {
  string user_id = 1;
  User user = 2;
  repeated string update_mask = 3; // 只更新指定字段
}
```

## 协议支持

### HTTP/gRPC 双协议支持

所有服务同时支持 HTTP 和 gRPC 协议：

- **gRPC**: 高性能的二进制协议，适用于服务间通信
- **HTTP**: RESTful API，适用于前端和第三方集成
- **统一配置**: 限流、缓存、熔断器等配置对两种协议均生效

### 配置优化亮点

1. **限流策略升级**
   - 从 `requests_per_minute` 升级为 `requests_per_second`
   - 支持突发流量控制 (`burst` 参数)
   - 针对不同操作类型的差异化限流

2. **缓存策略优化**
   - 从简单的 `cacheable` 标记升级为 `key_pattern` 模式
   - 支持动态缓存键生成
   - 细粒度的缓存控制

3. **批量操作支持**
   - 统一的 `/batch` 端点设计
   - 优化的批量操作限流策略
   - 事务性批量处理

## 性能优化配置

### 1. 限流配置

```protobuf
option (rate_limit) = {
  requests_per_second: 100
  burst: 200
};
```

### 2. 缓存配置

```protobuf
option (cache) = {
  ttl_seconds: 300
  key_pattern: "user:{user_id}"
  max_age_seconds: 60
};
```

### 3. 熔断器配置

```protobuf
option (circuit_breaker) = {
  failure_threshold: 5
  timeout_seconds: 30
  recovery_timeout_seconds: 60
};
```

## 最佳实践

### 1. 错误处理

```go
// 客户端错误处理示例
resp, err := client.GetUser(ctx, &user.GetUserRequest{
    UserId: "user123",
})
if err != nil {
    return err
}

switch result := resp.Result.(type) {
case *user.GetUserResponse_User:
    // 处理成功响应
    user := result.User
    // ...
case *user.GetUserResponse_Error:
    // 处理业务错误
    errorResp := result.Error
    log.Errorf("业务错误: %s (代码: %s)", errorResp.Message, errorResp.Code)
    return fmt.Errorf("获取用户失败: %s", errorResp.Message)
}
```

### 2. 分页查询

```go
// 分页查询示例
resp, err := client.ListUsers(ctx, &user.ListUsersRequest{
    Paging: &user.PagingRequest{
        Page: 1,
        Size: 20,
    },
    SortOptions: []*user.SortOption{
        {
            Field: "created_at",
            Direction: user.SortDirection_SORT_DIRECTION_DESC,
        },
    },
    Fields: []string{"user_id", "username", "email", "status"}, // 字段筛选
})
```

### 3. 批量操作

```go
// 批量创建用户示例
resp, err := client.BatchCreateUsers(ctx, &user.BatchCreateUsersRequest{
    Users: []*user.CreateUserRequest{
        {Username: "user1", Email: "user1@example.com"},
        {Username: "user2", Email: "user2@example.com"},
    },
})

if err != nil {
    return err
}

switch result := resp.Result.(type) {
case *user.BatchCreateUsersResponse_Results:
    results := result.Results
    for i, result := range results.Results {
        if result.Success {
            log.Infof("用户 %d 创建成功", i)
        } else {
            log.Errorf("用户 %d 创建失败: %s", i, result.Error)
        }
    }
case *user.BatchCreateUsersResponse_Error:
    return fmt.Errorf("批量创建失败: %s", result.Error.Message)
}
```

### 4. HTTP API 调用示例

```bash
# 获取用户列表
curl -X GET "http://localhost:8080/v1/users?page=1&size=10" \
  -H "Authorization: Bearer <token>"

# 创建用户
curl -X POST "http://localhost:8080/v1/users" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "username": "newuser",
    "email": "newuser@example.com",
    "password": "securepassword"
  }'

# 批量创建用户
curl -X POST "http://localhost:8080/v1/users/batch" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "users": [
      {"username": "user1", "email": "user1@example.com"},
      {"username": "user2", "email": "user2@example.com"}
    ]
  }'

# 检查用户权限
curl -X POST "http://localhost:8080/v1/users/123/permissions/check" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "resource": "user:read",
    "action": "view"
  }'
```

## 版本说明

### 版本策略

- v1 是 API 的初始稳定版本
- 采用语义化版本控制 (Semantic Versioning)
- 保证向后兼容性

### 未来版本计划

1. **小版本更新**: 添加新功能，保持向后兼容
2. **补丁版本**: 错误修复和性能优化
3. **重大版本**: 仅在必要时进行不兼容更新

## 开发指南

### 1. 代码生成

```bash
# 生成 Go 代码
protoc --go_out=. --go-grpc_out=. \
  --go_opt=paths=source_relative \
  --go-grpc_opt=paths=source_relative \
  api/user/service/v1/*.proto

# 生成 HTTP 网关代码（支持 RESTful API）
protoc --grpc-gateway_out=. \
  --grpc-gateway_opt=paths=source_relative \
  --grpc-gateway_opt=generate_unbound_methods=true \
  api/user/service/v1/*.proto

# 生成 OpenAPI 文档
protoc --openapiv2_out=. \
  --openapiv2_opt=allow_merge=true \
  api/user/service/v1/*.proto
```

### 2. 服务实现

```go
type userServiceServer struct {
    user.UnimplementedUserServiceServer
    // 依赖注入
}

func (s *userServiceServer) CreateUser(ctx context.Context, req *user.CreateUserRequest) (*user.CreateUserResponse, error) {
    // 参数验证
    if err := req.Validate(); err != nil {
        return &user.CreateUserResponse{
            Result: &user.CreateUserResponse_Error{
                Error: &user.ErrorResponse{
                    Code: user.ErrorCode_ERROR_CODE_INVALID_ARGUMENT,
                    Message: "参数验证失败",
                    Details: err.Error(),
                },
            },
        }, nil
    }

    // 业务逻辑
    user, err := s.createUser(ctx, req)
    if err != nil {
        return &user.CreateUserResponse{
            Result: &user.CreateUserResponse_Error{
                Error: &user.ErrorResponse{
                    Code: user.ErrorCode_ERROR_CODE_INTERNAL_ERROR,
                    Message: "创建用户失败",
                    Details: err.Error(),
                },
            },
        }, nil
    }

    return &user.CreateUserResponse{
        Result: &user.CreateUserResponse_User{
            User: user,
        },
    }, nil
}
```

### 3. 中间件集成

```go
// 限流中间件
func RateLimitInterceptor() grpc.UnaryServerInterceptor {
    return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
        // 从方法选项中获取限流配置
        // 实现限流逻辑
        return handler(ctx, req)
    }
}

// 缓存中间件
func CacheInterceptor() grpc.UnaryServerInterceptor {
    return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
        // 从方法选项中获取缓存配置
        // 实现缓存逻辑
        return handler(ctx, req)
    }
}
```

## 监控和观测

### 1. 指标收集

- API 调用次数和延迟
- 错误率统计
- 限流触发次数
- 缓存命中率
- 熔断器状态

### 2. 日志记录

- 结构化日志
- 请求追踪 ID
- 审计日志
- 性能日志

### 3. 分布式追踪

- OpenTelemetry 集成
- 跨服务调用追踪
- 性能瓶颈分析

## 安全考虑

### 1. 认证和授权

- JWT Token 验证
- RBAC 权限控制
- API 密钥管理
- OAuth 2.0 集成

### 2. 数据保护

- 敏感数据加密
- PII 数据脱敏
- 审计日志
- 数据备份和恢复

### 3. 网络安全

- TLS 加密传输
- API 网关保护
- DDoS 防护
- 入侵检测

## 故障排查

### 1. 常见问题

**问题**: 限流错误
```
ERROR_CODE_RATE_LIMIT_EXCEEDED: 请求频率超过限制
```
**解决**: 检查客户端请求频率，实现指数退避重试

**问题**: 缓存失效
```
ERROR_CODE_CACHE_MISS: 缓存未命中
```
**解决**: 检查缓存配置和 TTL 设置

**问题**: 熔断器开启
```
ERROR_CODE_CIRCUIT_BREAKER_OPEN: 熔断器已开启
```
**解决**: 检查下游服务状态，等待熔断器恢复

### 2. 调试工具

- gRPC 反射服务
- Postman/Insomnia 测试
- grpcurl 命令行工具
- 分布式追踪系统

## 贡献指南

### 1. 开发流程

1. Fork 项目
2. 创建功能分支
3. 编写代码和测试
4. 提交 Pull Request
5. 代码审查
6. 合并到主分支

### 2. 代码规范

- 遵循 Protocol Buffers 风格指南
- 使用有意义的命名
- 添加详细的注释
- 编写单元测试
- 更新文档

### 3. 测试要求

- 单元测试覆盖率 > 80%
- 集成测试
- 性能测试
- 安全测试

## 云原生部署

### 1. Kubernetes 部署配置

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: user-service
  namespace: om-platform
spec:
  replicas: 3
  selector:
    matchLabels:
      app: user-service
  template:
    metadata:
      labels:
        app: user-service
      annotations:
        prometheus.io/scrape: "true"
        prometheus.io/port: "9090"
    spec:
      containers:
      - name: user-service
        image: om-platform/user-service:v1.0.0
        ports:
        - containerPort: 8080
          name: http
        - containerPort: 9090
          name: grpc
        - containerPort: 9091
          name: metrics
        resources:
          requests:
            cpu: 500m
            memory: 512Mi
          limits:
            cpu: 1000m
            memory: 1Gi
        readinessProbe:
          httpGet:
            path: /health/ready
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 5
        livenessProbe:
          httpGet:
            path: /health/live
            port: 8080
          initialDelaySeconds: 20
          periodSeconds: 10
        env:
        - name: DB_HOST
          valueFrom:
            configMapKeyRef:
              name: user-service-config
              key: db_host
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: user-service-secrets
              key: db_password
        volumeMounts:
        - name: config-volume
          mountPath: /etc/user-service
      volumes:
      - name: config-volume
        configMap:
          name: user-service-config
```

### 2. 服务网格集成

```yaml
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: user-service
  namespace: om-platform
spec:
  hosts:
  - "api.om-platform.com"
  gateways:
  - om-platform-gateway
  http:
  - match:
    - uri:
        prefix: /v1/users
    - uri:
        prefix: /v1/auth
    route:
    - destination:
        host: user-service
        port:
          number: 8080
    retries:
      attempts: 3
      perTryTimeout: 2s
    fault:
      delay:
        percentage:
          value: 0.1
        fixedDelay: 5s
    timeout: 10s
```

### 3. 自动伸缩配置

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: user-service-hpa
  namespace: om-platform
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: user-service
  minReplicas: 3
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
  - type: Pods
    pods:
      metric:
        name: requests_per_second
      target:
        type: AverageValue
        averageValue: 1000
  behavior:
    scaleDown:
      stabilizationWindowSeconds: 300
    scaleUp:
      stabilizationWindowSeconds: 60
```

## 性能调优指南

### 1. 数据库优化

- 使用连接池管理数据库连接
- 实现读写分离策略
- 针对热点数据实现多级缓存
- 定期执行EXPLAIN分析SQL查询性能

```go
// 数据库连接池配置示例
func initDBPool(ctx context.Context) (*pgxpool.Pool, error) {
    config, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
    if err != nil {
        return nil, fmt.Errorf("解析数据库配置失败: %w", err)
    }
    
    // 设置连接池参数
    config.MaxConns = 20
    config.MinConns = 5
    config.MaxConnLifetime = 30 * time.Minute
    config.MaxConnIdleTime = 5 * time.Minute
    config.HealthCheckPeriod = 1 * time.Minute
    
    // 创建连接池
    pool, err := pgxpool.ConnectConfig(ctx, config)
    if err != nil {
        return nil, fmt.Errorf("连接数据库失败: %w", err)
    }
    
    return pool, nil
}
```

### 2. 缓存策略

- 实现多级缓存架构（本地缓存 + Redis）
- 使用布隆过滤器减少缓存穿透
- 实现缓存预热机制
- 采用缓存异步更新策略

```go
// 多级缓存实现示例
type MultiLevelCache struct {
    localCache *ristretto.Cache
    redisCache *redis.Client
    mutex      sync.RWMutex
}

func (c *MultiLevelCache) Get(ctx context.Context, key string) (interface{}, error) {
    // 1. 查询本地缓存
    if val, found := c.localCache.Get(key); found {
        metrics.CacheHits.WithLabelValues("local").Inc()
        return val, nil
    }
    
    // 2. 查询Redis缓存
    val, err := c.redisCache.Get(ctx, key).Result()
    if err == nil {
        // 回填本地缓存
        var decoded interface{}
        if err := json.Unmarshal([]byte(val), &decoded); err == nil {
            c.localCache.Set(key, decoded, 1)
        }
        metrics.CacheHits.WithLabelValues("redis").Inc()
        return decoded, nil
    } else if err != redis.Nil {
        return nil, err
    }
    
    metrics.CacheMisses.Inc()
    return nil, fmt.Errorf("缓存未命中")
}
```

### 3. 并发控制

- 使用工作池限制并发请求数
- 实现自适应限流算法
- 针对热点数据实现分片锁
- 使用上下文控制请求超时

```go
// 工作池实现示例
type WorkerPool struct {
    tasks       chan func()
    concurrency int
    wg          sync.WaitGroup
}

func NewWorkerPool(concurrency int) *WorkerPool {
    pool := &WorkerPool{
        tasks:       make(chan func(), concurrency*10),
        concurrency: concurrency,
    }
    
    pool.wg.Add(concurrency)
    for i := 0; i < concurrency; i++ {
        go func() {
            defer pool.wg.Done()
            for task := range pool.tasks {
                task()
            }
        }()
    }
    
    return pool
}

func (p *WorkerPool) Submit(task func()) error {
    select {
    case p.tasks <- task:
        return nil
    default:
        return fmt.Errorf("工作池已满")
    }
}
```

## 高可用部署架构

### 1. 多区域部署

- 实现跨区域数据同步
- 配置全局负载均衡
- 实现就近接入策略
- 区域故障自动切换

### 2. 灾备策略

- 定期数据备份与恢复演练
- 实现数据库主从复制
- 配置跨区域数据备份
- 制定完整的灾难恢复计划

### 3. 混沌工程实践

- 定期执行故障注入测试
- 模拟网络延迟和分区
- 测试服务过载场景
- 验证自动恢复机制

```go
// 混沌测试中间件示例
func ChaosTesting() grpc.UnaryServerInterceptor {
    return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
        // 根据配置决定是否注入故障
        if shouldInjectFailure() {
            failureType := selectFailureType()
            switch failureType {
            case "latency":
                // 注入延迟
                latency := randomLatency(100, 500) // 100-500ms
                time.Sleep(latency)
            case "error":
                // 注入错误
                return nil, status.Error(selectRandomErrorCode(), "混沌测试注入的错误")
            case "panic":
                // 注入panic (会被恢复)
                if rand.Float64() < 0.01 { // 1%的概率
                    panic("混沌测试注入的panic")
                }
            }
        }
        
        return handler(ctx, req)
    }
}
```

## API 版本管理

### 1. 版本兼容性策略

- 向后兼容性保证
- 字段弃用流程
- API 版本生命周期
- 兼容性测试套件

### 2. 版本迁移指南

- 平滑升级路径
- 客户端适配策略
- 版本共存期
- 迁移工具和脚本

### 3. API 变更管理

- 变更影响评估
- 变更通知机制
- 变更文档维护
- 自动化兼容性检查

## 联系我们

- 项目仓库: https://github.com/om-platform/api
- 问题反馈: https://github.com/om-platform/api/issues
- 邮件联系: api-team@om-platform.com
- 技术文档: https://docs.om-platform.com/api/user