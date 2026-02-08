---
name: pinche-fullstack
description: Pinche 春节拼车应用全栈开发指南。当开发该项目的后端 (Go/Gin)、前端 (Vue 3/Vite/Pinia/TailwindCSS)、数据库 (MySQL) 相关功能时，加载此技能获取项目架构、编码规范和开发流程指导。
---

# Pinche 全栈开发技能

## 概述

Pinche 是一个基于 **Vue 3 + Golang** 的全栈拼车应用，解决春节返乡出行问题。支持司机发布"车找人"和乘客发布"人找车"，通过智能匹配撮合双方。

## 触发条件

以下场景应加载此技能：
- 开发 Pinche 项目的后端接口
- 开发 Pinche 项目的前端页面
- 修改数据库表结构
- 添加 WebSocket 实时通知
- 处理主题系统或样式问题

## 项目结构

```
pinche/
├── server/                    # Go 后端服务
│   ├── cmd/main.go           # 入口
│   ├── config/               # 配置管理
│   └── internal/
│       ├── handler/          # HTTP 处理器
│       ├── service/          # 业务逻辑层
│       ├── repository/       # 数据访问层
│       ├── model/            # 数据模型
│       ├── middleware/       # 中间件
│       ├── router/           # 路由配置
│       ├── websocket/        # WebSocket Hub
│       └── database/         # 数据库连接
├── app/                       # Vue 用户端 (端口 3000)
├── admin/                     # Vue 管理端 (端口 3001)
└── sql/init.sql              # 数据库初始化
```

## 后端开发流程

### 分层架构

```
HTTP 请求 → Handler → Service → Repository → MySQL
              ↓          ↓          ↓
           参数校验    业务逻辑    SQL 查询
```

### 新增接口步骤

1. 在 `model/` 定义请求/响应结构体
2. 在 `repository/` 实现数据访问方法（必须使用参数化查询）
3. 在 `service/` 实现业务逻辑
4. 在 `handler/` 实现 HTTP 处理
5. 在 `router/router.go` 注册路由

详细代码模板参考 `references/backend_patterns.md`。

### API 响应格式

```go
// 统一响应
type Response struct {
    Code    int         `json:"code"`    // 0=成功
    Message string      `json:"message"`
    Data    interface{} `json:"data"`
}

c.JSON(http.StatusOK, model.Success(data))
c.JSON(http.StatusOK, model.Error(model.ErrCodeBadRequest, "错误信息"))
```

## 前端开发流程

### 新增页面步骤

1. 在 `views/` 创建 Vue 组件
2. 在 `router/index.js` 添加路由
3. 如需状态管理，在 `stores/` 创建 Pinia store

详细代码模板参考 `references/frontend_patterns.md`。

### 样式规范

- 使用 TailwindCSS 工具类
- 主题色使用 CSS 变量：`var(--theme-primary)`
- 移动端触摸目标最小 44px
- 支持三种主题：`light` | `dark` | `spring`

## 数据库设计

| 表名 | 说明 |
|------|------|
| users | 用户信息 |
| trips | 行程信息 |
| matches | 匹配记录 |
| messages | 私聊消息 |
| notifications | 通知 |
| announcements | 公告 |
| trip_grabs | 抢单记录 |

### 关键设计

- 使用 `open_id` 对外暴露用户标识，隐藏内部 `id`
- 行程状态：1=待匹配, 2=已匹配, 3=已完成, 4=已取消, 5=已封禁, 6=已过期
- 匹配状态：0=待确认, 1=成功, 2=失败

详细字段定义参考 `references/database_schema.md`。

## 开发命令

```bash
# 后端
cd server && go run cmd/main.go

# 用户端
cd app && npm run dev

# 管理端
cd admin && npm run dev

# 构建
cd app && npm run build
cd admin && npm run build
```

## 安全规范

- **必须**使用参数化 SQL 查询，禁止字符串拼接
- **必须**验证用户权限，特别是操作他人数据时
- 敏感数据使用 bcrypt 加密

## 性能规范

- 列表接口必须分页
- 大量数据查询添加索引
- 前端组件及时清理定时器和监听器

## 日志规范

- 使用 zap (sugared)+lumberjack 打印日志
- 日志级别：DEBUG > INFO > WARN > ERROR > FATAL
- 日志格式：时间戳 | 等级 | 文件名 | 行号 | 函数名 | 消息
- 日志文件按天分割，单日志文件大小限制200MB
- 请求和回复需要记录，敏感信息使用掩码
- 异常分支必须有日志记录，且错误日志必须包含堆栈
