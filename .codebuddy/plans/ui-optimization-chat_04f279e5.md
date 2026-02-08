---
name: ui-optimization-chat
overview: 优化登录/注册页面布局，实现行程列表按身份默认筛选对立方信息，以及在行程详情页增加即时通讯私聊功能（支持文字和图片消息）
design:
  architecture:
    framework: vue
  styleKeywords:
    - 移动端优先
    - 简约现代
    - 触摸友好
    - 气泡对话
  fontSystem:
    fontFamily: system-ui
    heading:
      size: 18px
      weight: 600
    subheading:
      size: 16px
      weight: 500
    body:
      size: 14px
      weight: 400
  colorSystem:
    primary:
      - "#3B82F6"
      - "#60A5FA"
    background:
      - "#F9FAFB"
      - "#FFFFFF"
      - "#F3F4F6"
    text:
      - "#1F2937"
      - "#6B7280"
      - "#9CA3AF"
    functional:
      - "#10B981"
      - "#EF4444"
todos:
  - id: optimize-login-register
    content: 优化Login.vue和Register.vue布局，将Logo区域紧凑化，表单上移
    status: completed
  - id: smart-filter-trips
    content: 修改Home.vue跳转逻辑和Trips.vue，实现根据身份默认筛选对立方行程
    status: completed
  - id: backend-message-model
    content: 创建后端消息模型message.go和数据库表结构
    status: completed
  - id: backend-message-repo-service
    content: 实现message_repo.go和message_service.go，完成消息CRUD操作
    status: completed
    dependencies:
      - backend-message-model
  - id: backend-message-handler
    content: 实现message_handler.go并注册路由，扩展hub.go支持私聊消息转发
    status: completed
    dependencies:
      - backend-message-repo-service
  - id: frontend-chat-page
    content: 创建Chat.vue聊天页面和message.js状态管理，扩展websocket.js
    status: completed
    dependencies:
      - backend-message-handler
  - id: integrate-chat-entry
    content: 在TripDetail.vue增加私聊入口，添加聊天页面路由
    status: completed
    dependencies:
      - frontend-chat-page
---

## 用户需求

针对现有拼车应用进行三项UI功能优化：

## 核心功能

### 1. 登录/注册页面布局优化

- 将输入表单区域上移，减少顶部Logo占用空间
- 符合现代移动端设计规范，提升单手操作便捷性
- 输入框更靠近屏幕中部位置，避免键盘遮挡

### 2. 行程列表智能筛选

- 用户选择身份（司机/乘客）后进入行程列表页面
- 默认展示对立方发布的信息：乘客身份显示司机行程，司机身份显示乘客行程
- 避免用户看到同类型信息造成混淆

### 3. 行程详情页即时通讯功能

- 在行程详情页增加与发布者私聊入口
- 支持文字消息发送和接收
- 支持图片消息发送和接收
- 基于现有WebSocket实现实时消息推送

## 技术栈

- 前端：Vue 3 + Vite + Tailwind CSS + Pinia
- 后端：Go + Gin + WebSocket + MySQL
- 无第三方UI组件库，使用自定义Tailwind样式

## 实现方案

### 1. 登录/注册页面布局优化

**方案**：调整flex布局比例，将Logo区域从flex-1改为固定高度，表单区域占据更多空间

- Login.vue: Logo区域使用固定pt-8，移除flex-1
- Register.vue: 同样优化，Logo区域紧凑化

### 2. 行程列表智能筛选

**方案**：通过路由参数传递用户选择的身份，Trips页面根据身份设置默认筛选条件

- Home.vue中"我是司机/乘客"按钮点击后，跳转到Trips页面并传递identity参数
- Trips.vue接收identity参数，设置filter.trip_type为对立方类型（司机选择后显示乘客行程trip_type=2，反之亦然）

### 3. 即时通讯功能

**方案**：扩展现有WebSocket实现私聊功能

**后端扩展**：

- 新增Message模型：存储私聊消息（sender_id, receiver_id, content, msg_type, created_at）
- 新增message_repo.go：消息数据访问层
- 新增message_service.go：消息业务逻辑
- 新增message_handler.go：消息API接口
- 扩展websocket/hub.go：支持私聊消息转发

**前端扩展**：

- 新增Chat.vue页面：私聊界面，支持文字和图片消息
- 新增message store：消息状态管理
- 扩展websocket.js：处理私聊消息类型
- TripDetail.vue：增加"联系发布者"按钮

## 实现注意事项

### 性能考虑

- 消息列表采用分页加载，避免一次性加载大量历史消息
- 图片上传前进行压缩，限制最大尺寸
- WebSocket消息采用JSON格式，保持轻量

### 安全考虑

- 私聊消息仅双方可见，查询时校验用户权限
- 图片上传校验文件类型和大小

## 目录结构

```
app/src/
├── views/
│   ├── Login.vue           # [MODIFY] 优化布局，Logo区域紧凑化
│   ├── Register.vue        # [MODIFY] 优化布局，Logo区域紧凑化
│   ├── Home.vue            # [MODIFY] 点击身份按钮跳转至Trips并传递identity参数
│   ├── Trips.vue           # [MODIFY] 接收identity参数，默认筛选对立方行程
│   ├── TripDetail.vue      # [MODIFY] 增加联系发布者按钮
│   └── Chat.vue            # [NEW] 私聊页面，支持文字和图片消息
├── stores/
│   └── message.js          # [NEW] 消息状态管理
├── utils/
│   └── websocket.js        # [MODIFY] 扩展支持私聊消息类型
└── router/
    └── index.js            # [MODIFY] 添加聊天页面路由

server/internal/
├── model/
│   └── message.go          # [NEW] Message模型定义
├── repository/
│   └── message_repo.go     # [NEW] 消息数据访问层
├── service/
│   └── message_service.go  # [NEW] 消息业务逻辑
├── handler/
│   └── message_handler.go  # [NEW] 消息API处理器
├── websocket/
│   └── hub.go              # [MODIFY] 扩展私聊消息处理
└── router/
    └── router.go           # [MODIFY] 添加消息相关路由

sql/
└── init.sql                # [MODIFY] 添加messages表结构
```

## 设计风格

延续现有应用的现代简约风格，保持一致性。采用移动端优先设计，确保触摸友好。

## 页面设计

### 1. 登录/注册页面优化

- Logo区域：缩小至顶部20%空间，使用固定padding而非flex-1
- 表单区域：白色圆角卡片从屏幕40%位置开始，给输入框更多可见空间
- 保持渐变蓝色背景，白色圆角表单卡片的视觉层次

### 2. 聊天页面设计

- **顶部导航栏**：返回按钮、对方昵称、更多操作
- **消息列表区域**：
- 自己消息：右侧蓝色气泡
- 对方消息：左侧灰色气泡
- 图片消息：圆角缩略图，点击可放大
- 时间分隔线：每隔5分钟或新日期显示时间
- **底部输入区域**：
- 左侧图片选择按钮
- 中间文字输入框
- 右侧发送按钮
- 安全区域底部padding

### 3. 行程详情页新增元素

- 发布者信息卡片下方增加"立即私聊"按钮
- 使用主色蓝色，与整体风格一致