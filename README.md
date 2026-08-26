# TeamFlow

**TeamFlow 企业级协作平台**

类 Jira / 飞书任务协作的 Go 后端服务，支持团队、项目、任务的一站式管理，内置实时通知与消息队列等企业级特性。

## ✨ 功能特性

- 🧑‍💼 **用户系统** — 注册 / 登录 / JWT 认证
- 🧩 **团队 / 项目 / 任务** — 完整的 CRUD 与状态流转
- 🔐 **RBAC 权限** — 角色与权限控制
- 💾 **Redis 缓存** — 热点数据缓存与加速
- 📁 **文件存储** — 本地文件上传
- 📝 **操作日志** — 关键操作审计追踪
- 🔔 **实时通知** — WebSocket 实时推送
- 📬 **消息队列** — MQ 异步解耦
- ⏱ **定时任务** — 周期性任务调度
- 🐳 **Docker 部署** — 一键容器化部署

## 🛠 技术栈

| 分类 | 技术 |
| ---- | ---- |
| 框架 | Gin |
| ORM | GORM |
| 数据库 | MySQL |
| 缓存 | Redis |
| 认证 | JWT |
| 实时通信 | WebSocket |
| 部署 | Docker |

## 📦 快速开始

```bash
# 启动依赖（MySQL / Redis / RabbitMQ）
docker-compose up -d

# 运行
go run cmd/main.go

# 构建
go build -o bin/teamflow cmd/main.go

# 测试
go test ./...
```

## 📁 项目结构

```
TeamFlow/
├── cmd/main.go              # 应用入口
├── config/config.yaml       # 配置文件
├── internal/
│   ├── controller/          # Gin handler
│   ├── service/             # 业务逻辑层
│   ├── repository/          # GORM 数据库操作层
│   ├── model/               # GORM 模型
│   ├── middleware/          # 中间件（JWT、日志、CORS 等）
│   ├── router/              # 路由注册
│   └── dto/                 # 请求/响应 DTO
├── pkg/
│   ├── database/            # MySQL 和 Redis 连接
│   ├── jwt/                 # Token 生成与解析
│   ├── response/            # 统一响应格式
│   └── utils/               # 工具函数
└── storage/uploads/         # 本地上传文件
```

## 🔗 数据流向

```
Router → Middleware → Controller → Service → Repository → DB/Redis
```

## 📚 文档

完整的开发教程文档位于 [`TeamFlow企业级开发教程/`](./TeamFlow企业级开发教程)，共 87 篇，从项目准备、数据库设计到 WebSocket、消息队列、Docker 部署，循序渐进。

## 📄 License

MIT
