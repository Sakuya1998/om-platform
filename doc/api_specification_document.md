# 智能运维平台 API 接口规范文档

## 1. 引言

本文档定义了智能运维平台各微服务API接口的设计规范、指导原则和最佳实践。旨在确保API的一致性、可维护性和易用性，为平台各模块的开发与集成提供统一标准。

**Go 版本要求**: >=1.18

## 2. 通用API设计原则

- **RESTful风格 (HTTP网关)**: 对外暴露的HTTP API应遵循RESTful设计原则，使用标准的HTTP方法 (GET, POST, PUT, DELETE, PATCH)。
- **gRPC优先 (服务间通信)**: 微服务之间的内部通信优先采用gRPC，以获得高性能和强类型约束。
- **无状态服务**: 服务应设计为无状态，便于水平扩展和负载均衡。
- **幂等性**: 对于所有创建、更新或删除资源的非查询类操作，应保证幂等性。
- **命名规范**: 采用清晰、一致的命名约定。URL路径、gRPC服务名、方法名、消息名和字段名应具有描述性。
  - URL路径: `kebab-case` (例如: `/v1/alarm-rules`)
  - gRPC服务名: `PascalCase` (例如: `AlarmService`)
  - gRPC方法名: `PascalCase` (例如: `CreateAlarm`)
  - Protobuf消息名: `PascalCase` (例如: `AlarmRule`)
  - Protobuf字段名: `snake_case` (例如: `rule_id`, `display_name`)
- **错误处理**: API应返回明确的错误码和错误信息。HTTP API使用标准HTTP状态码，gRPC使用标准gRPC状态码。
- **资源释放**: 服务实现中必须确保数据库连接、文件句柄等资源得到及时、正确的释放。

## 3. API 版本控制

- **策略**: 
  - HTTP API: 通过URL路径进行版本控制，例如 `/api/v1/users`, `/api/v2/users`。
  - gRPC API: 通过Protobuf包名进行版本控制，例如 `package api.monitoring.v1;`。
- **向后兼容**: 尽量保持API的向后兼容性。对于破坏性变更，必须升级API版本。
- **废弃策略**: 当旧版本API不再维护时，应提前通知并设定合理的废弃过渡期。

## 4. 认证与授权

- **认证机制**: 平台API应采用统一的认证机制，例如基于OAuth 2.0的Bearer Token或JWT。
- **授权模型**: 根据业务需求，实现基于角色的访问控制 (RBAC) 或基于属性的访问控制 (ABAC)。
- **敏感数据**: API请求和响应中不得包含明文密码等敏感信息。

## 5. 请求与响应格式

- **HTTP API**: 
  - 数据格式: `application/json`。
  - 字符编码: `UTF-8`。
- **gRPC API**: 
  - 数据格式: Protobuf二进制。
- **标准错误响应 (HTTP JSON示例)**:
  ```json
  {
    "code": 40001, // 业务错误码
    "message": "Invalid input parameter: 'name' is required.",
    "details": [] // 可选，更详细的错误信息
  }
  ```
- **分页**: 对于返回集合资源的API，应支持分页。
  - 请求参数: `page_size` (每页数量), `page_token` (下一页令牌) 或 `offset` (偏移量), `limit` (数量)。
  - 响应体: 包含数据列表和下一页的 `next_page_token` 或总数 `total_count`。
- **数据校验**: API应对输入参数进行严格校验，对无效输入返回明确的错误信息。

## 6. Protobuf 最佳实践

- **标准类型**: 使用 `google/protobuf/timestamp.proto` 处理时间戳，`google/protobuf/duration.proto` 处理时间段，`google/protobuf/empty.proto` 表示空请求/响应等。
- **HTTP映射**: 使用 `google/api/annotations.proto` 为gRPC服务方法定义HTTP/JSON转换规则。
- **包管理**: 
  - `package` 声明应清晰反映服务和版本，例如 `api.servicename.v1`。
  - `option go_package` 应指向项目内正确的Go包路径，例如 `om-platform/api/servicename/v1;v1`。
- **枚举**: 枚举值的第一个值应为 `XXX_UNSPECIFIED = 0;`，作为默认值。
- **字段命名**: 采用 `snake_case`。
- **注释**: 为服务、方法、消息和字段添加清晰的注释。

## 7. 服务API定义 (示例)

以下为平台各核心微服务的API定义示例，具体字段和方法需根据详细设计进一步完善。

### 7.1. 监控中枢 (Monitoring Service)

