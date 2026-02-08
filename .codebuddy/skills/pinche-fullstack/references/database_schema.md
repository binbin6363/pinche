# 数据库表结构

## 表概览

| 表名 | 说明 |
|------|------|
| users | 用户信息 |
| trips | 行程信息 |
| matches | 匹配记录 |
| messages | 私聊消息 |
| notifications | 通知 |
| announcements | 公告 |
| trip_grabs | 抢单记录 |

## users 用户表

```sql
CREATE TABLE users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    open_id VARCHAR(64) NOT NULL UNIQUE COMMENT '对外暴露的用户标识',
    phone VARCHAR(20) COMMENT '手机号',
    nickname VARCHAR(64) NOT NULL COMMENT '昵称',
    avatar VARCHAR(255) COMMENT '头像URL',
    identity TINYINT DEFAULT 0 COMMENT '身份：0=未设置，1=司机，2=乘客',
    status TINYINT DEFAULT 0 COMMENT '状态：0=正常，1=封禁',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_phone (phone),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

## trips 行程表

```sql
CREATE TABLE trips (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL COMMENT '发布者ID',
    trip_type TINYINT NOT NULL COMMENT '类型：1=车找人，2=人找车',
    departure_city VARCHAR(64) NOT NULL COMMENT '出发城市',
    departure_address VARCHAR(255) COMMENT '出发详细地址',
    destination_city VARCHAR(64) NOT NULL COMMENT '目的城市',
    destination_address VARCHAR(255) COMMENT '目的详细地址',
    departure_time DATETIME NOT NULL COMMENT '出发时间',
    seat_count INT DEFAULT 1 COMMENT '座位数/人数',
    price DECIMAL(10,2) COMMENT '价格',
    description TEXT COMMENT '备注说明',
    status TINYINT DEFAULT 1 COMMENT '状态：1=待匹配，2=已匹配，3=已完成，4=已取消，5=已封禁',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_trip_type (trip_type),
    INDEX idx_status (status),
    INDEX idx_departure_time (departure_time),
    INDEX idx_departure_city (departure_city),
    INDEX idx_destination_city (destination_city),
    FOREIGN KEY (user_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

## matches 匹配表

```sql
CREATE TABLE matches (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    trip_id BIGINT UNSIGNED NOT NULL COMMENT '行程ID',
    requester_id BIGINT UNSIGNED NOT NULL COMMENT '请求者ID',
    owner_id BIGINT UNSIGNED NOT NULL COMMENT '行程发布者ID',
    status TINYINT DEFAULT 0 COMMENT '状态：0=待确认，1=成功，2=失败',
    message TEXT COMMENT '留言',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_trip_id (trip_id),
    INDEX idx_requester_id (requester_id),
    INDEX idx_owner_id (owner_id),
    INDEX idx_status (status),
    FOREIGN KEY (trip_id) REFERENCES trips(id),
    FOREIGN KEY (requester_id) REFERENCES users(id),
    FOREIGN KEY (owner_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

## messages 消息表

```sql
CREATE TABLE messages (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    sender_id BIGINT UNSIGNED NOT NULL COMMENT '发送者ID',
    receiver_id BIGINT UNSIGNED NOT NULL COMMENT '接收者ID',
    content TEXT NOT NULL COMMENT '消息内容',
    is_read TINYINT DEFAULT 0 COMMENT '是否已读：0=未读，1=已读',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_sender_id (sender_id),
    INDEX idx_receiver_id (receiver_id),
    INDEX idx_is_read (is_read),
    FOREIGN KEY (sender_id) REFERENCES users(id),
    FOREIGN KEY (receiver_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

## notifications 通知表

```sql
CREATE TABLE notifications (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    type VARCHAR(32) NOT NULL COMMENT '通知类型',
    title VARCHAR(128) NOT NULL COMMENT '标题',
    content TEXT COMMENT '内容',
    data JSON COMMENT '附加数据',
    is_read TINYINT DEFAULT 0 COMMENT '是否已读',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_is_read (is_read),
    FOREIGN KEY (user_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

## announcements 公告表

```sql
CREATE TABLE announcements (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(128) NOT NULL COMMENT '标题',
    content TEXT NOT NULL COMMENT '内容',
    type TINYINT DEFAULT 1 COMMENT '类型：1=普通，2=紧急，3=系统',
    is_active TINYINT DEFAULT 1 COMMENT '是否启用',
    sort_order INT DEFAULT 0 COMMENT '排序权重',
    start_time DATETIME COMMENT '开始展示时间',
    end_time DATETIME COMMENT '结束展示时间',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_is_active (is_active),
    INDEX idx_start_time (start_time),
    INDEX idx_end_time (end_time)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

## trip_grabs 抢单表

```sql
CREATE TABLE trip_grabs (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    trip_id BIGINT UNSIGNED NOT NULL COMMENT '行程ID',
    user_id BIGINT UNSIGNED NOT NULL COMMENT '抢单用户ID',
    message TEXT COMMENT '留言',
    status TINYINT DEFAULT 0 COMMENT '状态：0=待处理，1=接受，2=拒绝',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_trip_user (trip_id, user_id),
    INDEX idx_trip_id (trip_id),
    INDEX idx_user_id (user_id),
    INDEX idx_status (status),
    FOREIGN KEY (trip_id) REFERENCES trips(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

## 状态码说明

### 用户状态 (users.status)
| 值 | 说明 |
|----|------|
| 0 | 正常 |
| 1 | 封禁 |

### 行程状态 (trips.status)
| 值 | 说明 |
|----|------|
| 1 | 待匹配 |
| 2 | 已匹配 |
| 3 | 已完成 |
| 4 | 已取消 |
| 5 | 已封禁 |

### 匹配状态 (matches.status)
| 值 | 说明 |
|----|------|
| 0 | 待确认 |
| 1 | 成功 |
| 2 | 失败 |

### 行程类型 (trips.trip_type)
| 值 | 说明 |
|----|------|
| 1 | 车找人（司机发布） |
| 2 | 人找车（乘客发布） |

## 设计要点

1. **open_id vs id**：对外接口使用 `open_id` 标识用户，隐藏内部自增 `id`
2. **软删除**：使用 `status` 字段标记封禁/取消，不物理删除
3. **时间索引**：`departure_time` 添加索引优化时间范围查询
4. **城市索引**：出发/目的城市添加索引优化路线匹配
