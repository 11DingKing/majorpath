# MajorPath

MajorPath 为高校专业选择与招录映射未来发展提供业务服务。咨询师创建学生画像，系统把画像与专业、招生批次、职业路径和院校计划关联，形成可追踪的推荐方案；学生确认方案后进入跟进流程，后台 worker 处理过期提醒并保留审计记录。

## 运行

```bash
GOTOOLCHAIN=local go test ./... -count=1
GOTOOLCHAIN=local go run ./cmd/server
```

默认监听 `:8080`，`/healthz` 为存活检查，`/readyz` 验证数据库与迁移状态。首次启动从 `migrations` 目录按版本创建 SQLite 数据库。

## 身份与业务角色

`POST /v1/auth/login` 返回可撤销会话 token，`POST /v1/auth/logout` 撤销当前会话。`counselor` 可以维护学生与推荐方案，`reviewer` 可以审核方案并查看审计记录。会话过期、撤销和角色权限在 service 与 HTTP 层均有验证。

## 主要流程

1. 登录后创建学生画像并登记兴趣标签。
2. 浏览专业、院校招生计划与职业路径，创建推荐方案。
3. 方案依次经过 `draft -> submitted -> approved -> archived`，审核需要乐观锁版本。
4. 确认方案生成跟进任务，worker 在截止时间到达后写入提醒与审计事件。

数据库包含 users、sessions、students、interests、majors、universities、admission_plans、career_paths、recommendations、recommendation_items、followups、audit_events、idempotency_keys 共 13 张关联表。
