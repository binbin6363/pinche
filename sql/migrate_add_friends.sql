-- 好友关系表迁移脚本
-- 用于存储用户之间的好友申请和好友关系

USE pinche;

-- 好友关系表
CREATE TABLE IF NOT EXISTS friends (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '记录ID',
    user_id BIGINT UNSIGNED NOT NULL COMMENT '申请者ID(内部)',
    friend_id BIGINT UNSIGNED NOT NULL COMMENT '被申请者ID(内部)',
    status TINYINT NOT NULL DEFAULT 0 COMMENT '状态: 0-待确认 1-已同意 2-已拒绝',
    message VARCHAR(200) NOT NULL DEFAULT '' COMMENT '申请留言',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_user_friend (user_id, friend_id),
    KEY idx_friend_id (friend_id),
    KEY idx_status (status),
    KEY idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='好友关系表';
