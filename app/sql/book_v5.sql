-- =====================================================================
-- book_v5.sql —— 绿联 NAS 登录认证支持
-- Engine: InnoDB / utf8mb4_unicode_ci
-- Dependency: 依赖 system-manage.sql (sys_user)
-- Notes:
--   1. sys_user 表新增 ugreen_user_id 字段，用于映射绿联NAS用户
--   2. 配合后端 UgreenAuth 中间件和 UgreenAuthService 使用
-- =====================================================================

-- 1. sys_user 新增 ugreen_user_id 字段
ALTER TABLE sys_user
    ADD COLUMN ugreen_user_id VARCHAR(64) DEFAULT NULL COMMENT '绿联NAS用户ID' AFTER avatar;

CREATE INDEX idx_sys_user_ugreen_user_id ON sys_user(ugreen_user_id);