- **职责**: 负责告警管理、指标数据处理、日志接入等。
- **示例 Proto (`api/monitoring/v1/alarm.proto`)**: (参考开发计划文档中的示例)
  ```protobuf
  syntax = "proto3";

  package api.monitoring.v1;

  import "google/api/annotations.proto";
  import "google/protobuf/timestamp.proto";
  import "google/protobuf/empty.proto";

  option go_package = "om-platform/api/monitoring/v1;v1";
  option java_multiple_files = true;
  option java_package = "api.monitoring.v1";

  // 告警服务
  service AlarmService {
    // 创建告警规则
    rpc CreateAlarmRule (CreateAlarmRuleRequest) returns (AlarmRule) {
      option (google.api.http) = {
        post: "/v1/monitoring/alarm-rules",
        body: "*"
      };
    }
    // 获取告警规则
    rpc GetAlarmRule (GetAlarmRuleRequest) returns (AlarmRule) {
      option (google.api.http) = {
        get: "/v1/monitoring/alarm-rules/{id}"
      };
    }
    // 更新告警规则
    rpc UpdateAlarmRule (UpdateAlarmRuleRequest) returns (AlarmRule) {
      option (google.api.http) = {
        put: "/v1/monitoring/alarm-rules/{alarm_rule.id}",
        body: "alarm_rule"
      };
    }
    // 删除告警规则
    rpc DeleteAlarmRule (DeleteAlarmRuleRequest) returns (google.protobuf.Empty) {
      option (google.api.http) = {
        delete: "/v1/monitoring/alarm-rules/{id}"
      };
    }
    // 列出告警规则
    rpc ListAlarmRules (ListAlarmRulesRequest) returns (ListAlarmRulesReply) {
      option (google.api.http) = {
        get: "/v1/monitoring/alarm-rules"
      };
    }

    // 上报告警事件
    rpc ReportAlarmEvent (ReportAlarmEventRequest) returns (AlarmEvent) {
        option (google.api.http) = {
            post: "/v1/monitoring/alarm-events",
            body: "*"
        };
    }
    // 查询告警事件
    rpc QueryAlarmEvents (QueryAlarmEventsRequest) returns (QueryAlarmEventsReply) {
        option (google.api.http) = {
            get: "/v1/monitoring/alarm-events"
        };
    }
  }

  message AlarmRule {
    string id = 1;          // 规则ID
    string name = 2;        // 规则名称
    string description = 3; // 描述
    string severity = 4;    // 告警级别 (e.g., CRITICAL, WARNING, INFO)
    string query = 5;       // 指标查询语句 (e.g., PromQL)
    string condition = 6;   // 触发条件 (e.g., "> 80")
    google.protobuf.Duration for_duration = 7; // 持续时间
    bool enabled = 8;       // 是否启用
    google.protobuf.Timestamp created_at = 9;
    google.protobuf.Timestamp updated_at = 10;
  }

  message CreateAlarmRuleRequest {
    string name = 1;
    string description = 2;
    string severity = 3;
    string query = 4;
    string condition = 5;
    google.protobuf.Duration for_duration = 6;
    bool enabled = 7;
  }

  message GetAlarmRuleRequest {
    string id = 1;
  }

  message UpdateAlarmRuleRequest {
    AlarmRule alarm_rule = 1;
  }

  message DeleteAlarmRuleRequest {
    string id = 1;
  }

  message ListAlarmRulesRequest {
    int32 page_size = 1;
    string page_token = 2;
    // 可以添加过滤条件，例如按名称、级别等
  }

  message ListAlarmRulesReply {
    repeated AlarmRule alarm_rules = 1;
    string next_page_token = 2;
    int32 total_size = 3;
  }

  message AlarmEvent {
    string id = 1;
    string rule_id = 2; //关联的告警规则ID
    string name = 3; // 告警名称 (可来自规则)
    string severity = 4;
    string message = 5;
    map<string, string> labels = 6; // 告警标签，如instance, job等
    google.protobuf.Timestamp triggered_at = 7;
    google.protobuf.Timestamp resolved_at = 8;
    string status = 9; // (e.g., FIRING, RESOLVED)
  }

  message ReportAlarmEventRequest {
    string rule_id = 1;
    string name = 2;
    string severity = 3;
    string message = 4;
    map<string, string> labels = 5;
  }

  message QueryAlarmEventsRequest {
    int32 page_size = 1;
    string page_token = 2;
    string rule_id_filter = 3; // 按规则ID过滤
    string severity_filter = 4; // 按级别过滤
    string status_filter = 5; // 按状态过滤
    google.protobuf.Timestamp start_time = 6;
    google.protobuf.Timestamp end_time = 7;
  }

  message QueryAlarmEventsReply {
    repeated AlarmEvent alarm_events = 1;
    string next_page_token = 2;
    int32 total_size = 3;
  }
  ```

### 7.2. 自动化引擎 (Automation Service)

