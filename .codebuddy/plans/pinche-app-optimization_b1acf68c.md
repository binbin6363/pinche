---
name: pinche-app-optimization
overview: 全面优化免费顺风车拼车应用的前端页面和后端接口，重构首页、行程页、匹配页和个人中心页面，实现身份切换、头像修改、聊天列表等功能。
todos:
  - id: refactor-home
    content: 改版首页Home.vue，移除四个快捷入口按钮，保留最新动态列表，新增平台公告区域
    status: completed
  - id: add-identity-store
    content: 扩展user.js状态管理，新增identity身份字段和setIdentity方法，支持持久化
    status: completed
  - id: refactor-trips
    content: 改版行程页Trips.vue，移除筛选Tab，根据用户身份自动筛选对立方行程
    status: completed
    dependencies:
      - add-identity-store
  - id: refactor-matches-to-messages
    content: 重写Matches.vue为消息会话列表页，调用fetchConversations展示聊天记录
    status: completed
  - id: refactor-profile
    content: 改版Profile.vue，新增身份切换功能和头像修改功能
    status: completed
    dependencies:
      - add-identity-store
  - id: update-main-layout
    content: 更新MainLayout.vue底部导航，将"匹配"改为"消息"
    status: completed
    dependencies:
      - refactor-matches-to-messages
---

## 用户需求

优化免费顺风车拼车应用的 UI 和功能体验，主要改版内容：

## 核心功能

### 1. 首页改版

- 去掉四个快捷入口大按钮（司机发布、找乘客、找司机、乘客发布）
- 保留最新匹配动态展示区
- 新增平台重要通知/公告展示区
- 保留温馨提示区域

### 2. 行程页面改版

- 去掉"全部""司机""乘客"三个筛选 Tab
- 根据当前用户身份自动展示对立方行程（司机看乘客行程、乘客看司机行程）
- 行程卡片严禁显示用户联系方式（手机号等）
- 点击行程卡片可进入详情页，支持发起私聊

### 3. 匹配页面改版为"消息"页面

- 原"匹配"Tab 改为"消息"功能
- 展示聊天会话列表（私聊记录）
- 显示最后一条消息预览、未读数、时间

### 4. 我的页面改版

- 新增身份切换功能（司机/乘客）
- 新增头像修改功能（点击头像可更换）
- 保留我的行程列表
- 保留消息通知入口
- 保留退出登录功能

### 5. 信息发布与展示分离

- 发布页独立，展示页清晰呈现行程必要信息
- 行程卡片展示：出发地、目的地、时间、座位数、费用
- 不展示发布者联系方式

## 技术栈

- 前端：Vue 3 + Vite + Tailwind CSS + Pinia
- 后端：Go + Gin（已有完整 API 支持）
- 通信：WebSocket（已有实时消息推送）

## 实现方案

### 1. 首页改版

- 移除 `Home.vue` 中的四个快捷入口按钮区域
- 保留最新行程列表，改为"最新动态"模块
- 新增"平台公告"卡片区域展示重要通知

### 2. 行程页面改版

- 移除 `Trips.vue` 中的筛选 Tab（全部/司机/乘客按钮）
- 新增用户身份状态管理（`user.js` 添加 `identity` 字段：1=司机，2=乘客）
- 根据身份自动设置 `trip_type` 参数（司机身份查 trip_type=2，乘客身份查 trip_type=1）
- 行程卡片移除联系方式展示，保留私聊入口

### 3. 匹配页面改版为消息页面

- 重写 `Matches.vue` 为会话列表页面
- 调用已有的 `messageStore.fetchConversations()` 获取会话列表
- 展示对方头像、昵称、最后消息、未读数、时间
- 点击进入 Chat 页面

### 4. 我的页面改版

- 新增身份切换组件（司机/乘客切换按钮）
- 头像区域改为可点击，支持选择图片上传
- 调用已有的 `userStore.updateProfile({ avatar })` 更新头像

### 5. 状态管理扩展

- `user.js` 新增 `identity` 字段及 `setIdentity` 方法
- 身份信息持久化到 localStorage

## 实现要点

1. **身份状态设计**：用户身份（identity）存储在前端 userStore，支持切换，默认值为 2（乘客）
2. **行程筛选逻辑**：Trips 页面根据 identity 自动计算对立方 trip_type 进行筛选
3. **头像上传**：使用 FileReader 转 base64（演示），生产环境需接入 OSS
4. **消息页面**：复用已有 message store 的 fetchConversations API
5. **隐私保护**：行程列表/详情页不展示用户手机号等联系方式

## 目录结构

```
app/src/
├── views/
│   ├── Home.vue           # [MODIFY] 移除快捷入口，保留动态和公告
│   ├── Trips.vue          # [MODIFY] 移除筛选Tab，根据身份自动筛选
│   ├── Matches.vue        # [MODIFY] 重写为消息会话列表页
│   └── Profile.vue        # [MODIFY] 新增身份切换和头像修改
├── stores/
│   └── user.js            # [MODIFY] 新增identity字段和setIdentity方法
└── layouts/
    └── MainLayout.vue     # [MODIFY] 底部导航"匹配"改为"消息"图标和文案
```