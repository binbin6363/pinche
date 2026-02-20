-- 用户表添加联系方式、紧急联系人、车辆信息字段
-- 执行前请备份数据库

USE pinche;

-- contact info
ALTER TABLE users ADD COLUMN contact_phone VARCHAR(20) NOT NULL DEFAULT '' COMMENT '联系手机号' AFTER province;
ALTER TABLE users ADD COLUMN contact_wechat VARCHAR(50) NOT NULL DEFAULT '' COMMENT '联系微信号' AFTER contact_phone;

-- emergency contact
ALTER TABLE users ADD COLUMN emergency_contact_name VARCHAR(50) NOT NULL DEFAULT '' COMMENT '紧急联系人姓名' AFTER contact_wechat;
ALTER TABLE users ADD COLUMN emergency_contact_phone VARCHAR(20) NOT NULL DEFAULT '' COMMENT '紧急联系人电话' AFTER emergency_contact_name;
ALTER TABLE users ADD COLUMN emergency_contact_relation VARCHAR(20) NOT NULL DEFAULT '' COMMENT '紧急联系人关系' AFTER emergency_contact_phone;

-- car info
ALTER TABLE users ADD COLUMN car_number VARCHAR(10) NOT NULL DEFAULT '' COMMENT '车牌号' AFTER emergency_contact_relation;
ALTER TABLE users ADD COLUMN car_brand VARCHAR(50) NOT NULL DEFAULT '' COMMENT '车辆品牌' AFTER car_number;
ALTER TABLE users ADD COLUMN car_model VARCHAR(50) NOT NULL DEFAULT '' COMMENT '车辆型号' AFTER car_brand;
ALTER TABLE users ADD COLUMN car_color VARCHAR(20) NOT NULL DEFAULT '' COMMENT '车辆颜色' AFTER car_model;