- **职责**: 负责作业编排、自动化任务执行、自愈剧本管理。
- **示例 Proto (`api/automation/v1/job.proto`)**:
  ```protobuf
  syntax = "proto3";

  package api.automation.v1;

  import "google/api/annotations.proto";
  import "google/protobuf/timestamp.proto";
  import "google/protobuf/empty.proto";
  import "google/protobuf/struct.proto"; // 用于表示任意JSON结构的任务参数

  option go_package = "om-platform/api/automation/v1;v1";

  service JobService {
    // 创建剧本
    rpc CreatePlaybook (CreatePlaybookRequest) returns (Playbook) {
      option (google.api.http) = {
        post: "/v1/automation/playbooks",
        body: "*"
      };
    }
    // 获取剧本
    rpc GetPlaybook (GetPlaybookRequest) returns (Playbook) {
      option (google.api.http) = {
        get: "/v1/automation/playbooks/{id}"
      };
    }
    // ... 其他 Playbook CRUD ...

    // 执行作业 (基于剧本或直接任务)
    rpc ExecuteJob (ExecuteJobRequest) returns (JobExecution) {
      option (google.api.http) = {
        post: "/v1/automation/jobs/execute",
        body: "*"
      };
    }
    // 获取作业执行状态
    rpc GetJobExecution (GetJobExecutionRequest) returns (JobExecution) {
      option (google.api.http) = {
        get: "/v1/automation/jobs/executions/{id}"
      };
    }
    // ... 其他 JobExecution 操作 (如取消, 列出历史) ...
  }

  message Playbook {
    string id = 1;
    string name = 2;
    string description = 3;
    repeated PlaybookStep steps = 4; // 剧本步骤
    google.protobuf.Timestamp created_at = 5;
    google.protobuf.Timestamp updated_at = 6;
  }

  message PlaybookStep {
    string id = 1;
    string name = 2;
    string type = 3; // 步骤类型 (e.g., SCRIPT, HTTP_REQUEST, MANUAL_APPROVAL)
    google.protobuf.Struct parameters = 4; // 步骤参数
    // ... 其他步骤相关字段，如超时、重试策略 ...
  }

  message CreatePlaybookRequest {
    string name = 1;
    string description = 2;
    repeated PlaybookStep steps = 3;
  }

  message GetPlaybookRequest {
    string id = 1;
  }

  message JobExecution {
    string id = 1;
    string playbook_id = 2; // 可选，如果基于剧本
    string job_name = 3;
    string status = 4; // (e.g., PENDING, RUNNING, SUCCEEDED, FAILED, CANCELED)
    google.protobuf.Struct input_parameters = 5;
    google.protobuf.Struct output_results = 6;
    string logs = 7; // 执行日志或日志存储链接
    google.protobuf.Timestamp started_at = 8;
    google.protobuf.Timestamp finished_at = 9;
  }

  message ExecuteJobRequest {
    string playbook_id = 1; // 可选
    string job_name = 2;    // 如果不是基于剧本，则需要提供任务名称或类型
    google.protobuf.Struct input_parameters = 3;
    // ... 其他执行选项，如目标主机、凭证等 ...
  }

  message GetJobExecutionRequest {
    string id = 1;
  }
  ```

### 7.3. 配置管理数据库 (CMDB Service)

- **职责**: 负责资产管理、配置项 (CI) 管理、变更追踪、拓扑关系维护。
- **示例 Proto (`api/cmdb/v1/ci.proto`)**:
  ```protobuf
  syntax = "proto3";

  package api.cmdb.v1;

  import "google/api/annotations.proto";
  import "google/protobuf/timestamp.proto";
  import "google/protobuf/empty.proto";
  import "google/protobuf/struct.proto";

  option go_package = "om-platform/api/cmdb/v1;v1";

  service CiService {
    // 创建配置项类型 (模型)
    rpc CreateCiType (CreateCiTypeRequest) returns (CiType) {
      option (google.api.http) = {
        post: "/v1/cmdb/ci-types",
        body: "*"
      };
    }
    // 获取配置项类型
    rpc GetCiType (GetCiTypeRequest) returns (CiType) {
      option (google.api.http) = {
        get: "/v1/cmdb/ci-types/{id}"
      };
    }
    // ... 其他 CiType CRUD ...

    // 创建配置项实例
    rpc CreateCi (CreateCiRequest) returns (Ci) {
      option (google.api.http) = {
        post: "/v1/cmdb/cis",
        body: "*"
      };
    }
    // 获取配置项实例
    rpc GetCi (GetCiRequest) returns (Ci) {
      option (google.api.http) = {
        get: "/v1/cmdb/cis/{id}"
      };
    }
    // 更新配置项实例
    rpc UpdateCi (UpdateCiRequest) returns (Ci) {
      option (google.api.http) = {
        put: "/v1/cmdb/cis/{ci.id}",
        body: "ci"
      };
    }
    // ... 其他 Ci CRUD, ListCi, SearchCi ...

    // 管理CI关系
    rpc CreateCiRelation (CreateCiRelationRequest) returns (CiRelation) {
        option (google.api.http) = {
            post: "/v1/cmdb/ci-relations",
            body: "*"
        };
    }
    // ... 其他 CiRelation 操作 ...
  }

  message CiType {
    string id = 1;
    string name = 2; // e.g., Server, Database, Application
    string description = 3;
    repeated CiTypeAttribute attributes = 4; // 类型定义的属性
    google.protobuf.Timestamp created_at = 5;
    google.protobuf.Timestamp updated_at = 6;
  }

  message CiTypeAttribute {
    string name = 1;        // 属性名 (e.g., ip_address, os_version)
    string type = 2;        // 属性类型 (e.g., STRING, INTEGER, ENUM)
    bool required = 3;
    string description = 4;
    // ... 其他元数据，如校验规则、默认值 ...
  }

  message CreateCiTypeRequest {
    string name = 1;
    string description = 2;
    repeated CiTypeAttribute attributes = 3;
  }

  message GetCiTypeRequest {
    string id = 1;
  }

  message Ci {
    string id = 1;
    string ci_type_id = 2; // 关联的CI类型ID
    string name = 3;       // CI实例名称
    google.protobuf.Struct attributes = 4; // CI实例的属性值 (key-value)
    google.protobuf.Timestamp created_at = 5;
    google.protobuf.Timestamp updated_at = 6;
    // ... 其他元数据，如状态、负责人 ...
  }

  message CreateCiRequest {
    string ci_type_id = 1;
    string name = 2;
    google.protobuf.Struct attributes = 3;
  }

  message GetCiRequest {
    string id = 1;
  }

  message UpdateCiRequest {
    Ci ci = 1;
  }

  message CiRelation {
    string id = 1;
    string source_ci_id = 2;
    string target_ci_id = 3;
    string relation_type = 4; // e.g., CONNECTS_TO, RUNS_ON, DEPENDS_ON
    google.protobuf.Timestamp created_at = 5;
  }

  message CreateCiRelationRequest {
    string source_ci_id = 1;
    string target_ci_id = 2;
    string relation_type = 3;
  }
  ```

