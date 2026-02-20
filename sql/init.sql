-- 拼车应用数据库初始化脚本
-- 创建数据库
CREATE DATABASE IF NOT EXISTS pinche DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE pinche;

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '用户ID(内部)',
    open_id VARCHAR(32) NOT NULL COMMENT '用户OpenID(对外)',
    phone VARCHAR(20) NOT NULL COMMENT '手机号',
    password VARCHAR(128) NOT NULL COMMENT '密码哈希',
    nickname VARCHAR(50) NOT NULL DEFAULT '' COMMENT '昵称',
    avatar VARCHAR(512) NOT NULL DEFAULT '' COMMENT '头像URL',
    gender TINYINT NOT NULL DEFAULT 0 COMMENT '性别: 0-未知 1-男 2-女',
    status TINYINT NOT NULL DEFAULT 0 COMMENT '状态: 0-正常 1-封禁',
    city VARCHAR(50) NOT NULL DEFAULT '' COMMENT '城市',
    province VARCHAR(50) NOT NULL DEFAULT '' COMMENT '省份',
    -- contact info
    contact_phone VARCHAR(20) NOT NULL DEFAULT '' COMMENT '联系手机号',
    contact_wechat VARCHAR(50) NOT NULL DEFAULT '' COMMENT '联系微信号',
    -- emergency contact
    emergency_contact_name VARCHAR(50) NOT NULL DEFAULT '' COMMENT '紧急联系人姓名',
    emergency_contact_phone VARCHAR(20) NOT NULL DEFAULT '' COMMENT '紧急联系人电话',
    emergency_contact_relation VARCHAR(20) NOT NULL DEFAULT '' COMMENT '紧急联系人关系',
    -- car info
    car_number VARCHAR(10) NOT NULL DEFAULT '' COMMENT '车牌号',
    car_brand VARCHAR(50) NOT NULL DEFAULT '' COMMENT '车辆品牌',
    car_model VARCHAR(50) NOT NULL DEFAULT '' COMMENT '车辆型号',
    car_color VARCHAR(20) NOT NULL DEFAULT '' COMMENT '车辆颜色',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_open_id (open_id),
    UNIQUE KEY uk_phone (phone),
    KEY idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- 行程表 (司机发布)
