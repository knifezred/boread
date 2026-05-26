-- =====================================================================
-- System Manage Schema (v5) —— RBAC + 数据权限
-- Source:  /app/web/src/typings/api/system-manage.d.ts
-- Engine:  InnoDB / utf8mb4_unicode_ci
-- Requires: MySQL 8.0.13+ (使用函数索引 IFNULL(deleted_at,'1970-01-01'))
-- Scope:   仅后台管理（部门、角色、菜单、按钮、用户、登录日志、操作日志、字典）
--          前台阅读者表 reader 见 business.sql
-- Notes:
--   1. 所有业务表统一加 deleted_at 软删字段（GORM 约定）
--   2. 关联表统一使用 id 外键
--   3. 时间字段 DATETIME(3) 毫秒精度
--   4. 数据权限：sys_role.data_scope + sys_role_dept
--   5. create_by / update_by 字段名保留 (兼容前端约定),
--      类型 BIGINT 存用户 id, 避免用户改名后历史脏数据
--   6. is_system 标记系统内置数据, 防止误删 (角色/菜单/字典)
--   7. 唯一索引一律 (业务键, IFNULL(deleted_at,'1970-01-01')) 函数索引,
--      避免"软删后无法用同名重建"的经典坑
-- =====================================================================

CREATE DATABASE IF NOT EXISTS `boread`
  DEFAULT CHARACTER SET utf8mb4
  DEFAULT COLLATE utf8mb4_unicode_ci;

USE `boread`;