### 7.4. IT服务管理 (ITSM Service)

- **职责**: 负责工单流转、知识库管理、SLA管理。
- **示例 Proto (`api/itsm/v1/ticket.proto`)**:
  ```protobuf
  syntax = "proto3";

  package api.itsm.v1;

  import "google/api/annotations.proto";
  import "google/protobuf/timestamp.proto";
  import "google/protobuf/empty.proto";

  option go_package = "om-platform/api/itsm/v1;v1";

  service TicketService {
    // 创建工单
    rpc CreateTicket (CreateTicketRequest) returns (Ticket) {
      option (google.api.http) = {
        post: "/v1/itsm/tickets",
        body: "*"
      };
    }
    // 获取工单详情
    rpc GetTicket (GetTicketRequest) returns (Ticket) {
      option (google.api.http) = {
        get: "/v1/itsm/tickets/{id}"
      };
    }
    // 更新工单状态或内容
    rpc UpdateTicket (UpdateTicketRequest) returns (Ticket) {
      option (google.api.http) = {
        put: "/v1/itsm/tickets/{ticket.id}",
        body: "ticket"
      };
    }
    // 列出工单
    rpc ListTickets (ListTicketsRequest) returns (ListTicketsReply) {
      option (google.api.http) = {
        get: "/v1/itsm/tickets"
      };
    }
    // ... 其他工单操作，如添加评论、分配处理人 ...
  }

  message Ticket {
    string id = 1;
    string title = 2;
    string description = 3;
    string status = 4;        // e.g., OPEN, IN_PROGRESS, RESOLVED, CLOSED
    string priority = 5;      // e.g., HIGH, MEDIUM, LOW
    string category = 6;      // 工单分类
    string reporter_id = 7;   // 报告人ID
    string assignee_id = 8;   // 处理人ID
    google.protobuf.Timestamp created_at = 9;
    google.protobuf.Timestamp updated_at = 10;
    google.protobuf.Timestamp resolved_at = 11;
    // ... 其他字段，如关联CI、附件等 ...
  }

  message CreateTicketRequest {
    string title = 1;
    string description = 2;
    string priority = 3;
    string category = 4;
    // reporter_id 通常从认证信息中获取
  }

  message GetTicketRequest {
    string id = 1;
  }

  message UpdateTicketRequest {
    Ticket ticket = 1;
    // 或者只包含可更新的字段
    // string status = 2;
    // string assignee_id = 3;
    // string comment = 4;
  }

  message ListTicketsRequest {
    int32 page_size = 1;
    string page_token = 2;
    string status_filter = 3;
    string priority_filter = 4;
    string assignee_id_filter = 5;
  }

  message ListTicketsReply {
    repeated Ticket tickets = 1;
    string next_page_token = 2;
    int32 total_size = 3;
  }
  ```

### 7.5. 安全管控 (Security Service)

- **职责**: 负责漏洞管理、合规审计、访问控制策略管理。
- **示例 Proto (`api/security/v1/policy.proto`)**:
  ```protobuf
  syntax = "proto3";

  package api.security.v1;

  import "google/api/annotations.proto";
  import "google/protobuf/timestamp.proto";
  import "google/protobuf/empty.proto";

  option go_package = "om-platform/api/security/v1;v1";

  service PolicyService {
    // 创建访问控制策略
    rpc CreatePolicy (CreatePolicyRequest) returns (Policy) {
      option (google.api.http) = {
        post: "/v1/security/policies",
        body: "*"
      };
    }
    // 获取策略
    rpc GetPolicy (GetPolicyRequest) returns (Policy) {
      option (google.api.http) = {
        get: "/v1/security/policies/{id}"
      };
    }
    // ... 其他 Policy CRUD ...

    // 检查权限 (PDP - Policy Decision Point)
    rpc CheckPermission (CheckPermissionRequest) returns (CheckPermissionReply) {
      option (google.api.http) = {
        post: "/v1/security/permissions/check",
        body: "*"
      };
    }
  }

  message Policy {
    string id = 1;
    string name = 2;
    string description = 3;
    repeated PolicyRule rules = 4; // 策略规则 (e.g., OPA Rego, or custom structure)
    google.protobuf.Timestamp created_at = 5;
    google.protobuf.Timestamp updated_at = 6;
  }

  message PolicyRule {
    string effect = 1; // ALLOW / DENY
    repeated string actions = 2; // e.g., "read", "write", "execute"
    repeated string resources = 3; // e.g., "/api/monitoring/alarms/*"
    string condition = 4; // 可选的条件表达式
  }

  message CreatePolicyRequest {
    string name = 1;
    string description = 2;
    repeated PolicyRule rules = 3;
  }

  message GetPolicyRequest {
    string id = 1;
  }

  message CheckPermissionRequest {
    string subject_id = 1; // 用户或服务主体ID
    string action = 2;     // 请求的操作
    string resource = 3;   // 请求的资源
    // map<string, string> context = 4; // 可选的上下文属性
  }

  message CheckPermissionReply {
    bool allowed = 1;
    string reason = 2; // 如果不允许，说明原因
  }
  ```

### 7.6. 灾备恢复 (BCDR Service)

