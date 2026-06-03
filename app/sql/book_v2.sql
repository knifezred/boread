-- =====================================================================
-- book_v2.sql —— 文件解析与章节管理
-- Engine: InnoDB / utf8mb4_unicode_ci
-- Dependency: 依赖 book_v1.sql 中的 book 表
-- Notes:
--   1. 章节内容存本地文件, DB 只存 byte_offset + byte_length 索引
--   2. book_file 支持多文件聚合, is_primary 标记默认版本
--   3. book_chapter_rule 定义章节切分规则, 支持全局默认和单书覆盖
--   4. book_content_filter_rule 定义入库/出库时的内容净化规则
-- =====================================================================

USE `boread`;

-- ---------------------------------------------------------------------
-- Table: book_file (物理文件)
-- 同一本书可能有多个 txt 副本, 解析后的字节偏移索引
-- pread 读取依赖此表的 content_path
-- ---------------------------------------------------------------------
DROP TABLE IF EXISTS `book_file`;
CREATE TABLE `book_file` (
  `id`                BIGINT UNSIGNED NOT NULL AUTO_INCREMENT          COMMENT '主键 id',
  `book_id`           BIGINT UNSIGNED NOT NULL                          COMMENT '归属作品 id (book.id)',
  `original_name`     VARCHAR(255)    NOT NULL                          COMMENT '原始文件名 (如 凡人修仙传_1-50.txt)',
  `source_type`       CHAR(1)         NOT NULL DEFAULT '1'              COMMENT '来源: 1-用户上传, 2-管理员上传, 3-本地扫描',
  `source_format`     VARCHAR(16)     NULL     DEFAULT NULL             COMMENT '源格式: txt, epub, mobi, pdf...',
  `source_file_url`   VARCHAR(512)    NULL     DEFAULT NULL             COMMENT '源文件 URL (保留原始上传件)',
  `content_path`      VARCHAR(512)    NULL     DEFAULT NULL             COMMENT '解析后的正文文件相对路径 (storage/books 下)',
  `content_size`      BIGINT UNSIGNED NOT NULL DEFAULT 0                COMMENT '正文文件大小 (bytes)',
  `content_md5`       CHAR(32)        NULL     DEFAULT NULL             COMMENT '正文文件 MD5 (秒传判重 + 完整性校验)',
  `content_charset`   VARCHAR(16)     NOT NULL DEFAULT 'utf-8'          COMMENT '正文文件字符编码',
  `content_version`   INT UNSIGNED    NOT NULL DEFAULT 1                COMMENT '文件内容版本号 (重新编码/清洗时递增, 配套 book_chapter 索引失效检测)',
  `chapter_count`     INT UNSIGNED    NOT NULL DEFAULT 0                COMMENT '该文件解析出的章节数 (冗余)',
  `is_primary`        TINYINT(1)      NOT NULL DEFAULT 0                COMMENT '是否为主版本 (阅读器默认使用, 同一 book 最多一个)',
  `file_status`       CHAR(1)         NOT NULL DEFAULT '1'              COMMENT '解析状态: 1-待处理, 2-处理中, 3-解析成功, 4-解析失败',
  `parse_message`     VARCHAR(512)    NULL     DEFAULT NULL             COMMENT '解析结果/失败原因',
  `create_by`         BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '创建人 (存 sys_user.id)',
  `create_time`       DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `update_by`         BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '更新人 (存 sys_user.id)',
  `update_time`       DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at`        DATETIME(3)     NULL     DEFAULT NULL             COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_content_md5` (`content_md5`),
  KEY `idx_book_id` (`book_id`),
  KEY `idx_book_primary` (`book_id`, `is_primary`),
  KEY `idx_file_status` (`file_status`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='书籍物理文件表';

-- ---------------------------------------------------------------------
-- Table: book_chapter (聚合章节索引)
-- 内容存文件 (book_file.content_path), pread 切片读
-- 章节按 (book_id, chapter_no) 聚合, 重复覆盖
-- 每个章节记录来自哪个文件, 便于溯源
-- ---------------------------------------------------------------------
DROP TABLE IF EXISTS `book_chapter`;
CREATE TABLE `book_chapter` (
  `id`            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT          COMMENT '主键 id',
  `book_id`       BIGINT UNSIGNED NOT NULL                          COMMENT '归属作品 id (book.id)',
  `file_id`       BIGINT UNSIGNED NOT NULL                          COMMENT '来源文件 id (book_file.id)',
  `volume_no`     INT UNSIGNED    NULL     DEFAULT NULL             COMMENT '卷序号 (1 开始)',
  `volume_title`  VARCHAR(255)    NULL     DEFAULT NULL             COMMENT '卷标题',
  `chapter_no`    INT UNSIGNED    NOT NULL                          COMMENT '章节序号 (1 开始)',
  `title`         VARCHAR(255)    NOT NULL                          COMMENT '章节标题',
  `byte_offset`   BIGINT UNSIGNED NOT NULL                          COMMENT '在 book_file.content_path 文件中的起始字节偏移',
  `byte_length`   INT UNSIGNED    NOT NULL                          COMMENT '章节字节长度 (UTF-8 编码后)',
  `word_count`    INT UNSIGNED    NOT NULL DEFAULT 0                COMMENT '字符数 (展示用,需过滤空格等无效字符)',
  `is_vip`        TINYINT(1)      NOT NULL DEFAULT 0                COMMENT '是否 VIP 章节',
  `status`        CHAR(1)         NOT NULL DEFAULT '1'              COMMENT '状态: 1-发布, 2-草稿, 3-下架',
  `create_by`     BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '创建人 (存 sys_user.id)',
  `create_time`   DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `update_by`     BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '更新人 (存 sys_user.id)',
  `update_time`   DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at`    DATETIME(3)     NULL     DEFAULT NULL             COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_book_file_chapter` (`book_id`, `file_id`, `chapter_no`),
  KEY `idx_book_id` (`book_id`),
  KEY `idx_file_id` (`file_id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='章节索引表（按文件存储，支持多版本切换）';

-- ---------------------------------------------------------------------
-- Table: book_upload (上传/解析任务)
-- 异步解析 epub/txt 的诊断记录, 失败原因可追溯
-- ---------------------------------------------------------------------
DROP TABLE IF EXISTS `book_upload`;
CREATE TABLE `book_upload` (
  `id`                BIGINT UNSIGNED NOT NULL AUTO_INCREMENT          COMMENT '主键 id',
  `book_id`           BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '关联作品 id (解析成功后回填)',
  `original_name`     VARCHAR(255)    NOT NULL                          COMMENT '原始文件名',
  `file_path`         VARCHAR(512)    NOT NULL                          COMMENT '文件存储路径',
  `file_size`         BIGINT UNSIGNED NOT NULL DEFAULT 0                COMMENT '文件大小 (bytes)',
  `file_md5`          CHAR(32)        NULL     DEFAULT NULL             COMMENT '文件 MD5',
  `source_format`     VARCHAR(16)     NULL     DEFAULT NULL             COMMENT '源格式: txt, epub',
  `parse_status`      CHAR(1)         NOT NULL DEFAULT '1'              COMMENT '解析状态: 1-待解析, 2-解析中, 3-解析成功, 4-解析失败',
  `parse_message`     VARCHAR(512)    NULL     DEFAULT NULL             COMMENT '解析结果/失败原因',
  `chapter_count`     INT UNSIGNED    NULL     DEFAULT NULL             COMMENT '解析出的章节数 (成功时)',
  `create_by`         BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '创建人 (存 sys_user.id)',
  `create_time`       DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `update_by`         BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '更新人 (存 sys_user.id)',
  `update_time`       DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at`        DATETIME(3)     NULL     DEFAULT NULL             COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_book_id` (`book_id`),
  KEY `idx_file_md5` (`file_md5`),
  KEY `idx_parse_status` (`parse_status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='书籍上传/解析任务表';

-- ---------------------------------------------------------------------
-- Table: book_chapter_rule (章节识别规则)
-- 用于解析 txt 时切分章节, 支持多条规则按优先级匹配
-- 存储系统默认规则和用户自定义规则
-- ---------------------------------------------------------------------
DROP TABLE IF EXISTS `book_chapter_rule`;
CREATE TABLE `book_chapter_rule` (
  `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `rule_name`       VARCHAR(64)     NOT NULL                       COMMENT '规则名称，如：标准章节、数字序号等',
  `rule_type`       CHAR(1)         NOT NULL DEFAULT '1'           COMMENT '规则类型: 1-系统默认, 2-用户自定义',
  `user_id`         BIGINT UNSIGNED NULL                           COMMENT '用户ID，rule_type=2时必填，标识该规则属于哪个用户',
  `title_pattern`   VARCHAR(512)    NOT NULL                       COMMENT '章节标题正则表达式，使用Go RE2语法，如：第(\\d+)章\\s+(.+)',
  `group_pattern`   VARCHAR(512)    NULL                           COMMENT '分卷标题正则表达式，使用Go RE2语法，如：第(\\d+)卷\\s+(.+)，空表示不分卷',
  `min_chapter_len` INT UNSIGNED    NOT NULL DEFAULT 100           COMMENT '章节最小字符数，低于此阈值则过滤（用于排除目录页、导航行等误匹配）',
  `max_chapter_len` INT UNSIGNED    NOT NULL DEFAULT 100000        COMMENT '章节最大字符数，超过此阈值则过滤（防止未正确切分的大块内容）',
  `sort_order`      INT             NOT NULL DEFAULT 0             COMMENT '全局排序字段，数字越小优先级越高，系统匹配时按此顺序依次尝试规则',
  `description`     VARCHAR(255)    NULL                           COMMENT '规则说明/示例，如：匹配"第1章 标题"格式的章节',
  `status`          CHAR(1)         NOT NULL DEFAULT '1'           COMMENT '状态: 1-启用, 2-禁用',
  `create_by`       BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '创建人 (存 sys_user.id)',
  `create_time`     DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `update_by`       BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '更新人 (存 sys_user.id)',
  `update_time`     DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at`      DATETIME(3)     NULL     DEFAULT NULL             COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_type_user` (`rule_type`, `user_id`)                     COMMENT '按规则类型和用户查询',
  KEY `idx_sort_status` (`sort_order`, `status`)                   COMMENT '按排序和状态查询活跃规则',
  KEY `idx_deleted_at` (`deleted_at`)                              COMMENT '软删除查询索引'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='章节规则定义表';


-- =====================================================
-- Table: book_chapter_rule_rel (书籍章节规则关联表)
-- 定义书籍与章节规则的绑定关系，支持多规则组合匹配和单规则直接使用
-- =====================================================
CREATE TABLE `book_chapter_rule_rel` (
  `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `book_id`         BIGINT UNSIGNED NOT NULL                       COMMENT '书籍ID，关联书籍主表',
  `reader_id`       BIGINT UNSIGNED NOT NULL                       COMMENT '读者ID，关联sys_user表',
  `rule_id`         BIGINT UNSIGNED NOT NULL                       COMMENT '规则ID，关联chapter_rule表',
  `create_by`       BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '创建人 (存 sys_user.id)',
  `create_time`     DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `update_by`       BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '更新人 (存 sys_user.id)',
  `update_time`     DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at`      DATETIME(3)     NULL     DEFAULT NULL             COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_book_rule` (`book_id`, `reader_id`)                 COMMENT '同一本书同一读者不能重复关联规则',
  KEY `idx_book_reader_id` (`book_id`, `reader_id`)              COMMENT '按书籍和读者查询关联',
  CONSTRAINT `fk_book_rule_rel_rule` FOREIGN KEY (`rule_id`) REFERENCES `chapter_rule` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='书籍规则关联表';


-- DROP TABLE IF EXISTS `book_chapter_rule`;
-- CREATE TABLE `book_chapter_rule` (
--   `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT          COMMENT '主键 id',
--   `rule_name`       VARCHAR(64)     NOT NULL                          COMMENT '规则名称',
--   `scope_type`      CHAR(1)         NOT NULL DEFAULT '1'              COMMENT '作用域: 1-全局默认, 2-单书覆盖',
--   `book_id`         BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '关联书 id (scope_type=2 时必填)',
--   `pattern`         VARCHAR(512)    NOT NULL                          COMMENT '章节标题正则 (Go RE2 语法)',
--   `title_group`     INT             NOT NULL DEFAULT 0                COMMENT '正则中作为标题的捕获组序号, 0=整行',
--   `min_chapter_len` INT UNSIGNED    NOT NULL DEFAULT 100              COMMENT '章节最小字符数 (过滤误匹配如目录页)',
--   `max_chapter_len` INT UNSIGNED    NOT NULL DEFAULT 100000           COMMENT '章节最大字符数 (过大可能是未切分)',
--   `priority`        INT             NOT NULL DEFAULT 0                COMMENT '优先级 (越大越先匹配)',
--   `description`     VARCHAR(255)    NULL     DEFAULT NULL             COMMENT '说明 / 示例',
--   `status`          CHAR(1)         NOT NULL DEFAULT '1'              COMMENT '状态: 1-启用, 2-禁用',
--   `create_by`       BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '创建人 (存 sys_user.id)',
--   `create_time`     DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
--   `update_by`       BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '更新人 (存 sys_user.id)',
--   `update_time`     DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
--   `deleted_at`      DATETIME(3)     NULL     DEFAULT NULL             COMMENT '软删除时间',
--   PRIMARY KEY (`id`),
--   KEY `idx_scope` (`scope_type`, `book_id`),
--   KEY `idx_priority` (`priority`),
--   KEY `idx_deleted_at` (`deleted_at`)
-- ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='章节识别规则表';

-- ---------------------------------------------------------------------
-- Table: book_content_filter_rule (内容净化规则)
-- 入库时(解析阶段) 或 出库时(读章节时) 应用
-- action: 1-替换  2-拦截整章  3-标记审核
-- ---------------------------------------------------------------------
DROP TABLE IF EXISTS `book_content_filter_rule`;
CREATE TABLE `book_content_filter_rule` (
  `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT          COMMENT '主键 id',
  `rule_name`       VARCHAR(64)     NOT NULL                          COMMENT '规则名称',
  `match_type`      CHAR(1)         NOT NULL DEFAULT '1'              COMMENT '匹配方式: 1-关键词, 2-正则',
  `pattern`         VARCHAR(512)    NOT NULL                          COMMENT '关键词或正则',
  `action`          CHAR(1)         NOT NULL DEFAULT '1'              COMMENT '动作: 1-替换, 2-拦截整章, 3-标记审核',
  `replacement`     VARCHAR(255)    NOT NULL DEFAULT '***'            COMMENT '替换文本 (action=1 时使用)',
  `apply_stage`     CHAR(1)         NOT NULL DEFAULT '1'              COMMENT '应用阶段: 1-入库时(解析阶段一次性), 2-出库时(读章节实时)',
  `category`        VARCHAR(32)     NULL     DEFAULT NULL             COMMENT '分类标签: politics/porn/violence/ad...',
  `severity`        CHAR(1)         NOT NULL DEFAULT '1'              COMMENT '严重程度: 1-低, 2-中, 3-高',
  `description`     VARCHAR(255)    NULL     DEFAULT NULL             COMMENT '说明',
  `status`          CHAR(1)         NOT NULL DEFAULT '1'              COMMENT '状态: 1-启用, 2-禁用',
  `create_by`       BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '创建人 (存 sys_user.id)',
  `create_time`     DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `update_by`       BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '更新人 (存 sys_user.id)',
  `update_time`     DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at`      DATETIME(3)     NULL     DEFAULT NULL             COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_stage_status` (`apply_stage`, `status`),
  KEY `idx_category` (`category`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='内容净化规则表';
