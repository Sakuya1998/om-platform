# 智能运维平台 (Intelligent OM Platform)

基于 Kratos 框架构建的企业级智能运维管理平台，提供用户管理、组织架构、权限控制等核心功能。

## 项目特性

- 🚀 基于 Kratos v2 微服务框架
- 🔐 完整的用户认证与授权体系
- 🏢 灵活的组织架构管理
- 🎯 细粒度的权限控制
- 📡 支持 HTTP/gRPC 双协议
- 🔄 统一的 API 设计规范
- 📊 内置缓存与限流机制

## 快速开始

### 环境要求

- Go 1.19+
- Protocol Buffers 3.0+
- Kratos CLI v2

### 安装 Kratos CLI
```bash
go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
```

### 项目构建
```bash
# 下载依赖
make init

# 生成 API 代码
make api

# 构建项目
make build

# 运行服务
./bin/server -conf ./configs
```
## API 服务

### 用户服务 (User Service)

提供完整的用户管理功能，包括：

- **用户管理**: 用户CRUD操作、批量操作、状态管理
- **认证服务**: 登录/登出、令牌管理、密码重置
- **账户管理**: 个人资料、安全设置、第三方账户绑定
- **组织管理**: 组织架构、成员管理、层级关系
- **部门管理**: 部门CRUD、成员分配、树形结构
- **角色管理**: 角色定义、权限分配、用户角色绑定
- **权限管理**: 权限定义、权限检查、批量验证

### API 端点

所有API都支持HTTP和gRPC两种协议：

- **HTTP**: `http://localhost:8000/v1/`
- **gRPC**: `localhost:9000`

主要端点包括：
- `/v1/users` - 用户管理
- `/v1/auth` - 认证服务
- `/v1/organizations` - 组织管理
- `/v1/departments` - 部门管理
- `/v1/roles` - 角色管理
- `/v1/permissions` - 权限管理

## 项目结构

```
om-platform/
├── api/                    # API定义
│   └── user/service/v1/   # 用户服务API
├── app/                   # 应用代码
│   └── user/service/      # 用户服务实现
├── deploy/                # 部署配置
├── doc/                   # 项目文档
├── pkg/                   # 公共包
└── third_party/           # 第三方proto文件
```

## 开发工具

### Makefile 命令
```bash
# 下载和更新依赖
make init

# 生成API文件 (pb.go, http, grpc, validate, swagger)
make api

# 生成所有文件
make all
```

### Wire 依赖注入
```bash
# 安装 wire
go get github.com/google/wire/cmd/wire

# 生成依赖注入代码
cd cmd/server
wire
```

## 部署

### Docker 部署
```bash
# 构建镜像
docker build -t om-platform .

# 运行容器
docker run --rm -p 8000:8000 -p 9000:9000 \
  -v /path/to/configs:/data/conf \
  om-platform
```

### 配置说明

服务默认监听端口：
- HTTP: `8000`
- gRPC: `9000`

配置文件位置：`./configs/config.yaml`

## 技术特性

### 缓存策略
- 支持多级缓存配置
- 智能缓存键模式
- 可配置TTL时间

### 限流保护
- 基于令牌桶算法
- 支持突发流量处理
- 细粒度限流控制

### API 设计
- RESTful 风格设计
- 统一错误码体系
- 完整的参数验证
- 支持批量操作

## 文档

详细的项目规划和API规范请参考 `doc` 目录下的文档：

- [智能运维平台开发计划](./doc/intelligent_om_platform_development_plan.md)
- [用户API优化方案](./doc/user_api_optimization_proposal.md)
- [用户API优化示例](./doc/user_api_optimization_examples.md)

## 贡献

欢迎提交 Issue 和 Pull Request 来帮助改进项目。

## 许可证

本项目采用 MIT 许可证，详情请参阅 [LICENSE](./LICENSE) 文件。