- **职责**: 负责灾备预案管理、演练验证、自动化切换。
- **示例 Proto (`api/bcdr/v1/plan.proto`)**:
  ```protobuf
  syntax = "proto3";

  package api.bcdr.v1;

  import "google/api/annotations.proto";
  import "google/protobuf/timestamp.proto";
  import "google/protobuf/empty.proto";

  option go_package = "om-platform/api/bcdr/v1;v1";

  service PlanService {
    // 创建灾备预案
    rpc CreateRecoveryPlan (CreateRecoveryPlanRequest) returns (RecoveryPlan) {
      option (google.api.http) = {
        post: "/v1/bcdr/plans",
        body: "*"
      };
    }
    // 获取灾备预案
    rpc GetRecoveryPlan (GetRecoveryPlanRequest) returns (RecoveryPlan) {
      option (google.api.http) = {
        get: "/v1/bcdr/plans/{id}"
      };
    }
    // ... 其他 RecoveryPlan CRUD ...

    // 执行灾备演练/切换
    rpc ExecuteRecoveryAction (ExecuteRecoveryActionRequest) returns (RecoveryActionExecution) {
      option (google.api.http) = {
        post: "/v1/bcdr/actions/execute",
        body: "*"
      };
    }
    // 获取执行状态
    rpc GetRecoveryActionExecution (GetRecoveryActionExecutionRequest) returns (RecoveryActionExecution) {
      option (google.api.http) = {
        get: "/v1/bcdr/actions/executions/{id}"
      };
    }
  }

  message RecoveryPlan {
    string id = 1;
    string name = 2;
    string description = 3;
    repeated RecoveryStep steps = 4; // 恢复步骤
    string rpo_target = 5; // 恢复点目标 (e.g., "15 minutes")
    string rto_target = 6; // 恢复时间目标 (e.g., "1 hour")
    google.protobuf.Timestamp created_at = 7;
    google.protobuf.Timestamp updated_at = 8;
  }

  message RecoveryStep {
    string id = 1;
    string name = 2;
    string description = 3;
    string action_type = 4; // e.g., FAILOVER_DB, RESTORE_BACKUP, UPDATE_DNS
    // ... 其他步骤参数 ...
  }

  message CreateRecoveryPlanRequest {
    string name = 1;
    string description = 2;
    repeated RecoveryStep steps = 3;
    string rpo_target = 4;
    string rto_target = 5;
  }

  message GetRecoveryPlanRequest {
    string id = 1;
  }

  message RecoveryActionExecution {
    string id = 1;
    string plan_id = 2;
    string action_type = 3; // DRILL (演练) / FAILOVER (切换) / FAILBACK (切回)
    string status = 4; // PENDING, RUNNING, SUCCEEDED, FAILED
    google.protobuf.Timestamp started_at = 5;
    google.protobuf.Timestamp finished_at = 6;
  }

  message ExecuteRecoveryActionRequest {
    string plan_id = 1;
    string action_type = 2;
  }

  message GetRecoveryActionExecutionRequest {
    string id = 1;
  }
  ```

### 7.4. 自动化引擎 (Automation Engine Service)