CREATE TABLE IF NOT EXISTS trips (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '行程ID',
    user_id BIGINT UNSIGNED NOT NULL COMMENT '发布者ID',
    trip_type TINYINT NOT NULL COMMENT '类型: 1-司机 2-乘客',
    departure_city VARCHAR(50) NOT NULL COMMENT '出发城市',
    departure_province VARCHAR(50) NOT NULL DEFAULT '' COMMENT '出发省份',
    departure_address VARCHAR(255) NOT NULL COMMENT '出发详细地址',
    departure_lat DECIMAL(10, 7) NOT NULL COMMENT '出发纬度',
    departure_lng DECIMAL(10, 7) NOT NULL COMMENT '出发经度',
    destination_city VARCHAR(50) NOT NULL COMMENT '目的城市',
    destination_province VARCHAR(50) NOT NULL DEFAULT '' COMMENT '目的省份',
    destination_address VARCHAR(255) NOT NULL COMMENT '目的详细地址',
    destination_lat DECIMAL(10, 7) NOT NULL COMMENT '目的纬度',
    destination_lng DECIMAL(10, 7) NOT NULL COMMENT '目的经度',
    departure_time DATETIME NOT NULL COMMENT '出发时间',
    seats INT NOT NULL DEFAULT 1 COMMENT '座位数(司机)/需要座位数(乘客)',
    price DECIMAL(10, 2) NOT NULL DEFAULT 0 COMMENT '费用(元)',
    remark VARCHAR(500) NOT NULL DEFAULT '' COMMENT '备注',
    images TEXT COMMENT '行程图片(JSON数组)',
    status TINYINT NOT NULL DEFAULT 1 COMMENT '状态: 1-待匹配 2-已匹配 3-已完成 4-已取消 5-已封禁',
    view_count INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '浏览次数',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    KEY idx_user_id (user_id),
    KEY idx_status (status),
    KEY idx_departure_time (departure_time),
    KEY idx_departure_city (departure_city),
    KEY idx_destination_city (destination_city),
    KEY idx_departure_province (departure_province),
    KEY idx_destination_province (destination_province)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='行程表';

-- 匹配记录表
CREATE TABLE IF NOT EXISTS matches (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '匹配ID',
    driver_trip_id BIGINT UNSIGNED NOT NULL COMMENT '司机行程ID',
    passenger_trip_id BIGINT UNSIGNED NOT NULL COMMENT '乘客行程ID',
    driver_id BIGINT UNSIGNED NOT NULL COMMENT '司机ID',
    passenger_id BIGINT UNSIGNED NOT NULL COMMENT '乘客ID',
    match_score DECIMAL(5, 2) NOT NULL DEFAULT 0 COMMENT '匹配得分',
    driver_status TINYINT NOT NULL DEFAULT 0 COMMENT '司机确认状态: 0-待确认 1-已接受 2-已拒绝',
    passenger_status TINYINT NOT NULL DEFAULT 0 COMMENT '乘客确认状态: 0-待确认 1-已接受 2-已拒绝',
    status TINYINT NOT NULL DEFAULT 0 COMMENT '匹配状态: 0-待确认 1-匹配成功 2-匹配失败',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_trips (driver_trip_id, passenger_trip_id),
    KEY idx_driver_id (driver_id),
    KEY idx_passenger_id (passenger_id),
    KEY idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='匹配记录表';

-- 通知表
CREATE TABLE IF NOT EXISTS notifications (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '通知ID',
    user_id BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    match_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '匹配ID',
    trip_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '行程ID',
    title VARCHAR(100) NOT NULL COMMENT '通知标题',
    content TEXT NOT NULL COMMENT '通知内容',
    is_read TINYINT NOT NULL DEFAULT 0 COMMENT '是否已读: 0-未读 1-已读',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (id),
    KEY idx_user_id (user_id),
    KEY idx_match_id (match_id),
    KEY idx_trip_id (trip_id),
    KEY idx_is_read (is_read)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='通知表';

-- 私聊消息表
CREATE TABLE IF NOT EXISTS messages (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '消息ID',
    sender_id BIGINT UNSIGNED NOT NULL COMMENT '发送者ID',
    receiver_id BIGINT UNSIGNED NOT NULL COMMENT '接收者ID',
    content TEXT NOT NULL COMMENT '消息内容(文字或图片URL)',
    msg_type TINYINT NOT NULL DEFAULT 1 COMMENT '消息类型: 1-文字 2-图片 3-系统消息',
    is_read TINYINT NOT NULL DEFAULT 0 COMMENT '是否已读: 0-未读 1-已读',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (id),
    KEY idx_sender_id (sender_id),
    KEY idx_receiver_id (receiver_id),
    KEY idx_conversation (sender_id, receiver_id, created_at),
    KEY idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='私聊消息表';

-- 公告表
CREATE TABLE IF NOT EXISTS announcements (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '公告ID',
    title VARCHAR(100) NOT NULL COMMENT '标题',
    content TEXT NOT NULL COMMENT '内容',
    type TINYINT NOT NULL DEFAULT 1 COMMENT '类型: 1-普通 2-重要 3-紧急',
    is_active TINYINT NOT NULL DEFAULT 1 COMMENT '是否启用: 0-禁用 1-启用',
    sort_order INT NOT NULL DEFAULT 0 COMMENT '排序权重',
    start_time DATETIME NOT NULL COMMENT '开始时间',
    end_time DATETIME NOT NULL COMMENT '结束时间',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    KEY idx_is_active (is_active),
    KEY idx_start_time (start_time),
    KEY idx_end_time (end_time)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='公告表';

-- 行程抢单记录表
CREATE TABLE IF NOT EXISTS trip_grabs (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '抢单ID',
    trip_id BIGINT UNSIGNED NOT NULL COMMENT '行程ID',
    user_id BIGINT UNSIGNED NOT NULL COMMENT '抢单用户ID',
    message VARCHAR(200) NOT NULL DEFAULT '' COMMENT '留言',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_trip_user (trip_id, user_id),
    KEY idx_trip_id (trip_id),
    KEY idx_user_id (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='行程抢单记录表';

-- 行程修改审核表
CREATE TABLE IF NOT EXISTS trip_updates (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '审核ID',
    trip_id BIGINT UNSIGNED NOT NULL COMMENT '行程ID',
    user_id BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    update_type TINYINT NOT NULL COMMENT '修改类型: 1-起终点 2-时间',
    old_value TEXT NOT NULL COMMENT '原值(JSON)',
    new_value TEXT NOT NULL COMMENT '新值(JSON)',
    status TINYINT NOT NULL DEFAULT 0 COMMENT '状态: 0-待审核 1-已通过 2-已拒绝',
    reject_reason VARCHAR(200) NOT NULL DEFAULT '' COMMENT '拒绝原因',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    KEY idx_trip_id (trip_id),
    KEY idx_user_id (user_id),
    KEY idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='行程修改审核表';

-- 插入系统用户 (用于系统通知)
INSERT INTO users (id, open_id, phone, password, nickname, avatar, gender, status) VALUES 
(1, 'system_000000000000000000', '00000000000', '', '系统通知', '', 0, 0)
ON DUPLICATE KEY UPDATE nickname = '系统通知';

-- 插入默认公告
INSERT INTO announcements (title, content, type, is_active, sort_order, start_time, end_time) VALUES 
('欢迎使用春节拼车平台', '本平台为公益性质，不收取任何费用。请在"我的"页面选择您的身份（司机/乘客），然后在"行程"页面查看匹配信息。祝您旅途愉快！', 1, 1, 100, '2024-01-01 00:00:00', '2026-12-31 23:59:59')
ON DUPLICATE KEY UPDATE content = content;
