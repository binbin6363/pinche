-- 为 messages 表添加 duration 字段（语音消息时长）
ALTER TABLE messages ADD COLUMN duration INT NOT NULL DEFAULT 0 COMMENT '语音消息时长(秒)' AFTER msg_type;

-- 更新 msg_type 注释
ALTER TABLE messages MODIFY COLUMN msg_type TINYINT NOT NULL DEFAULT 1 COMMENT '消息类型: 1-文字 2-图片 3-语音 4-表情';