```protobuf
service AutomationEngineService {
  option (google.api.default_host) = "automation.example.com";

  // 运行一个自动化作业
  rpc RunJob (RunJobRequest) returns (JobRun) {
    option (google.api.http) = {
      post: "/v1/automation/jobs/run",
      body: "*"
    };
  }

  // 获取作业运行状态
  rpc GetJobRunStatus (GetJobRunStatusRequest) returns (JobRun) {
    option (google.api.http) = {
      get: "/v1/automation/job-runs/{job_run_id}"
    };
  }

  // 列出作业定义
  rpc ListJobs (ListJobsRequest) returns (ListJobsResponse) {
    option (google.api.http) = {
      get: "/v1/automation/jobs"
    };
  }
  
  // 创建作业定义
  rpc CreateJob (CreateJobRequest) returns (Job) {
    option (google.api.http) = {
      post: "/v1/automation/jobs",
      body: "job"
    };
  }

  // 获取作业定义
  rpc GetJob (GetJobRequest) returns (Job) {
    option (google.api.http) = {
      get: "/v1/automation/jobs/{job_id}"
    };
  }

  // 更新作业定义
  rpc UpdateJob (UpdateJobRequest) returns (Job) {
    option (google.api.http) = {
      put: "/v1/automation/jobs/{job.id}",
      body: "job"
    };
  }

  // 删除作业定义
  rpc DeleteJob (DeleteJobRequest) returns (google.protobuf.Empty) {
    option (google.api.http) = {
      delete: "/v1/automation/jobs/{job_id}"
    };
  }

  // 创建工作流定义
  rpc CreateWorkflow (CreateWorkflowRequest) returns (Workflow) {
    option (google.api.http) = {
      post: "/v1/automation/workflows",
      body: "workflow"
    };
  }

  // 获取工作流定义
  rpc GetWorkflow (GetWorkflowRequest) returns (Workflow) {
    option (google.api.http) = {
      get: "/v1/automation/workflows/{workflow_id}"
    };
  }

  // 更新工作流定义
  rpc UpdateWorkflow (UpdateWorkflowRequest) returns (Workflow) {
    option (google.api.http) = {
      put: "/v1/automation/workflows/{workflow.id}",
      body: "workflow"
    };
  }

  // 删除工作流定义
  rpc DeleteWorkflow (DeleteWorkflowRequest) returns (google.protobuf.Empty) {
    option (google.api.http) = {
      delete: "/v1/automation/workflows/{workflow_id}"
    };
  }

  // 列出工作流定义
  rpc ListWorkflows (ListWorkflowsRequest) returns (ListWorkflowsResponse) {
    option (google.api.http) = {
      get: "/v1/automation/workflows"
    };
  }

  // 运行工作流
  rpc RunWorkflow (RunWorkflowRequest) returns (WorkflowRun) {
    option (google.api.http) = {
      post: "/v1/automation/workflows/{workflow_id}/run",
      body: "*"
    };
  }

  // 获取工作流运行状态
  rpc GetWorkflowRunStatus (GetWorkflowRunStatusRequest) returns (WorkflowRun) {
    option (google.api.http) = {
      get: "/v1/automation/workflow-runs/{workflow_run_id}"
    };
  }
}

message Job {
  string id = 1;
  string name = 2;
  string description = 3;
  string type = 4; // e.g., script, ansible_playbook, container
  map<string, string> parameters = 5; // 作业参数定义
  string content = 6; // 作业内容，如脚本
  google.protobuf.Timestamp created_at = 7;
  google.protobuf.Timestamp updated_at = 8;
}

message RunJobRequest {
  string job_id = 1;
  map<string, string> arguments = 2; // 运行时参数
  string triggered_by = 3; // 触发者
}

message JobRun {
  string id = 1;
  string job_id = 2;
  string status = 3; // e.g., PENDING, RUNNING, SUCCEEDED, FAILED
  string output = 4; // 作业输出日志
  google.protobuf.Timestamp started_at = 5;
  google.protobuf.Timestamp finished_at = 6;
  map<string, string> arguments = 7;
}

message GetJobRunStatusRequest {
  string job_run_id = 1;
}

message ListJobsRequest {
  int32 page_size = 1;
  string page_token = 2;
  string filter = 3; // e.g., type='script'
}

message ListJobsResponse {
  repeated Job jobs = 1;
  string next_page_token = 2;
  int32 total_size = 3;
}

message CreateJobRequest {
  Job job = 1;
}

message GetJobRequest {
  string job_id = 1;
}

message UpdateJobRequest {
  Job job = 1;
}

message DeleteJobRequest {
  string job_id = 1;
}

message Workflow {
  string id = 1;
  string name = 2;
  string description = 3;
  repeated WorkflowStep steps = 4;
  google.protobuf.Timestamp created_at = 5;
  google.protobuf.Timestamp updated_at = 6;
}

message WorkflowStep {
  string name = 1;
  string job_id = 2;
  map<string, string> arguments_mapping = 3; // 参数映射
  repeated string depends_on = 4; // 依赖步骤名称
}

message CreateWorkflowRequest {
  Workflow workflow = 1;
}

message GetWorkflowRequest {
  string workflow_id = 1;
}

message UpdateWorkflowRequest {
  Workflow workflow = 1;
}

message DeleteWorkflowRequest {
  string workflow_id = 1;
}

message ListWorkflowsRequest {
  int32 page_size = 1;
  string page_token = 2;
  string filter = 3;
}

message ListWorkflowsResponse {
  repeated Workflow workflows = 1;
  string next_page_token = 2;
  int32 total_size = 3;
}

message RunWorkflowRequest {
  string workflow_id = 1;
  map<string, string> initial_arguments = 2; // 工作流启动参数
  string triggered_by = 3;
}

message WorkflowRun {
  string id = 1;
  string workflow_id = 2;
  string status = 3; // e.g., PENDING, RUNNING, SUCCEEDED, FAILED
  repeated JobRun job_runs = 4;
  google.protobuf.Timestamp started_at = 5;
  google.protobuf.Timestamp finished_at = 6;
}

message GetWorkflowRunStatusRequest {
  string workflow_run_id = 1;
}
```

### 7.5. 日志服务 (Log Service)

```protobuf
service LogService {
  option (google.api.default_host) = "log.example.com";

  // 接收日志条目 (单个或批量)
  rpc IngestLog (IngestLogRequest) returns (IngestLogResponse) {
    option (google.api.http) = {
      post: "/v1/logs/ingest",
      body: "*"
    };
  }

  // 查询日志
  rpc QueryLogs (QueryLogsRequest) returns (QueryLogsResponse) {
    option (google.api.http) = {
      post: "/v1/logs/query", // 使用POST以支持复杂的查询体
      body: "*"
      // 或者 GET /v1/logs?query=...&start_time=...&end_time=...&limit=...
    };
  }

  // 获取日志统计信息 (如日志量、错误分布等)
  rpc GetLogStats (GetLogStatsRequest) returns (GetLogStatsResponse) {
    option (google.api.http) = {
      get: "/v1/logs/stats"
    };
  }
}

message LogEntry {
  string id = 1; // 可选，由日志服务生成或客户端提供
  google.protobuf.Timestamp timestamp = 2;
  string service_name = 3;
  string host_name = 4;
  string level = 5; // e.g., INFO, ERROR, WARN, DEBUG
  string message = 6;
  map<string, google.protobuf.Value> fields = 7; // 结构化日志字段
  string trace_id = 8; // 关联追踪ID
  string span_id = 9;  // 关联Span ID
}

message IngestLogRequest {
  repeated LogEntry entries = 1;
}

message IngestLogResponse {
  int32 ingested_count = 1;
  repeated string failed_ids = 2; // 记录 ingest 失败的日志ID (如果客户端提供了ID)
}

message QueryLogsRequest {
  string query_string = 1; // 查询语句 (e.g., Lucene syntax, PromQL-like for logs)
  google.protobuf.Timestamp start_time = 2;
  google.protobuf.Timestamp end_time = 3;
  int32 limit = 4;
  string page_token = 5; // 用于分页
  enum SortOrder {
    ASC = 0;
    DESC = 1;
  }
  SortOrder sort_order = 6;
}

message QueryLogsResponse {
  repeated LogEntry entries = 1;
  string next_page_token = 2;
  int32 total_count = 3; // 匹配的总条目数 (可能为估算值)
}

message GetLogStatsRequest {
  string query_string = 1; // 针对特定查询的统计
  google.protobuf.Timestamp start_time = 2;
  google.protobuf.Timestamp end_time = 3;
  repeated string group_by_fields = 4; // 按字段分组统计
}

message LogStatPoint {
  map<string, string> dimensions = 1; // 分组维度
  int64 count = 2;
}

message GetLogStatsResponse {
  repeated LogStatPoint stats = 1;
}
```

