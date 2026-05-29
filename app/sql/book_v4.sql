-- =====================================================================
-- book_v4.sql —— 高级特性: 角色系统 + 阅读统计
-- Engine: InnoDB / utf8mb4_unicode_ci
-- Dependency: 依赖 book_v1.sql (book) 和 book_v2.sql (book_chapter)
-- Notes:
--   1. 角色系统支持主角/配角分类, 别名逗号分隔
--   2. book_character_chapter 记录角色出场章节, 支持按角色筛选章节
--   3. 阅读统计分三层: 原子事件 -> 日聚合 -> 书-日聚合
--   4. 周/月/年/总计 = SUM(GROUP BY) reader_read_stat_daily, 无需额外建表
-- =====================================================================

USE `boread`;

-- ---------------------------------------------------------------------
-- Table: book_character (小说角色)
-- 一个作品有多个角色 (主角/配角), 支持按角色搜索和筛选
-- ---------------------------------------------------------------------
DROP TABLE IF EXISTS `book_character`;
CREATE TABLE `book_character` (
  `id`            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT          COMMENT '主键 id',
  `book_id`       BIGINT UNSIGNED NOT NULL                          COMMENT '归属作品 id (book.id)',
  `name`          VARCHAR(128)    NOT NULL                          COMMENT '角色名',
  `alias`         VARCHAR(255)    NULL     DEFAULT NULL             COMMENT '别名 (花名/字/号, 逗号分隔)',
  `role_type`     CHAR(1)         NOT NULL DEFAULT '2'              COMMENT '类型: 1-主角, 2-配角, 3-反派, 4-龙套',
  `avatar`        VARCHAR(512)    NULL     DEFAULT NULL             COMMENT '角色头像 URL',
  `intro`         TEXT            NULL     DEFAULT NULL             COMMENT '角色简介',
  `sort_order`    INT             NOT NULL DEFAULT 0                COMMENT '排序 (主角靠前)',
  `create_by`     BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '创建人 (存 sys_user.id)',
  `create_time`   DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `update_by`     BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '更新人 (存 sys_user.id)',
  `update_time`   DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at`    DATETIME(3)     NULL     DEFAULT NULL             COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_book_id` (`book_id`),
  KEY `idx_name` (`name`),
  KEY `idx_role_type` (`role_type`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='小说角色表';

-- ---------------------------------------------------------------------
-- Table: book_character_chapter (角色-章节出场)
-- 记录角色在哪些章节出场, 支持按角色筛选章节
-- ---------------------------------------------------------------------
DROP TABLE IF EXISTS `book_character_chapter`;
CREATE TABLE `book_character_chapter` (
  `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT          COMMENT '主键 id',
  `character_id`    BIGINT UNSIGNED NOT NULL                          COMMENT '角色 id (book_character.id)',
  `chapter_id`      BIGINT UNSIGNED NOT NULL                          COMMENT '章节 id (book_chapter.id)',
  `appearance_desc` VARCHAR(255)    NULL     DEFAULT NULL             COMMENT '出场描述 (可选, 如"首次登场")',
  `create_by`       BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '创建人 (存 sys_user.id)',
  `create_time`     DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `update_by`       BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '更新人 (存 sys_user.id)',
  `update_time`     DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at`      DATETIME(3)     NULL     DEFAULT NULL             COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_char_chapter` (`character_id`, `chapter_id`),
  KEY `idx_chapter_id` (`chapter_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色-章节出场表';

-- ---------------------------------------------------------------------
-- Table: reader_read_event (阅读事件 / 原子日志)
-- 一次"上报阅读心跳" 一行: 客户端每 30-60s 心跳一次, 或翻章时上报
-- 设计要点:
--   不直接拿来跑周/月/年统计 (扫描过多)
--   只用作"明细溯源" + 当日聚合的数据源
--   按 event_date 分区可显著加速查询 (后期数据量大再分区)
-- ---------------------------------------------------------------------
DROP TABLE IF EXISTS `reader_read_event`;
CREATE TABLE `reader_read_event` (
  `id`             BIGINT UNSIGNED NOT NULL AUTO_INCREMENT          COMMENT '主键 id',
  `reader_id`      BIGINT UNSIGNED NOT NULL                          COMMENT '阅读者 id',
  `book_id`        BIGINT UNSIGNED NOT NULL                          COMMENT '书籍 id',
  `chapter_id`     BIGINT UNSIGNED NOT NULL                          COMMENT '章节 id',
  `duration_sec`   INT UNSIGNED    NOT NULL DEFAULT 0                COMMENT '本次阅读时长 (秒)',
  `word_count`     INT UNSIGNED    NOT NULL DEFAULT 0                COMMENT '本次阅读字数 (服务端从 book_chapter.word_count 回填, 客户端上报仅供参考)',
  `event_date`     DATE            NOT NULL                          COMMENT '事件日期 (聚合维度, 冗余于 event_time)',
  `event_time`     DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '事件时间',
  `create_by`      BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '创建人 (存 sys_user.id)',
  `create_time`    DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `update_by`      BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '更新人 (存 sys_user.id)',
  `update_time`    DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at`     DATETIME(3)     NULL     DEFAULT NULL             COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_reader_date` (`reader_id`, `event_date`),
  KEY `idx_book_id` (`book_id`),
  KEY `idx_event_date` (`event_date`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='阅读事件原子表';

-- ---------------------------------------------------------------------
-- Table: reader_read_stat_daily (按日聚合)
-- 凌晨定时任务把昨日 reader_read_event 聚合到这张表, 用 INSERT ... ON DUPLICATE KEY UPDATE 累加
-- 周/月/年/总 = SUM(GROUP BY) 此表, 无需额外建表
-- ---------------------------------------------------------------------
DROP TABLE IF EXISTS `reader_read_stat_daily`;
CREATE TABLE `reader_read_stat_daily` (
  `id`             BIGINT UNSIGNED NOT NULL AUTO_INCREMENT          COMMENT '主键 id',
  `reader_id`      BIGINT UNSIGNED NOT NULL                          COMMENT '阅读者 id',
  `stat_date`      DATE            NOT NULL                          COMMENT '统计日期',
  `read_duration`  INT UNSIGNED    NOT NULL DEFAULT 0                COMMENT '总阅读时长 (秒)',
  `read_words`     INT UNSIGNED    NOT NULL DEFAULT 0                COMMENT '总阅读字数',
  `book_count`     INT UNSIGNED    NOT NULL DEFAULT 0                COMMENT '阅读书数 (去重)',
  `chapter_count`  INT UNSIGNED    NOT NULL DEFAULT 0                COMMENT '阅读章数 (去重)',
  `session_count`  INT UNSIGNED    NOT NULL DEFAULT 0                COMMENT '阅读会话次数 (心跳次数)',
  `create_by`      BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '创建人 (存 sys_user.id)',
  `create_time`    DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `update_by`      BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '更新人 (存 sys_user.id)',
  `update_time`    DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at`     DATETIME(3)     NULL     DEFAULT NULL             COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_reader_date` (`reader_id`, `stat_date`),
  KEY `idx_stat_date` (`stat_date`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='读者阅读日统计';

-- ---------------------------------------------------------------------
-- Table: reader_read_stat_book (按读者-书-日聚合)
-- 与 reader_read_stat_daily 解耦: 后者按"读者全部书"汇总, 这张按"读者+单本书"明细
-- 个人页"我读过的书排行"用这张, 全局活跃统计用 daily
-- ---------------------------------------------------------------------
DROP TABLE IF EXISTS `reader_read_stat_book`;
CREATE TABLE `reader_read_stat_book` (
  `id`             BIGINT UNSIGNED NOT NULL AUTO_INCREMENT          COMMENT '主键 id',
  `reader_id`      BIGINT UNSIGNED NOT NULL                          COMMENT '阅读者 id',
  `book_id`        BIGINT UNSIGNED NOT NULL                          COMMENT '书籍 id',
  `stat_date`      DATE            NOT NULL                          COMMENT '统计日期',
  `read_duration`  INT UNSIGNED    NOT NULL DEFAULT 0                COMMENT '当日该书阅读时长 (秒)',
  `read_words`     INT UNSIGNED    NOT NULL DEFAULT 0                COMMENT '当日该书阅读字数',
  `chapter_count`  INT UNSIGNED    NOT NULL DEFAULT 0                COMMENT '当日该书阅读章数 (去重)',
  `create_by`      BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '创建人 (存 sys_user.id)',
  `create_time`    DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `update_by`      BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '更新人 (存 sys_user.id)',
  `update_time`    DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at`     DATETIME(3)     NULL     DEFAULT NULL             COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_reader_book_date` (`reader_id`, `book_id`, `stat_date`),
  KEY `idx_book_date` (`book_id`, `stat_date`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='读者-书-日 阅读统计';
