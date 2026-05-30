-- =====================================================================
-- book_v1.sql —— 书库核心: 分类 / 标签 / 作品主表
-- Engine: InnoDB / utf8mb4_unicode_ci
-- Dependency: 依赖 system-manage.sql 中的 sys_user 表
-- Notes:
--   1. sys_user_profile 扩展读者信息, 依赖 sys_user.id
--   2. book_category 自关联树, ancestors 加速子树查询
--   3. book 是聚合后的唯一书籍, title+author 作为聚合标识
--   4. book_tag_rel 连接 book 和 book_tag
-- =====================================================================

USE `boread`;

-- ---------------------------------------------------------------------
-- Table: sys_user_profile (用户扩展信息)
-- 依赖 sys_user.id, 无独立认证
-- ---------------------------------------------------------------------
DROP TABLE IF EXISTS `sys_user_profile`;
CREATE TABLE `sys_user_profile` (
  `user_id`       BIGINT UNSIGNED NOT NULL                          COMMENT '用户 id (FK -> sys_user.id)',
  `nickname`      VARCHAR(64)     NOT NULL DEFAULT ''               COMMENT '昵称 (覆盖 sys_user.nick_name)',
  `avatar`        VARCHAR(255)    NULL     DEFAULT NULL             COMMENT '头像 URL (覆盖 sys_user.avatar)',
  `signature`     VARCHAR(255)    NULL     DEFAULT NULL             COMMENT '个性签名',
  `create_time`   DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `update_time`   DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `create_by`     BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '创建人 (存 sys_user.id)',
  `update_by`     BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '更新人 (存 sys_user.id)',
  `deleted_at`    DATETIME(3)     NULL     DEFAULT NULL             COMMENT '软删除时间',
  PRIMARY KEY (`user_id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户扩展信息表';

-- ---------------------------------------------------------------------
-- Table: book_category (书籍分类)
-- 自关联树, 例: 文学 -> 小说 -> 玄幻
-- ---------------------------------------------------------------------
DROP TABLE IF EXISTS `book_category`;
CREATE TABLE `book_category` (
  `id`            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT          COMMENT '主键 id',
  `parent_id`     BIGINT UNSIGNED NOT NULL DEFAULT 0                COMMENT '父分类 id',
  `ancestors`     VARCHAR(1024)   NOT NULL DEFAULT ''                COMMENT '祖先路径 (逗号分隔, 如 0,1,5), 用于快速查子树',
  `category_name` VARCHAR(64)     NOT NULL                          COMMENT '分类名称',
  `category_code` VARCHAR(64)     NOT NULL                          COMMENT '分类编码',
  `description`   VARCHAR(255)    NULL     DEFAULT NULL             COMMENT '分类描述',
  `sort_order`    INT             NOT NULL DEFAULT 0                COMMENT '排序',
  `status`        CHAR(1)         NOT NULL DEFAULT '1'              COMMENT '状态: 1-启用, 2-禁用',
  `create_by`     BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '创建人 (存 sys_user.id)',
  `create_time`   DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `update_by`     BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '更新人 (存 sys_user.id)',
  `update_time`   DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at`    DATETIME(3)     NULL     DEFAULT NULL             COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_category_code` (`category_code`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='书籍分类表';

-- ---------------------------------------------------------------------
-- Table: book_tag (书籍标签)
-- ---------------------------------------------------------------------
DROP TABLE IF EXISTS `book_tag`;
CREATE TABLE `book_tag` (
  `id`          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键 id',
  `tag_name`    VARCHAR(64)     NOT NULL                COMMENT '标签名',
  `description` VARCHAR(255)    NULL     DEFAULT NULL   COMMENT '标签描述',
  `usage_count` INT UNSIGNED    NOT NULL DEFAULT 0      COMMENT '引用计数 (冗余, 便于排序热门)',
  `create_by`   BIGINT UNSIGNED NULL     DEFAULT NULL   COMMENT '创建人 (存 sys_user.id)',
  `create_time` DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `update_by`   BIGINT UNSIGNED NULL     DEFAULT NULL   COMMENT '更新人 (存 sys_user.id)',
  `update_time` DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at`  DATETIME(3)     NULL     DEFAULT NULL   COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_tag_name` (`tag_name`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='书籍标签表';

-- ---------------------------------------------------------------------
-- Table: book (作品 / 聚合后的唯一书籍)
-- 一个小说只有一条, title+author 作为业务聚合标识
-- 书架/阅读进度/笔记/评论等全部指向这张表
-- 文件级字段已下放到 book_file, 支持多文件聚合
-- ---------------------------------------------------------------------
DROP TABLE IF EXISTS `book`;
CREATE TABLE `book` (
  `id`                BIGINT UNSIGNED NOT NULL AUTO_INCREMENT          COMMENT '主键 id',
  `title`             VARCHAR(255)    NOT NULL                          COMMENT '书名',
  `author`            VARCHAR(128)    NOT NULL DEFAULT ''               COMMENT '作者',
  `cover`             VARCHAR(512)    NULL     DEFAULT NULL             COMMENT '封面 URL',
  `intro`             TEXT            NULL     DEFAULT NULL             COMMENT '简介',
  `category_id`       BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '分类 id (book_category.id)',
  `language`          VARCHAR(16)     NOT NULL DEFAULT 'zh-CN'          COMMENT '语言',
  `serial_status`     CHAR(1)         NOT NULL DEFAULT '1'              COMMENT '连载状态: 1-连载中, 2-已完结, 3-断更',
  `visibility`        CHAR(1)         NOT NULL DEFAULT '1'              COMMENT '可见性: 1-公开, 2-仅自己, 3-部门内',
  `primary_file_id`   BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '主版本文件 id (book_file.id), 阅读器默认读取此文件的字节索引',
  `total_chapters`    INT UNSIGNED    NOT NULL DEFAULT 0                COMMENT '聚合后总章节数 (冗余)',
  `total_words`       INT UNSIGNED    NOT NULL DEFAULT 0                COMMENT '聚合后总字数 (冗余)',
  `aggregate_status`  CHAR(1)         NOT NULL DEFAULT '1'              COMMENT '聚合状态: 1-单文件(无需聚合), 2-多文件聚合中, 3-聚合完成',
  `avg_rating`        DECIMAL(2,1)    NOT NULL DEFAULT 0.0              COMMENT '平均评分 (0.0-5.0, 冗余)',
  `rating_count`      INT UNSIGNED    NOT NULL DEFAULT 0                COMMENT '评分人数 (冗余)',
  `owner_id`          BIGINT UNSIGNED NOT NULL                          COMMENT '创建者 id (数据权限)',
  `dept_id`           BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '创建者所属部门 id',
  `status`            CHAR(1)         NOT NULL DEFAULT '1'              COMMENT '上架状态: 1-已上架, 2-下架, 3-审核中, 4-审核拒绝',
  `create_by`         BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '创建人 (存 sys_user.id)',
  `create_time`       DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `update_by`         BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '更新人 (存 sys_user.id)',
  `update_time`       DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at`        DATETIME(3)     NULL     DEFAULT NULL             COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_title` (`title`),
  KEY `idx_author` (`author`),
  KEY `idx_title_author` (`title`, `author`),
  KEY `idx_category` (`category_id`),
  KEY `idx_dept_id` (`dept_id`),
  KEY `idx_status_visibility` (`status`, `visibility`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='作品表 (聚合后的唯一书籍)';

-- ---------------------------------------------------------------------
-- Table: book_tag_rel (书籍-标签 关联)
-- ---------------------------------------------------------------------
DROP TABLE IF EXISTS `book_tag_rel`;
CREATE TABLE `book_tag_rel` (
  `id`      BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键 id',
  `book_id` BIGINT UNSIGNED NOT NULL                COMMENT '书籍 id',
  `tag_id`  BIGINT UNSIGNED NOT NULL                COMMENT '标签 id',
  `create_by`     BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '创建人 (存 sys_user.id)',
  `create_time`   DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `update_by`     BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '更新人 (存 sys_user.id)',
  `update_time`   DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at`    DATETIME(3)     NULL     DEFAULT NULL             COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_book_tag` (`book_id`, `tag_id`),
  KEY `idx_tag_id` (`tag_id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='书籍标签关联表';