-- ---------------------------------------------------------------------
-- Table: sys_dept (部门 / 组织)
-- 自关联树形结构, parent_id = 0 为顶层
-- ---------------------------------------------------------------------
DROP TABLE IF EXISTS `sys_dept`;
CREATE TABLE `sys_dept` (
  `id`          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT          COMMENT '主键 id',
  `parent_id`   BIGINT UNSIGNED NOT NULL DEFAULT 0                COMMENT '父部门 id, 顶层为 0',
  `ancestors`   VARCHAR(512)    NOT NULL DEFAULT ''               COMMENT '祖级链: 0,1,3 (便于查全部下级, 配合 LIKE prefix% 走前缀索引)',
  `dept_name`   VARCHAR(64)     NOT NULL                          COMMENT '部门名称',
  `dept_code`   VARCHAR(64)     NOT NULL                          COMMENT '部门编码',
  `leader`      VARCHAR(64)     NULL     DEFAULT NULL             COMMENT '部门负责人',
  `sort_order`  INT             NOT NULL DEFAULT 0                COMMENT '排序',
  `status`      CHAR(1)         NOT NULL DEFAULT '1'              COMMENT '状态: 1-启用, 2-禁用',
  `create_by`   BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '创建人 (存 sys_user.id, 字段名保留兼容前端)',
  `create_time` DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `update_by`   BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '更新人 (存 sys_user.id)',
  `update_time` DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at`  DATETIME(3)     NULL     DEFAULT NULL             COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_dept_code` (`dept_code`, ((IFNULL(`deleted_at`, '1970-01-01 00:00:00')))),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_ancestors` (`ancestors`(64)),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='部门表';

-- ---------------------------------------------------------------------
-- Table: sys_role
-- 含数据权限范围 data_scope, is_system 防止系统内置角色被误删
-- ---------------------------------------------------------------------
DROP TABLE IF EXISTS `sys_role`;
CREATE TABLE `sys_role` (
  `id`          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT          COMMENT '主键 id',
  `role_name`   VARCHAR(64)     NOT NULL                          COMMENT '角色名称',
  `role_code`   VARCHAR(64)     NOT NULL                          COMMENT '角色编码',
  `role_desc`   VARCHAR(255)    NULL     DEFAULT NULL             COMMENT '角色描述',
  `data_scope`  CHAR(1)         NOT NULL DEFAULT '5'              COMMENT '数据权限: 1-全部, 2-自定义部门, 3-本部门, 4-本部门及子部门, 5-仅本人',
  `is_system`   TINYINT(1)      NOT NULL DEFAULT 0                COMMENT '是否系统内置(1=不可删/不可改 code)',
  `sort_order`  INT             NOT NULL DEFAULT 0                COMMENT '排序',
  `status`      CHAR(1)         NOT NULL DEFAULT '1'              COMMENT '状态: 1-启用, 2-禁用',
  `create_by`   BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '创建人 (存 sys_user.id)',
  `create_time` DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `update_by`   BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '更新人 (存 sys_user.id)',
  `update_time` DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at`  DATETIME(3)     NULL     DEFAULT NULL             COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_role_code` (`role_code`, ((IFNULL(`deleted_at`, '1970-01-01 00:00:00')))),
  KEY `idx_role_name` (`role_name`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色表';

-- ---------------------------------------------------------------------
-- Table: sys_role_dept
-- 角色-部门 关联表 (仅 data_scope=2 自定义部门时使用)
-- ---------------------------------------------------------------------
DROP TABLE IF EXISTS `sys_role_dept`;
CREATE TABLE `sys_role_dept` (
  `id`        BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键 id',
  `role_id`   BIGINT UNSIGNED NOT NULL                COMMENT '角色 id',
  `dept_id`   BIGINT UNSIGNED NOT NULL                COMMENT '部门 id',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_role_dept` (`role_id`, `dept_id`),
  KEY `idx_dept_id` (`dept_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色数据权限部门关联表';

-- ---------------------------------------------------------------------
-- Table: sys_user (后台管理用户)
-- 含: 乐观锁 version / 登录风控 (pwd_error_count + locked_until)
--     密码策略 (pwd_updated_at 用于过期判断)
-- 字段命名说明: user_* 前缀与前端 Api.SystemManage.User 类型对齐,
--               nick_name 历史命名, 不加 user_ 前缀 (前后端均沿用)
-- 前台阅读者请见 business.sql -> reader 表
-- ---------------------------------------------------------------------
DROP TABLE IF EXISTS `sys_user`;
CREATE TABLE `sys_user` (
  `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT          COMMENT '主键 id',
  `dept_id`         BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '所属部门 id',
  `user_name`       VARCHAR(64)     NOT NULL                          COMMENT '用户名 (登录账号)',
  `password`        VARCHAR(128)    NOT NULL                          COMMENT '密码 (bcrypt)',
  `pwd_updated_at`  DATETIME(3)     NULL     DEFAULT NULL             COMMENT '密码最后修改时间 (用于过期策略)',
  `pwd_error_count` SMALLINT UNSIGNED NOT NULL DEFAULT 0              COMMENT '连续密码错误次数 (登录成功重置)',
  `locked_until`    DATETIME(3)     NULL     DEFAULT NULL             COMMENT '账号锁定到何时 (NULL=未锁)',
  `user_gender`     CHAR(1)         NULL     DEFAULT NULL             COMMENT '性别: 1-男, 2-女',
  `nick_name`       VARCHAR(64)     NOT NULL DEFAULT ''               COMMENT '昵称',
  `user_phone`      VARCHAR(20)     NULL     DEFAULT NULL             COMMENT '手机号',
  `user_email`      VARCHAR(128)    NULL     DEFAULT NULL             COMMENT '邮箱',
  `avatar`          VARCHAR(255)    NULL     DEFAULT NULL             COMMENT '头像 URL',
  `last_login_time` DATETIME(3)     NULL     DEFAULT NULL             COMMENT '最后登录时间',
  `last_login_ip`   VARCHAR(64)     NULL     DEFAULT NULL             COMMENT '最后登录 IP',
  `status`          CHAR(1)         NOT NULL DEFAULT '1'              COMMENT '状态: 1-启用, 2-禁用',
  `version`         BIGINT UNSIGNED NOT NULL DEFAULT 0                COMMENT '乐观锁版本号 (GORM optimistic_locking)',
  `create_by`       BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '创建人 (存 sys_user.id)',
  `create_time`     DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `update_by`       BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '更新人 (存 sys_user.id)',
  `update_time`     DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at`      DATETIME(3)     NULL     DEFAULT NULL             COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_name` (`user_name`, ((IFNULL(`deleted_at`, '1970-01-01 00:00:00')))),
  KEY `idx_dept_id` (`dept_id`),
  KEY `idx_user_phone` (`user_phone`),
  KEY `idx_user_email` (`user_email`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='后台用户表';

-- ---------------------------------------------------------------------
-- Table: sys_user_role
-- 关联表: 直接物理删, 无 deleted_at
-- ---------------------------------------------------------------------
DROP TABLE IF EXISTS `sys_user_role`;
CREATE TABLE `sys_user_role` (
  `id`          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键 id',
  `user_id`     BIGINT UNSIGNED NOT NULL                COMMENT '用户 id',
  `role_id`     BIGINT UNSIGNED NOT NULL                COMMENT '角色 id',
  `create_time` DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_role` (`user_id`, `role_id`),
  KEY `idx_role_id` (`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户角色关联表';

-- ---------------------------------------------------------------------
-- Table: sys_menu
-- is_system 防止系统内置菜单被误删
-- ---------------------------------------------------------------------
DROP TABLE IF EXISTS `sys_menu`;
CREATE TABLE `sys_menu` (
  `id`                  BIGINT UNSIGNED NOT NULL AUTO_INCREMENT          COMMENT '主键 id',
  `parent_id`           BIGINT UNSIGNED NOT NULL DEFAULT 0                COMMENT '父菜单 id, 顶层为 0',
  `menu_type`           CHAR(1)         NOT NULL                          COMMENT '菜单类型: 1-目录, 2-菜单',
  `menu_name`           VARCHAR(64)     NOT NULL                          COMMENT '菜单名称',
  `route_name`          VARCHAR(128)    NOT NULL                          COMMENT '路由名称',
  `route_path`          VARCHAR(255)    NOT NULL                          COMMENT '路由路径',
  `component`           VARCHAR(255)    NULL     DEFAULT NULL             COMMENT '组件路径',
  `icon`                VARCHAR(128)    NULL     DEFAULT NULL             COMMENT '图标名称',
  `icon_type`           CHAR(1)         NOT NULL DEFAULT '1'              COMMENT '图标类型: 1-iconify, 2-local',
  `i18n_key`            VARCHAR(255)    NULL     DEFAULT NULL             COMMENT '国际化 key',
  `keep_alive`          TINYINT(1)      NOT NULL DEFAULT 0                COMMENT '是否缓存路由',
  `constant`            TINYINT(1)      NOT NULL DEFAULT 0                COMMENT '是否为常量路由',
  `sort_order`          INT             NOT NULL DEFAULT 0                COMMENT '排序',
  `href`                VARCHAR(255)    NULL     DEFAULT NULL             COMMENT '外链地址',
  `hide_in_menu`        TINYINT(1)      NOT NULL DEFAULT 0                COMMENT '是否在菜单中隐藏',
  `active_menu`         VARCHAR(128)    NULL     DEFAULT NULL             COMMENT '高亮的菜单 routeName',
  `multi_tab`           TINYINT(1)      NOT NULL DEFAULT 0                COMMENT '是否支持多页签',
  `fixed_index_in_tab`  INT             NULL     DEFAULT NULL             COMMENT '在 tab 中的固定位置',
  `query`               JSON            NULL     DEFAULT NULL             COMMENT '路由 query 参数',
  `is_system`           TINYINT(1)      NOT NULL DEFAULT 0                COMMENT '是否系统内置(1=不可删/不可改 route_name)',
  `status`              CHAR(1)         NOT NULL DEFAULT '1'              COMMENT '状态: 1-启用, 2-禁用',
  `create_by`           BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '创建人 (存 sys_user.id)',
  `create_time`         DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `update_by`           BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '更新人 (存 sys_user.id)',
  `update_time`         DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at`          DATETIME(3)     NULL     DEFAULT NULL             COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_route_name` (`route_name`, ((IFNULL(`deleted_at`, '1970-01-01 00:00:00')))),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='菜单表';

-- ---------------------------------------------------------------------
-- Table: sys_menu_button
-- 加 deleted_at: 按钮删除会断 sys_role_button 引用, 软删后可恢复审计
-- ---------------------------------------------------------------------
DROP TABLE IF EXISTS `sys_menu_button`;
CREATE TABLE `sys_menu_button` (
  `id`          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键 id',
  `menu_id`     BIGINT UNSIGNED NOT NULL                COMMENT '菜单 id',
  `button_code` VARCHAR(64)     NOT NULL                COMMENT '按钮编码',
  `button_desc` VARCHAR(255)    NULL     DEFAULT NULL   COMMENT '按钮描述',
  `create_time` DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `update_time` DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at`  DATETIME(3)     NULL     DEFAULT NULL   COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_menu_button` (`menu_id`, `button_code`, ((IFNULL(`deleted_at`, '1970-01-01 00:00:00')))),
  KEY `idx_button_code` (`button_code`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='菜单按钮表';

-- ---------------------------------------------------------------------
-- Table: sys_role_menu
-- ---------------------------------------------------------------------
DROP TABLE IF EXISTS `sys_role_menu`;
CREATE TABLE `sys_role_menu` (
  `id`          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键 id',
  `role_id`     BIGINT UNSIGNED NOT NULL                COMMENT '角色 id',
  `menu_id`     BIGINT UNSIGNED NOT NULL                COMMENT '菜单 id',
  `create_time` DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_role_menu` (`role_id`, `menu_id`),
  KEY `idx_menu_id` (`menu_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色菜单关联表';

-- ---------------------------------------------------------------------
-- Table: sys_role_button
-- ---------------------------------------------------------------------
DROP TABLE IF EXISTS `sys_role_button`;
CREATE TABLE `sys_role_button` (
  `id`          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键 id',
  `role_id`     BIGINT UNSIGNED NOT NULL                COMMENT '角色 id',
  `button_id`   BIGINT UNSIGNED NOT NULL                COMMENT '按钮 id',
  `create_time` DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_role_button` (`role_id`, `button_id`),
  KEY `idx_button_id` (`button_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色按钮关联表';

-- ---------------------------------------------------------------------
-- Table: sys_login_log
-- 日志表: 追加型, 无需 deleted_at / update_time
-- login_result 区别于其他表 status (启用/禁用), 避免 GORM BaseModel 抽象冲突
-- user_agent 用 TEXT, 现代浏览器 UA 经常超过 512 字符
-- 索引: idx_user_time 用于"查某用户登录历史 ORDER BY DESC"高频场景
-- ---------------------------------------------------------------------
DROP TABLE IF EXISTS `sys_login_log`;
CREATE TABLE `sys_login_log` (
  `id`            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT          COMMENT '主键 id',
  `user_type`     CHAR(1)         NOT NULL DEFAULT '1'              COMMENT '用户类型: 1-后台 sys_user, 2-前台 reader',
  `user_id`       BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '用户 id',
  `user_name`     VARCHAR(64)     NOT NULL                          COMMENT '登录用户名',
  `login_ip`      VARCHAR(64)     NULL     DEFAULT NULL             COMMENT '登录 IP',
  `user_agent`    TEXT            NULL     DEFAULT NULL             COMMENT '用户代理',
  `login_type`    CHAR(1)         NOT NULL DEFAULT '1'              COMMENT '类型: 1-登录, 2-登出',
  `login_result`  CHAR(1)         NOT NULL                          COMMENT '结果: 1-成功, 2-失败',
  `message`       VARCHAR(255)    NULL     DEFAULT NULL             COMMENT '提示信息',
  `login_time`    DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '登录时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_time` (`user_type`, `user_id`, `login_time` DESC),
  KEY `idx_user_name` (`user_name`),
  KEY `idx_login_time` (`login_time` DESC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='登录日志表';

-- ---------------------------------------------------------------------
-- Table: sys_operation_log (操作日志)
-- 追加型, 通过 Gin 中间件统一拦截写入
-- 用于回答: "谁在何时改了什么, 改成了什么"
-- request_body 应在中间件层做密码/token 脱敏后再入库
-- ---------------------------------------------------------------------
DROP TABLE IF EXISTS `sys_operation_log`;
CREATE TABLE `sys_operation_log` (
  `id`             BIGINT UNSIGNED NOT NULL AUTO_INCREMENT          COMMENT '主键 id',
  `user_id`        BIGINT UNSIGNED NOT NULL                          COMMENT '操作人 sys_user.id',
  `user_name`      VARCHAR(64)     NOT NULL                          COMMENT '操作人名 (冗余, 防止改名后查不到)',
  `module`         VARCHAR(64)     NOT NULL                          COMMENT '业务模块: user/role/menu/dept/dict/...',
  `action`         VARCHAR(32)     NOT NULL                          COMMENT '动作: create/update/delete/grant/revoke/...',
  `target_id`     BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '操作对象 id',
  `target_name`   VARCHAR(128)    NULL     DEFAULT NULL             COMMENT '操作对象名 (冗余便于检索)',
  `request_url`   VARCHAR(255)    NULL     DEFAULT NULL             COMMENT '请求 URL',
  `request_method` VARCHAR(16)    NULL     DEFAULT NULL             COMMENT 'HTTP 方法',
  `request_body`   TEXT           NULL     DEFAULT NULL             COMMENT '请求体 JSON (脱敏)',
  `response_code`  INT            NULL     DEFAULT NULL             COMMENT '业务响应码',
  `client_ip`      VARCHAR(64)    NULL     DEFAULT NULL             COMMENT '客户端 IP',
  `user_agent`     TEXT           NULL     DEFAULT NULL             COMMENT '用户代理',
  `cost_ms`        INT UNSIGNED   NOT NULL DEFAULT 0                COMMENT '耗时毫秒',
  `operate_time`   DATETIME(3)    NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '操作时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_time` (`user_id`, `operate_time` DESC),
  KEY `idx_module_action` (`module`, `action`),
  KEY `idx_target` (`module`, `target_id`),
  KEY `idx_operate_time` (`operate_time` DESC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='操作日志表';

-- ---------------------------------------------------------------------
-- Table: sys_dict (字典分类)
-- 例: gender, enable_status, book_serial_status
-- is_system 防止系统内置字典被误删
-- ---------------------------------------------------------------------
DROP TABLE IF EXISTS `sys_dict`;
CREATE TABLE `sys_dict` (
  `id`          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT          COMMENT '主键 id',
  `dict_name`   VARCHAR(64)     NOT NULL                          COMMENT '字典名称',
  `dict_code`   VARCHAR(64)     NOT NULL                          COMMENT '字典编码 (业务唯一键)',
  `dict_desc`   VARCHAR(255)    NULL     DEFAULT NULL             COMMENT '字典描述',
  `is_system`   TINYINT(1)      NOT NULL DEFAULT 0                COMMENT '是否系统内置(1=不可删/不可改 dict_code)',
  `status`      CHAR(1)         NOT NULL DEFAULT '1'              COMMENT '状态: 1-启用, 2-禁用',
  `create_by`   BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '创建人 (存 sys_user.id)',
  `create_time` DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `update_by`   BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '更新人 (存 sys_user.id)',
  `update_time` DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at`  DATETIME(3)     NULL     DEFAULT NULL             COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_dict_code` (`dict_code`, ((IFNULL(`deleted_at`, '1970-01-01 00:00:00')))),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='字典分类表';

-- ---------------------------------------------------------------------
-- Table: sys_dict_item (字典项)
-- 仅 (dict_id, item_value) 业务唯一, label 由前端 UI 校验提示, 不做强约束
-- 理由: 运营场景下 label 偶尔重复 (近义词) 是真实需求, 强约束会卡录入
-- ---------------------------------------------------------------------
DROP TABLE IF EXISTS `sys_dict_item`;
CREATE TABLE `sys_dict_item` (
  `id`          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT          COMMENT '主键 id',
  `dict_id`     BIGINT UNSIGNED NOT NULL                          COMMENT '字典 id',
  `item_label`  VARCHAR(128)    NOT NULL                          COMMENT '显示标签 (UI 校验, 无强唯一)',
  `item_value`  VARCHAR(128)    NOT NULL                          COMMENT '实际值',
  `item_desc`   VARCHAR(255)    NULL     DEFAULT NULL             COMMENT '描述',
  `sort_order`  INT             NOT NULL DEFAULT 0                COMMENT '排序',
  `status`      CHAR(1)         NOT NULL DEFAULT '1'              COMMENT '状态: 1-启用, 2-禁用',
  `create_by`   BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '创建人 (存 sys_user.id)',
  `create_time` DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `update_by`   BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '更新人 (存 sys_user.id)',
  `update_time` DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at`  DATETIME(3)     NULL     DEFAULT NULL             COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_dict_value` (`dict_id`, `item_value`, ((IFNULL(`deleted_at`, '1970-01-01 00:00:00')))),
  KEY `idx_dict_id` (`dict_id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='字典项表';