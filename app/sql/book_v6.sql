
-- =====================================================================
-- 系统配置表
-- Engine: InnoDB / utf8mb4_unicode_ci
-- =====================================================================

DROP TABLE IF EXISTS `sys_setting`;
CREATE TABLE `sys_setting` (
    `id`            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `category`      VARCHAR(64)     NOT NULL COMMENT '配置分类: system/clean/dedup/recommend/sync等',
    `key`           VARCHAR(128)    NOT NULL COMMENT '配置键名',
    `value`         LONGTEXT        NOT NULL COMMENT '配置值(JSON格式存储)',
    `value_type`    VARCHAR(16)     NOT NULL DEFAULT 'string' COMMENT '值类型: string/number/boolean/json/array',
    `description`   VARCHAR(255)    NULL     DEFAULT NULL COMMENT '配置说明',
    `editable`      TINYINT(1)      NOT NULL DEFAULT 1 COMMENT '是否可编辑: 0-只读,1-可编辑',
    `is_system`     TINYINT(1)      NOT NULL DEFAULT 0 COMMENT '是否系统内置: 0-用户配置,1-系统内置(不可删除)',
    `status`        CHAR(1)         NOT NULL DEFAULT '1'              COMMENT '状态: 1-启用, 2-禁用',
    `create_by`     BIGINT UNSIGNED NULL     DEFAULT NULL COMMENT '创建人',
    `create_time`   DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    `update_by`     BIGINT UNSIGNED NULL     DEFAULT NULL COMMENT '更新人',
    `update_time`   DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
    `deleted_at`    DATETIME(3)     NULL     DEFAULT NULL COMMENT '软删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_category_key` (`category`, `key`, `deleted_at`),
    KEY `idx_category` (`category`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统配置表';