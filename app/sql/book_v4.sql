-- =====================================================================
-- book_v4.sql —— 高级特性: 角色系统 + 阅读事件
-- Engine: InnoDB / utf8mb4_unicode_ci
-- Dependency: 依赖 book_v1.sql (book) 和 book_v2.sql (book_chapter)
-- Notes:
--   1. 角色支持主角/配角分类, 别名以 | 包裹存储, 查找时 LIKE '%|xxx|%'
--   2. book_character_chapter 已注释, 需用时取消注释即可
--   3. 阅读事件 reader_read_event 是唯一事实源, 周/月/年统计由服务层实时 SUM
--   4. 阅读统计不做预聚合, 直接 SQL 实时查询
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
  `alias`         VARCHAR(255)    NULL     DEFAULT NULL             COMMENT '别名 (花名/字/号, 格式: |别名1|别名2|..., LIKE 查询 %|目标|%)',
  `role_type`     CHAR(1)         NOT NULL DEFAULT '2'              COMMENT '类型: 1-主角, 2-配角, 3-反派, 4-龙套',
  `avatar`        VARCHAR(512)    NULL     DEFAULT NULL             COMMENT '角色头像 URL',
  `intro`         TEXT            NULL     DEFAULT NULL             COMMENT '角色简介',
  `extra`         JSON            NULL     DEFAULT NULL             COMMENT '扩展属性 (JSON, 如生卒年/门派/种族等)',
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
-- Table: book_character_rel (角色关系)
-- 两个角色之间可以有多重关系 (如 A 是 B 的师父, 也是 B 的父亲)
-- character_a -> character_b 表示有向关系, 非对称关系请确定方向存入
-- 对称关系 (如夫妻/兄弟) 建议按 (较小ID, 较大ID) 存储, 查询时 UNION 正反向
-- ---------------------------------------------------------------------
DROP TABLE IF EXISTS `book_character_rel`;
CREATE TABLE `book_character_rel` (
  `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT          COMMENT '主键 id',
  `book_id`         BIGINT UNSIGNED NOT NULL                          COMMENT '归属作品 id (book.id, 冗余便于按书查询)',
  `character_a_id`  BIGINT UNSIGNED NOT NULL                          COMMENT '角色 A id (book_character.id)',
  `character_b_id`  BIGINT UNSIGNED NOT NULL                          COMMENT '角色 B id (book_character.id)',
  `relation_type`   VARCHAR(32)     NOT NULL                          COMMENT '关系类型: 师徒/父子/夫妻/兄弟/仇敌/挚友...',
  `relation_desc`   VARCHAR(255)    NULL     DEFAULT NULL             COMMENT '关系描述 (可选, 如"授业恩师")',
  `sort_order`      INT             NOT NULL DEFAULT 0                COMMENT '排序',
  `create_by`       BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '创建人 (存 sys_user.id)',
  `create_time`     DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `update_by`       BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '更新人 (存 sys_user.id)',
  `update_time`     DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at`      DATETIME(3)     NULL     DEFAULT NULL             COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_rel_type` (`character_a_id`, `character_b_id`, `relation_type`),
  KEY `idx_book_id` (`book_id`),
  KEY `idx_character_b` (`character_b_id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色关系表';

-- ---------------------------------------------------------------------
-- [暂时注释] book_character_chapter (角色-章节出场)
-- 记录角色在哪些章节出场, 支持按角色筛选章节
-- ---------------------------------------------------------------------
-- CREATE TABLE `book_character_chapter` (
--   `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT          COMMENT '主键 id',
--   `character_id`    BIGINT UNSIGNED NOT NULL                          COMMENT '角色 id (book_character.id)',
--   `chapter_id`      BIGINT UNSIGNED NOT NULL                          COMMENT '章节 id (book_chapter.id)',
--   `appearance_desc` VARCHAR(255)    NULL     DEFAULT NULL             COMMENT '出场描述 (可选, 如"首次登场")',
--   `create_by`       BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '创建人 (存 sys_user.id)',
--   `create_time`     DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
--   `update_by`       BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '更新人 (存 sys_user.id)',
--   `update_time`     DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
--   `deleted_at`      DATETIME(3)     NULL     DEFAULT NULL             COMMENT '软删除时间',
--   PRIMARY KEY (`id`),
--   UNIQUE KEY `uk_char_chapter` (`character_id`, `chapter_id`),
--   KEY `idx_chapter_id` (`chapter_id`)
-- ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色-章节出场表';

-- ---------------------------------------------------------------------
-- Table: reader_read_event (阅读事件 / 原子日志)
-- 一次"上报阅读心跳" 一行: 客户端每 30-60s 心跳一次, 或翻章时上报
-- 设计要点:
--   - 纯追加写入, 无更新/删除, 所以不带 create_by/update_by/deleted_at
--   - session_id 标识单次阅读会话 (从打开书到离开/超时), 用于聚合"会话次数"
--   - word_count 由客户端上报本次心跳区间实际阅读字数 (非回填整章字数)
--   - 周/月/年/总计聚合由服务层实时 SUM, 不做预聚合中间表
-- ---------------------------------------------------------------------
DROP TABLE IF EXISTS `reader_read_event`;
CREATE TABLE `reader_read_event` (
  `id`             BIGINT UNSIGNED NOT NULL AUTO_INCREMENT          COMMENT '主键 id',
  `reader_id`      BIGINT UNSIGNED NOT NULL                          COMMENT '阅读者 id',
  `book_id`        BIGINT UNSIGNED NOT NULL                          COMMENT '书籍 id',
  `chapter_id`     BIGINT UNSIGNED NOT NULL                          COMMENT '章节 id',
  `session_id`     VARCHAR(36)     NOT NULL                          COMMENT '阅读会话 UUID (一次打开书的完整阅读)服务端生成',
  `duration_sec`   INT UNSIGNED    NOT NULL DEFAULT 0                COMMENT '本次心跳区间阅读时长 (秒)',
  `word_count`     INT UNSIGNED    NOT NULL DEFAULT 0                COMMENT '本次心跳区间滚动字数 (客户端上报, 非整章回填)',
  `event_date`     DATE            NOT NULL                          COMMENT '事件日期 (聚合分区键, 冗余于 event_time)',
  `event_time`     DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '事件时间',
  `device_type`    VARCHAR(32)     NOT NULL DEFAULT 'web'             COMMENT '设备类型: web/ios/android/pc',
  PRIMARY KEY (`id`),
  KEY `idx_reader_date` (`reader_id`, `event_date`),
  KEY `idx_book_date` (`book_id`, `event_date`),
  KEY `idx_session` (`session_id`),
  KEY `idx_event_date` (`event_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='阅读事件原子表(追加日志,无更新)';

-- ---------------------------------------------------------------------
-- [已删除] reader_read_stat_daily (按日聚合)
-- 原设计: 凌晨定时任务聚合 reader_read_event 到此表
-- 删除原因: 避免定时任务运维负担 + 预聚合表的重复维护
-- 替代方案: 由服务层 SELECT SUM/COUNT(DISTINCT) FROM reader_read_event WHERE event_date BETWEEN ? AND ? GROUP BY reader_id
-- ---------------------------------------------------------------------

-- ---------------------------------------------------------------------
-- [已删除] reader_read_stat_book (按读者-书-日聚合)
-- 替代方案: SELECT SUM(duration_sec), SUM(word_count), COUNT(DISTINCT chapter_id)
--           FROM reader_read_event WHERE reader_id=? AND book_id=? AND event_date BETWEEN ? AND ?
-- ---------------------------------------------------------------------