### 7.6. 认证授权服务 (Auth Service)

```protobuf
service AuthService {
  option (google.api.default_host) = "auth.example.com";

  // 用户登录
  rpc Login (LoginRequest) returns (LoginResponse) {
    option (google.api.http) = {
      post: "/v1/auth/login",
      body: "*"
    };
  }

  // 用户登出 (通常是客户端行为，服务端可选择实现使token失效)
  rpc Logout (LogoutRequest) returns (google.protobuf.Empty) {
    option (google.api.http) = {
      post: "/v1/auth/logout",
      body: "*"
    };
  }

  // 刷新访问令牌
  rpc RefreshToken (RefreshTokenRequest) returns (LoginResponse) {
    option (google.api.http) = {
      post: "/v1/auth/refresh-token",
      body: "*"
    };
  }

  // 验证令牌 (通常由API网关或服务自身调用)
  rpc ValidateToken (ValidateTokenRequest) returns (ValidateTokenResponse) {
    option (google.api.http) = {
      post: "/v1/auth/validate-token",
      body: "*"
    };
  }

  // 检查权限
  rpc CheckPermission (CheckPermissionRequest) returns (CheckPermissionResponse) {
    option (google.api.http) = {
      post: "/v1/auth/check-permission",
      body: "*"
    };
  }
}

message LoginRequest {
  string username = 1;
  string password = 2;
  // string grant_type = 3; // e.g., "password", "sso_ticket"
}

message LoginResponse {
  string access_token = 1;
  string token_type = 2; // e.g., "Bearer"
  int64 expires_in = 3;  // access_token 有效期 (秒)
  string refresh_token = 4;
  string user_id = 5;
}

message LogoutRequest {
  string access_token = 1; // 可选，用于服务端使特定token失效
}

message RefreshTokenRequest {
  string refresh_token = 1;
}

message ValidateTokenRequest {
  string access_token = 1;
}

message ValidateTokenResponse {
  bool valid = 1;
  string user_id = 2;
  repeated string scopes = 3; // 令牌具有的权限范围
  int64 expires_at = 4; // 令牌过期时间戳 (Unix seconds)
}

message CheckPermissionRequest {
  string user_id = 1; // 或从token中解析
  string permission_name = 2; // e.g., "cmdb:ci:create", "automation:job:run"
  string resource_id = 3; // 可选，针对特定资源的权限检查
}

message CheckPermissionResponse {
  bool granted = 1;
}
```

### 7.7. 用户管理服务 (User Management Service)

