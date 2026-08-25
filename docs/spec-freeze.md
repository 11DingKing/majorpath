# Foundation 规格冻结表

| 范畴 | 冻结结论 |
|---|---|
| 业务边界 | 面向咨询师与审核员的学生画像、专业/院校招录映射、职业路径推荐与跟进，不涉及电商、投票、CRM、库存或其它禁止题材。 |
| 持久化 | SQLite 真实关系数据库，database/sql + go-sqlite3；启动执行版本化 migration，外键、唯一约束、业务索引和时间字段完整。 |
| 表关系 | users/sessions/students/interests/majors/universities/admission_plans/career_paths/recommendations/recommendation_items/followups/audit_events/idempotency_keys，跨实体外键和级联策略。 |
| 事务 | 创建推荐方案在单事务中锁定学生、写方案项、幂等键和审计；审核与跟进状态转换同事务提交，失败回滚。 |
| 状态机 | recommendation draft→submitted→approved→archived；followup pending→running→done/failed，非法转换拒绝。 |
| 并发 | recommendation.version 乐观锁；计划名额使用条件更新；worker 使用租约字段避免重复领取。 |
| context | HTTP request context 贯穿 handler/service/repository；worker 使用可取消 context 和超时。 |
| worker | 到期跟进扫描、有限重试、指数退避、永久失败审计、优雅停止和重启恢复。 |
| 错误传播 | domain sentinel errors，repository 保留 errors.Is 链，HTTP 映射稳定 code/status/request_id。 |
| HTTP | 登录/退出、学生、专业、推荐、审核、跟进、审计、健康/就绪接口；请求 ID、日志、panic recovery。 |
| 身份权限 | 可撤销会话、过期检查、counselor/reviewer 两角色；service 与 HTTP 双层校验。 |
| Docker | 多阶段真实构建 ./cmd/server，migrations 复制到 /app/migrations，默认入口提供 healthcheck。 |
| 测试 | 领域、service、SQLite migration/事务/重启、HTTP 契约、并发、worker、分页和错误链；测试总量≥1500物理行。 |
| 规模 | compact_10：生产 Go≥2000 行、≥20 文件、≥10 package；测试 Go≥1500 行；不以空壳或重复结构凑数。 |
| 禁止题材 | 已核对 excluded-topics.md，大学专业选择与招录映射未来发展不属于禁止清单。 |
| 后续容量 | 10 个独立运行时边界：事务回滚、状态转换、乐观并发、名额冲突、幂等生命周期、context 取消、worker 重试/恢复、审计一致性、会话撤销/过期、分页组合查询。此阶段不植入 Bug、不建题目分支。 |
