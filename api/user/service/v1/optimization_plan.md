# 用户服务API定义优化方案

## 1. 现状分析

通过对`/api/user/service/v1/`目录下的API定义文件进行分析，发现以下问题：

### 1.1 结构问题

- **服务职责划分不清晰**：部分服务之间存在功能重叠，如`UserService`包含了账户状态管理功能，这些功能原本属于`AccountService`
- **消息类型重复定义**：多个proto文件中存在相似的消息类型定义，缺乏复用
- **缺少公共消息类型文件**：基础数据类型和共享消息类型没有集中管理
- **导入关系混乱**：部分文件导入了不必要的包，或者缺少必要的导入

### 1.2 命名问题

- **命名风格不一致**：字段命名混用驼峰和下划线风格
- **枚举值命名不规范**：部分枚举值使用全大写下划线分隔，部分使用驼峰
- **服务方法命名不统一**：相似功能的方法在不同服务中命名不一致

### 1.3 注释问题

- **注释不完整**：部分服务和方法缺少详细注释
- **注释格式不统一**：有的使用行内注释，有的使用块注释
- **字段描述重复**：同时使用注释和JSON注解描述同一字段

## 2. 优化方案

### 2.1 结构优化

#### 2.1.1 创建公共消息类型文件

创建`common.proto`文件，将共享的基础消息类型移至此文件：

```protobuf
syntax = "proto3";

package api.user.service.v1;

import "google/protobuf/timestamp.proto";

option go_package = "om-platform/api/user/service/v1;v1";
option java_multiple_files = true;
option java_package = "api.user.service.v1";

// 用户状态（基本状态）
enum UserAccountStatus {
  USER_ACCOUNT_OFF = 0; // 禁用
  USER_ACCOUNT_ON = 1;  // 启用
}

// 用户权限级别
enum UserAuthority {
  SYS_ADMIN = 0;     // 系统超级用户
  SYS_MANAGER = 1;   // 系统管理员
  CUSTOMER_USER = 2; // 普通用户
  GUEST_USER = 3;    // 游客
  REFRESH_TOKEN = 4; // 刷新令牌
}

// 用户性别
enum UserGender {
  SECRET = 0; // 未知
  MALE = 1;   // 男性
  FEMALE = 2; // 女性
}

// 基础审计信息
message AuditInfo {
  uint32 create_by = 1;                           // 创建者ID
  google.protobuf.Timestamp create_time = 2;      // 创建时间
  uint32 update_by = 3;                           // 更新者ID
  google.protobuf.Timestamp update_time = 4;      // 更新时间
}

// 基础租户信息
message TenantInfo {
  uint32 tenant_id = 1;                           // 租户ID
  string tenant_name = 2;                         // 租户名称
}
```

#### 2.1.2 重组服务职责

- **AuthService**：专注于认证、令牌和会话管理
- **UserService**：专注于用户基本信息管理
- **AccountService**：专注于账户状态管理（从UserService分离）
- **UserPreferenceService**：保持不变，专注于用户偏好设置
- **IdentityService**：专注于身份提供商集成
- **PermissionService**：专注于权限管理
- **RoleService**：专注于角色管理
- **OperationLogService**：专注于操作日志管理

### 2.2 命名规范优化

#### 2.2.1 统一字段命名风格

- 所有字段名使用下划线命名法（snake_case）
- JSON标签使用驼峰命名法（camelCase）
- 保持proto字段名与Go字段名的一致性

#### 2.2.2 统一枚举值命名风格

- 所有枚举值使用全大写下划线分隔（SCREAMING_SNAKE_CASE）
- 枚举类型名使用驼峰命名法（PascalCase）

#### 2.2.3 统一服务方法命名

- 列表查询：`List<Resource>`
- 详情查询：`Get<Resource>`
- 创建操作：`Create<Resource>`
- 更新操作：`Update<Resource>`
- 删除操作：`Delete<Resource>`
- 批量操作：`Batch<Operation><Resource>`

### 2.3 注释完善

#### 2.3.1 统一注释格式

- 服务定义：使用多行注释，描述服务职责和功能范围
- 方法定义：使用单行注释，描述方法功能、参数和返回值
- 消息类型：使用多行注释，描述消息用途和重要字段
- 字段定义：使用行尾注释，简要描述字段含义

#### 2.3.2 避免重复描述

- 移除重复的字段描述，优先使用JSON注解
- 注释内容应补充而非重复字段名含义

## 3. 实施路径

### 3.1 第一阶段：结构优化

1. 创建`common.proto`文件，提取共享消息类型
2. 调整服务职责划分，创建`account.proto`文件
3. 更新导入关系，确保依赖正确

### 3.2 第二阶段：命名规范化

1. 统一字段命名风格
2. 统一枚举值命名风格
3. 统一服务方法命名

### 3.3 第三阶段：注释完善

1. 完善服务和方法注释
2. 统一注释格式
3. 移除重复描述

## 4. 优化效果

- **代码质量提升**：结构清晰，命名规范，注释完善
- **开发效率提高**：减少重复定义，提高代码复用
- **维护成本降低**：服务职责明确，接口文档完整
- **兼容性保障**：保持向后兼容，不影响现有功能