# 春节拼车应用

一个基于 Vue.js + Golang 的简易拼车应用，专门解决春节返乡期间的出行难题。

## 🌐 在线体验

**网站地址**：[https://pinche.chat](https://pinche.chat)

## 功能特性

- **行程发布**：司机/乘客可发布包含路线、时间、座位数等信息的行程
- **智能匹配**：基于路线相似度、时间匹配等条件自动匹配司机与乘客
- **实时通知**：通过 WebSocket 实时推送匹配通知
- **双向确认**：司机和乘客双方确认后才正式成立拼车
- **联系交换**：拼车成功后自动交换双方联系方式

## 技术架构

- **前端**：Vue 3 + Vite + TailwindCSS + Pinia
- **后端**：Golang + Gin + MySQL
- **实时通信**：WebSocket

## 项目结构

```
pinche/
├── app/                    # 前端 Vue 应用
│   ├── src/
│   │   ├── views/          # 页面组件
│   │   ├── layouts/        # 布局组件
│   │   ├── stores/         # Pinia 状态管理
│   │   ├── router/         # 路由配置
│   │   ├── utils/          # 工具函数
│   │   └── styles/         # 样式文件
│   └── package.json
├── server/                 # 后端 Go 服务
│   ├── cmd/                # 入口文件
│   ├── config/             # 配置
│   └── internal/
│       ├── database/       # 数据库连接
│       ├── handler/        # HTTP 处理器
│       ├── middleware/     # 中间件
│       ├── model/          # 数据模型
│       ├── repository/     # 数据访问层
│       ├── router/         # 路由配置
│       ├── service/        # 业务逻辑层
│       └── websocket/      # WebSocket 处理
├── sql/                    # 数据库脚本
│   └── init.sql            # 初始化SQL
└── doc/                    # 文档目录
```

## 快速开始

### 1. 数据库初始化

```bash
mysql -u root -p < sql/init.sql
```

### 2. 启动后端服务

```bash
cd server
go mod tidy
go run cmd/main.go
```

环境变量配置（可选）：
- `SERVER_PORT`：服务端口，默认 8080
- `DB_HOST`：数据库地址，默认 127.0.0.1
- `DB_PORT`：数据库端口，默认 3306
- `DB_USER`：数据库用户，默认 root
- `DB_PASSWORD`：数据库密码
- `DB_NAME`：数据库名，默认 pinche
- `JWT_SECRET`：JWT 密钥
- `JWT_EXPIRE_HOUR`：Token 过期时间（小时），默认 168

### 3. 启动前端应用

```bash
cd app
npm install
npm run dev
```

访问 http://localhost:3000

## API 接口

### 用户模块
- `POST /api/user/register` - 用户注册
- `POST /api/user/login` - 用户登录
- `GET /api/user/profile` - 获取个人信息
- `PUT /api/user/profile` - 更新个人信息

### 行程模块
- `GET /api/trips` - 获取行程列表
- `GET /api/trips/:id` - 获取行程详情
- `POST /api/trips` - 发布行程
- `GET /api/trips/my` - 获取我的行程
- `PUT /api/trips/:id/cancel` - 取消行程
- `DELETE /api/trips/:id` - 删除行程

### 匹配模块
- `GET /api/matches` - 获取我的匹配
- `GET /api/matches/:id` - 获取匹配详情
- `POST /api/matches/:id/confirm` - 确认匹配
- `GET /api/matches/:id/contact` - 获取联系方式

### 通知模块
- `GET /api/notifications` - 获取通知列表
- `PUT /api/notifications/:id/read` - 标记已读
- `PUT /api/notifications/read-all` - 全部已读

### WebSocket
- `GET /ws?token=xxx` - WebSocket 连接

## 匹配算法

匹配算法基于以下因素计算匹配得分：

1. **时间匹配**：出发时间差在 12 小时内，时间越接近得分越高
2. **路线匹配**：出发城市和目的城市必须相同
3. **位置距离**：使用 Haversine 公式计算地理距离，距离越近得分越高

匹配得分 = 基础分 × 时间系数 × 位置系数

当得分 >= 50 时，系统会自动创建匹配记录并通知双方。

## 免责声明

本平台仅提供信息匹配服务，不参与实际交易。用户请自行确认出行细节，注意人身和财产安全。如发生任何纠纷，请自行协商解决或寻求法律途径。