```protobuf
service UserManagementService {
  option (google.api.default_host) = "user.example.com";

  // 创建用户
  rpc CreateUser (CreateUserRequest) returns (User) {
    option (google.api.http) = {
      post: "/v1/users",
      body: "user"
    };
  }

  // 获取用户信息
  rpc GetUser (GetUserRequest) returns (User) {
    option (google.api.http) = {
      get: "/v1/users/{user_id}"
    };
  }

  // 更新用户信息
  rpc UpdateUser (UpdateUserRequest) returns (User) {
    option (google.api.http) = {
      put: "/v1/users/{user.id}",
      body: "user"
    };
  }

  // 删除用户
  rpc DeleteUser (DeleteUserRequest) returns (google.protobuf.Empty) {
    option (google.api.http) = {
      delete: "/v1/users/{user_id}"
    };
  }

  // 列出用户
  rpc ListUsers (ListUsersRequest) returns (ListUsersResponse) {
    option (google.api.http) = {
      get: "/v1/users"
    };
  }

  // 为用户分配角色
  rpc AssignRoleToUser (AssignRoleToUserRequest) returns (google.protobuf.Empty) {
    option (google.api.http) = {
      post: "/v1/users/{user_id}/roles",
      body: "*"
    };
  }

  // 从用户移除角色
  rpc RemoveRoleFromUser (RemoveRoleFromUserRequest) returns (google.protobuf.Empty) {
    option (google.api.http) = {
      delete: "/v1/users/{user_id}/roles/{role_id}"
    };
  }

  // 获取用户角色列表
  rpc ListUserRoles (ListUserRolesRequest) returns (ListUserRolesResponse) {
    option (google.api.http) = {
      get: "/v1/users/{user_id}/roles"
    };
  }

  // 创建角色
  rpc CreateRole (CreateRoleRequest) returns (Role) {
    option (google.api.http) = {
      post: "/v1/roles",
      body: "role"
    };
  }

  // 获取角色信息
  rpc GetRole (GetRoleRequest) returns (Role) {
    option (google.api.http) = {
      get: "/v1/roles/{role_id}"
    };
  }

  // 更新角色信息
  rpc UpdateRole (UpdateRoleRequest) returns (Role) {
    option (google.api.http) = {
      put: "/v1/roles/{role.id}",
      body: "role"
    };
  }

  // 删除角色
  rpc DeleteRole (DeleteRoleRequest) returns (google.protobuf.Empty) {
    option (google.api.http) = {
      delete: "/v1/roles/{role_id}"
    };
  }

  // 列出角色
  rpc ListRoles (ListRolesRequest) returns (ListRolesResponse) {
    option (google.api.http) = {
      get: "/v1/roles"
    };
  }
  
  // 为角色分配权限
  rpc AssignPermissionToRole (AssignPermissionToRoleRequest) returns (google.protobuf.Empty) {
    option (google.api.http) = {
      post: "/v1/roles/{role_id}/permissions",
      body: "*"
    };
  }

  // 从角色移除权限
  rpc RemovePermissionFromRole (RemovePermissionFromRoleRequest) returns (google.protobuf.Empty) {
    option (google.api.http) = {
      delete: "/v1/roles/{role_id}/permissions/{permission_name}"
    };
  }

  // 列出角色权限
  rpc ListRolePermissions (ListRolePermissionsRequest) returns (ListRolePermissionsResponse) {
    option (google.api.http) = {
      get: "/v1/roles/{role_id}/permissions"
    };
  }

  // 创建权限定义 (通常由系统预设，较少动态创建)
  rpc CreatePermission (CreatePermissionRequest) returns (Permission) {
      option (google.api.http) = {
          post: "/v1/permissions",
          body: "permission"
      };
  }

  // 列出所有权限定义
  rpc ListPermissions (ListPermissionsRequest) returns (ListPermissionsResponse) {
      option (google.api.http) = {
          get: "/v1/permissions"
      };
  }
}

message User {
  string id = 1;
  string username = 2;
  // password_hash string = 3; // 不应在API中返回
  string email = 4;
  string display_name = 5;
  bool is_active = 6;
  google.protobuf.Timestamp created_at = 7;
  google.protobuf.Timestamp updated_at = 8;
  map<string, string> metadata = 9; // 额外信息，如部门、电话等
}

message CreateUserRequest {
  User user = 1;
  string password = 2; // 创建时需要明文密码
}

message GetUserRequest {
  string user_id = 1;
}

message UpdateUserRequest {
  User user = 1;
  // string password = 2; // 可选，用于修改密码
}

message DeleteUserRequest {
  string user_id = 1;
}

message ListUsersRequest {
  int32 page_size = 1;
  string page_token = 2;
  string filter = 3; // e.g., "is_active=true AND email LIKE '%@example.com%'"
}

message ListUsersResponse {
  repeated User users = 1;
  string next_page_token = 2;
  int32 total_size = 3;
}

message Role {
  string id = 1;
  string name = 2; // e.g., "admin", "viewer", "editor"
  string description = 3;
  google.protobuf.Timestamp created_at = 4;
  google.protobuf.Timestamp updated_at = 5;
}

message AssignRoleToUserRequest {
  string user_id = 1;
  string role_id = 2;
}

message RemoveRoleFromUserRequest {
  string user_id = 1;
  string role_id = 2;
}

message ListUserRolesRequest {
  string user_id = 1;
}

message ListUserRolesResponse {
  repeated Role roles = 1;
}

message CreateRoleRequest {
  Role role = 1;
}

message GetRoleRequest {
  string role_id = 1;
}

message UpdateRoleRequest {
  Role role = 1;
}

message DeleteRoleRequest {
  string role_id = 1;
}

message ListRolesRequest {
  int32 page_size = 1;
  string page_token = 2;
}

message ListRolesResponse {
  repeated Role roles = 1;
  string next_page_token = 2;
  int32 total_size = 3;
}

message Permission {
    string name = 1; // e.g., "cmdb:ci:read", "automation:job:execute"
    string description = 2;
}

message AssignPermissionToRoleRequest {
    string role_id = 1;
    string permission_name = 2;
}

message RemovePermissionFromRoleRequest {
    string role_id = 1;
    string permission_name = 2;
}

message ListRolePermissionsRequest {
    string role_id = 1;
}

message ListRolePermissionsResponse {
    repeated Permission permissions = 1;
}

message CreatePermissionRequest {
    Permission permission = 1;
}

message ListPermissionsRequest {
    int32 page_size = 1;
    string page_token = 2;
}

message ListPermissionsResponse {
    repeated Permission permissions = 1;
    string next_page_token = 2;
    int32 total_size = 3;
}

```

## 8. 附录

- **HTTP状态码参考**: [MDN HTTP response status codes](https://developer.mozilla.org/en-US/docs/Web/HTTP/Status)
- **gRPC状态码参考**: [gRPC Status codes and their use in gRPC](https://grpc.github.io/grpc/core/md_doc_statuscodes.html)

---

本文档为API接口规范的初步定义，具体细节将随着项目进展和详细设计进行迭代和完善。