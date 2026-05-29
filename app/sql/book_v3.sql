-- =====================================================================
-- book_v3.sql —— 读者交互: 书架 / 阅读进度 / 笔记 / 书评 / 评论
-- Engine: InnoDB / utf8mb4_unicode_ci
-- Dependency: 依赖 book_v1.sql (book) 和 book_v2.sql (book_chapter)
-- Notes:
--   1. 书架支持分组和置顶, last_read_time 用于排序
--   2. 阅读进度持续覆写, 记录章节内字符偏移和全书百分比
--   3. 笔记/划线合并为 reader_note, note_type 区分类型
--   4. 书评独立于章节评论, 支持评分 1-5 星
--   5. 章节评论楼中楼, parent_id 自关联
-- =====================================================================

USE `boread`;

-- ---------------------------------------------------------------------
-- Table: reader_bookshelf (书架 / 个人收藏)
-- 一个读者对一本书最多一条
-- ---------------------------------------------------------------------
DROP TABLE IF EXISTS `reader_bookshelf`;
CREATE TABLE `reader_bookshelf` (
  `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT          COMMENT '主键 id',
  `reader_id`       BIGINT UNSIGNED NOT NULL                          COMMENT '阅读者 id',
  `book_id`         BIGINT UNSIGNED NOT NULL                          COMMENT '书籍 id',
  `group_name`      VARCHAR(64)     NOT NULL DEFAULT '默认'           COMMENT '分组名 (读者自定义)',
  `is_top`          TINYINT(1)      NOT NULL DEFAULT 0                COMMENT '是否置顶',
  `last_read_time`  DATETIME(3)     NULL     DEFAULT NULL             COMMENT '最后阅读时间',
  `add_time`        DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '加入书架时间',
  `create_by`       BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '创建人 (存 sys_user.id)',
  `create_time`     DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `update_by`       BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '更新人 (存 sys_user.id)',
  `update_time`     DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at`      DATETIME(3)     NULL     DEFAULT NULL             COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_reader_book` (`reader_id`, `book_id`),
  KEY `idx_book_id` (`book_id`),
  KEY `idx_last_read` (`reader_id`, `last_read_time`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='读者书架表';

-- ---------------------------------------------------------------------
-- Table: reader_read_progress (阅读进度)
-- 一个读者对一本书一条进度, 持续覆写
-- ---------------------------------------------------------------------
DROP TABLE IF EXISTS `reader_read_progress`;
CREATE TABLE `reader_read_progress` (
  `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT          COMMENT '主键 id',
  `reader_id`       BIGINT UNSIGNED NOT NULL                          COMMENT '阅读者 id',
  `book_id`         BIGINT UNSIGNED NOT NULL                          COMMENT '书籍 id',
  `file_id`         BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '当前使用文件 id (book_file.id), 文件切换时重新映射',
  `chapter_id`      BIGINT UNSIGNED NOT NULL                          COMMENT '当前章节 id',
  `chapter_no`      INT UNSIGNED    NOT NULL                          COMMENT '当前章节序号 (冗余, 便于排序展示)',
  `position`        INT UNSIGNED    NOT NULL DEFAULT 0                COMMENT '章内位置 (字符偏移量)',
  `percent`         DECIMAL(5,2)    NOT NULL DEFAULT 0.00             COMMENT '全书进度百分比 0.00-100.00',
  `read_duration`   INT UNSIGNED    NOT NULL DEFAULT 0                COMMENT '累计阅读时长 (秒)',
  `last_read_time`  DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '最后阅读时间',
  `create_time`     DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `create_by`       BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '创建人 (存 sys_user.id)',
  `update_by`       BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '更新人 (存 sys_user.id)',
  `update_time`     DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at`      DATETIME(3)     NULL     DEFAULT NULL             COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_reader_book` (`reader_id`, `book_id`),
  KEY `idx_chapter` (`chapter_id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='阅读进度表';

-- ---------------------------------------------------------------------
-- Table: reader_note (笔记 / 划线)
-- 合并设计: note_type 区分纯笔记和划线 (划线 = 选段 + 可选内容)
-- ---------------------------------------------------------------------
DROP TABLE IF EXISTS `reader_note`;
CREATE TABLE `reader_note` (
  `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT          COMMENT '主键 id',
  `reader_id`       BIGINT UNSIGNED NOT NULL                          COMMENT '阅读者 id',
  `book_id`         BIGINT UNSIGNED NOT NULL                          COMMENT '书籍 id',
  `chapter_id`      BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '章节 id (允许整书笔记 NULL)',
  `note_type`       CHAR(1)         NOT NULL DEFAULT '1'              COMMENT '类型: 1-笔记(无选段), 2-划线(有选段), 3-划线+笔记',
  `selected_text`   TEXT            NULL     DEFAULT NULL             COMMENT '选中文本 (划线时)',
  `start_offset`    INT UNSIGNED    NULL     DEFAULT NULL             COMMENT '选段起始偏移',
  `end_offset`      INT UNSIGNED    NULL     DEFAULT NULL             COMMENT '选段结束偏移',
  `highlight_color` VARCHAR(16)     NULL     DEFAULT NULL             COMMENT '高亮颜色 (划线时, 如 yellow/red)',
  `content`         TEXT            NULL     DEFAULT NULL             COMMENT '笔记内容',
  `visibility`      CHAR(1)         NOT NULL DEFAULT '2'              COMMENT '可见性: 1-公开, 2-仅自己',
  `create_by`       BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '创建人 (存 sys_user.id)',
  `create_time`     DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `update_by`       BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '更新人 (存 sys_user.id)',
  `update_time`     DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at`      DATETIME(3)     NULL     DEFAULT NULL             COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_reader_book` (`reader_id`, `book_id`),
  KEY `idx_chapter` (`chapter_id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='读者笔记/划线表';

-- ---------------------------------------------------------------------
-- Table: book_review (整本书评)
-- 一个读者对一本书可发多条书评
-- ---------------------------------------------------------------------
DROP TABLE IF EXISTS `book_review`;
CREATE TABLE `book_review` (
  `id`            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT          COMMENT '主键 id',
  `book_id`       BIGINT UNSIGNED NOT NULL                          COMMENT '书籍 id',
  `reader_id`     BIGINT UNSIGNED NOT NULL                          COMMENT '阅读者 id',
  `rating`        TINYINT UNSIGNED NULL    DEFAULT NULL             COMMENT '评分 1-5 星',
  `title`         VARCHAR(255)    NULL     DEFAULT NULL             COMMENT '书评标题',
  `content`       TEXT            NOT NULL                          COMMENT '书评内容',
  `like_count`    INT UNSIGNED    NOT NULL DEFAULT 0                COMMENT '点赞数',
  `reply_count`   INT UNSIGNED    NOT NULL DEFAULT 0                COMMENT '回复数 (冗余)',
  `owner_id`      BIGINT UNSIGNED NOT NULL                          COMMENT '创建者 id (数据权限)',
  `dept_id`       BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '创建者所属部门 id',
  `status`        CHAR(1)         NOT NULL DEFAULT '1'              COMMENT '状态: 1-正常, 2-隐藏, 3-审核中',
  `create_by`     BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '创建人 (存 sys_user.id)',
  `create_time`   DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `update_by`     BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '更新人 (存 sys_user.id)',
  `update_time`   DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at`    DATETIME(3)     NULL     DEFAULT NULL             COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_book_id` (`book_id`),
  KEY `idx_reader_id` (`reader_id`),
  KEY `idx_owner` (`owner_id`),
  KEY `idx_dept_id` (`dept_id`),
  KEY `idx_deleted_at` (`deleted_at`),
  CONSTRAINT `ck_review_rating` CHECK (`rating` IS NULL OR (`rating` BETWEEN 1 AND 5))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='书评表';

-- ---------------------------------------------------------------------
-- Table: chapter_comment (章节评论)
-- 楼中楼: parent_id 自关联
-- ---------------------------------------------------------------------
DROP TABLE IF EXISTS `chapter_comment`;
CREATE TABLE `chapter_comment` (
  `id`            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT          COMMENT '主键 id',
  `book_id`       BIGINT UNSIGNED NOT NULL                          COMMENT '书籍 id (冗余, 便于按书查所有评论)',
  `chapter_id`    BIGINT UNSIGNED NOT NULL                          COMMENT '章节 id',
  `reader_id`     BIGINT UNSIGNED NOT NULL                          COMMENT '阅读者 id',
  `parent_id`     BIGINT UNSIGNED NOT NULL DEFAULT 0                COMMENT '父评论 id (0=顶层)',
  `reply_to_id`   BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '回复的读者 id (@某人)',
  `content`       TEXT            NOT NULL                          COMMENT '评论内容',
  `like_count`    INT UNSIGNED    NOT NULL DEFAULT 0                COMMENT '点赞数',
  `owner_id`      BIGINT UNSIGNED NOT NULL                          COMMENT '创建者 id (数据权限)',
  `dept_id`       BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '创建者所属部门 id',
  `status`        CHAR(1)         NOT NULL DEFAULT '1'              COMMENT '状态: 1-正常, 2-隐藏, 3-审核中',
  `create_by`     BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '创建人 (存 sys_user.id)',
  `create_time`   DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `update_by`     BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '更新人 (存 sys_user.id)',
  `update_time`   DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at`    DATETIME(3)     NULL     DEFAULT NULL             COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_chapter` (`chapter_id`),
  KEY `idx_book_id` (`book_id`),
  KEY `idx_reader_id` (`reader_id`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_owner` (`owner_id`),
  KEY `idx_dept_id` (`dept_id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='章节评论表';
