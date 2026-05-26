-- =====================================================================
-- Business Schema —— 小说阅读平台业务表
-- Engine: InnoDB / utf8mb4_unicode_ci
-- Scope:  阅读者、书库、书架、章节、笔记/划线、评论/书评、分类标签
-- Notes:
--   1. 后台 sys_user 管运营，前台 reader 管阅读，分离
--   2. 业务表统一: owner_id (创建者) + dept_id (归属部门) + deleted_at
--      用于配合 sys_role.data_scope 实现数据权限过滤
--   3. 章节内容存本地文件 (storage/books/{hash分桶}/{book_id}.txt)
--      DB 只存 byte_offset + byte_length 索引, 用 pread 切片读章节
--   4. 笔记/划线合并为 reader_note + note_type 字段, 消除特殊情况
-- =====================================================================

USE `boread`;

-- ---------------------------------------------------------------------
-- Table: reader (前台阅读者)
-- 与 sys_user 解耦, 字段语义不同, 权限模型不同
-- ---------------------------------------------------------------------
DROP TABLE IF EXISTS `reader`;
CREATE TABLE `reader` (
  `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT          COMMENT '主键 id',
  `username`        VARCHAR(64)     NOT NULL                          COMMENT '登录账号',
  `password`        VARCHAR(128)    NOT NULL                          COMMENT '密码 (bcrypt)',
  `pwd_updated_at`  DATETIME(3)     NULL     DEFAULT NULL             COMMENT '密码最后修改时间',
  `pwd_error_count` SMALLINT UNSIGNED NOT NULL DEFAULT 0              COMMENT '连续密码错误次数 (登录成功重置)',
  `locked_until`    DATETIME(3)     NULL     DEFAULT NULL             COMMENT '账号锁定到何时 (NULL=未锁)',
  `nickname`        VARCHAR(64)     NOT NULL DEFAULT ''               COMMENT '昵称',
  `avatar`          VARCHAR(255)    NULL     DEFAULT NULL             COMMENT '头像 URL',
  `email`           VARCHAR(128)    NULL     DEFAULT NULL             COMMENT '邮箱',
  `phone`           VARCHAR(20)     NULL     DEFAULT NULL             COMMENT '手机号',
  `gender`          CHAR(1)         NULL     DEFAULT NULL             COMMENT '性别: 1-男, 2-女, 0-保密',
  `signature`       VARCHAR(255)    NULL     DEFAULT NULL             COMMENT '个性签名',
  `read_pref`       JSON            NULL     DEFAULT NULL             COMMENT '阅读偏好: 字号/主题/翻页等',
  `last_login_time` DATETIME(3)     NULL     DEFAULT NULL             COMMENT '最后登录时间',
  `last_login_ip`   VARCHAR(64)     NULL     DEFAULT NULL             COMMENT '最后登录 IP',
  `status`          CHAR(1)         NOT NULL DEFAULT '1'              COMMENT '状态: 1-启用, 2-禁用',
  `version`         BIGINT UNSIGNED NOT NULL DEFAULT 0                COMMENT '乐观锁版本号',
  `create_time`     DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `update_time`     DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at`      DATETIME(3)     NULL     DEFAULT NULL             COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_username` (`username`),
  KEY `idx_email` (`email`),
  KEY `idx_phone` (`phone`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='前台阅读者表';

-- ---------------------------------------------------------------------
-- Table: book_category (书籍分类)
-- 自关联树, 例: 文学 -> 小说 -> 玄幻
-- ---------------------------------------------------------------------
DROP TABLE IF EXISTS `book_category`;
CREATE TABLE `book_category` (
  `id`            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT          COMMENT '主键 id',
  `parent_id`     BIGINT UNSIGNED NOT NULL DEFAULT 0                COMMENT '父分类 id',
  `category_name` VARCHAR(64)     NOT NULL                          COMMENT '分类名称',
  `category_code` VARCHAR(64)     NOT NULL                          COMMENT '分类编码',
  `sort_order`    INT             NOT NULL DEFAULT 0                COMMENT '排序',
  `status`        CHAR(1)         NOT NULL DEFAULT '1'              COMMENT '状态: 1-启用, 2-禁用',
  `create_by`     BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '创建人 (存 sys_user.id)',
  `create_time`   DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
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
  `usage_count` INT UNSIGNED    NOT NULL DEFAULT 0      COMMENT '引用计数 (冗余, 便于排序热门)',
  `create_time` DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_tag_name` (`tag_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='书籍标签表';

-- ---------------------------------------------------------------------
-- Table: book (书库 / 小说主表)
-- uploader_type + uploader_id 区分后台管理员上传还是读者上传
-- owner_id + dept_id 用于数据权限过滤
-- ---------------------------------------------------------------------
DROP TABLE IF EXISTS `book`;
CREATE TABLE `book` (
  `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT          COMMENT '主键 id',
  `title`           VARCHAR(255)    NOT NULL                          COMMENT '书名',
  `author`          VARCHAR(128)    NOT NULL DEFAULT ''               COMMENT '作者',
  `cover`           VARCHAR(512)    NULL     DEFAULT NULL             COMMENT '封面 URL',
  `intro`           TEXT            NULL     DEFAULT NULL             COMMENT '简介',
  `category_id`     BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '分类 id (book_category.id)',
  `source_type`     CHAR(1)         NOT NULL DEFAULT '1'              COMMENT '来源: 1-用户上传, 2-管理员上传, 3-外部导入',
  `source_format`   VARCHAR(16)     NULL     DEFAULT NULL             COMMENT '源格式: txt, epub, mobi, pdf...',
  `source_file_url` VARCHAR(512)    NULL     DEFAULT NULL             COMMENT '源文件 URL (原始上传, 可选保留)',
  `content_path`    VARCHAR(512)    NULL     DEFAULT NULL             COMMENT '正文文件相对路径 (storage/books 下), 解析完成后写入',
  `content_size`    BIGINT UNSIGNED NOT NULL DEFAULT 0                COMMENT '正文文件大小 (bytes)',
  `content_md5`     CHAR(32)        NULL     DEFAULT NULL             COMMENT '正文文件 MD5 (秒传 + 完整性校验)',
  `content_charset` VARCHAR(16)     NOT NULL DEFAULT 'utf-8'          COMMENT '正文文件字符编码 (txt 原文可能 GBK)',
  `total_chapters`  INT UNSIGNED    NOT NULL DEFAULT 0                COMMENT '总章节数 (冗余)',
  `total_words`     INT UNSIGNED    NOT NULL DEFAULT 0                COMMENT '总字数 (冗余)',
  `language`        VARCHAR(16)     NOT NULL DEFAULT 'zh-CN'          COMMENT '语言',
  `serial_status`   CHAR(1)         NOT NULL DEFAULT '1'              COMMENT '连载状态: 1-连载中, 2-已完结, 3-断更',
  `visibility`      CHAR(1)         NOT NULL DEFAULT '1'              COMMENT '可见性: 1-公开(全平台), 2-仅自己, 3-部门内',
  `owner_id`        BIGINT UNSIGNED NOT NULL                          COMMENT '上传者 id (数据权限)',
  `owner_type`      CHAR(1)         NOT NULL DEFAULT '2'              COMMENT '上传者类型: 1-sys_user, 2-reader',
  `dept_id`         BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '上传者所属部门 id (数据权限, owner_type=1 时有效)',
  `view_count`      INT UNSIGNED    NOT NULL DEFAULT 0                COMMENT '阅读次数 (冗余)',
  `like_count`      INT UNSIGNED    NOT NULL DEFAULT 0                COMMENT '点赞数',
  `status`          CHAR(1)         NOT NULL DEFAULT '1'              COMMENT '审核状态: 1-已上架, 2-下架, 3-审核中, 4-审核拒绝',
  `create_time`     DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `update_time`     DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at`      DATETIME(3)     NULL     DEFAULT NULL             COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_title` (`title`),
  KEY `idx_author` (`author`),
  KEY `idx_category` (`category_id`),
  KEY `idx_owner` (`owner_type`, `owner_id`),
  KEY `idx_dept_id` (`dept_id`),
  KEY `idx_status_visibility` (`status`, `visibility`),
  KEY `idx_content_md5` (`content_md5`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='书库 / 小说表';

-- ---------------------------------------------------------------------
-- Table: book_tag_rel (书籍-标签 关联)
-- ---------------------------------------------------------------------
DROP TABLE IF EXISTS `book_tag_rel`;
CREATE TABLE `book_tag_rel` (
  `id`      BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键 id',
  `book_id` BIGINT UNSIGNED NOT NULL                COMMENT '书籍 id',
  `tag_id`  BIGINT UNSIGNED NOT NULL                COMMENT '标签 id',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_book_tag` (`book_id`, `tag_id`),
  KEY `idx_tag_id` (`tag_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='书籍标签关联表';

-- ---------------------------------------------------------------------
-- Table: book_chapter (章节索引)
-- 设计要点:
--   内容不入库, 整本书一个文件 (book.content_path)
--   章节用 byte_offset + byte_length 在该文件中切片
--   pread(fd, offset, length) 单次系统调用即可读章节, 零拷贝
-- ---------------------------------------------------------------------
DROP TABLE IF EXISTS `book_chapter`;
CREATE TABLE `book_chapter` (
  `id`            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT          COMMENT '主键 id',
  `book_id`       BIGINT UNSIGNED NOT NULL                          COMMENT '书籍 id',
  `chapter_no`    INT UNSIGNED    NOT NULL                          COMMENT '章节序号 (1 开始)',
  `title`         VARCHAR(255)    NOT NULL                          COMMENT '章节标题',
  `byte_offset`   BIGINT UNSIGNED NOT NULL                          COMMENT '在 book.content_path 文件中的起始字节偏移',
  `byte_length`   INT UNSIGNED    NOT NULL                          COMMENT '章节字节长度 (UTF-8 编码后)',
  `word_count`    INT UNSIGNED    NOT NULL DEFAULT 0                COMMENT '字符数 (展示用, 与字节数不同)',
  `is_vip`        TINYINT(1)      NOT NULL DEFAULT 0                COMMENT '是否 VIP 章节',
  `status`        CHAR(1)         NOT NULL DEFAULT '1'              COMMENT '状态: 1-发布, 2-草稿, 3-下架',
  `create_time`   DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `update_time`   DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at`    DATETIME(3)     NULL     DEFAULT NULL             COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_book_chapter_no` (`book_id`, `chapter_no`),
  KEY `idx_book_id` (`book_id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='书籍章节索引表 (内容存文件)';

-- ---------------------------------------------------------------------
-- Table: book_upload (上传/解析任务)
-- 异步解析 epub/txt 的诊断记录, 失败原因可追溯
-- ---------------------------------------------------------------------
DROP TABLE IF EXISTS `book_upload`;
CREATE TABLE `book_upload` (
  `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT          COMMENT '主键 id',
  `uploader_id`     BIGINT UNSIGNED NOT NULL                          COMMENT '上传者 id',
  `uploader_type`   CHAR(1)         NOT NULL                          COMMENT '上传者类型: 1-sys_user, 2-reader',
  `original_name`   VARCHAR(255)    NOT NULL                          COMMENT '原始文件名',
  `file_url`        VARCHAR(512)    NOT NULL                          COMMENT '文件存储 URL',
  `file_size`       BIGINT UNSIGNED NOT NULL DEFAULT 0                COMMENT '文件大小 (bytes)',
  `file_md5`        CHAR(32)        NULL     DEFAULT NULL             COMMENT '文件 MD5 (秒传判重)',
  `file_format`     VARCHAR(16)     NOT NULL                          COMMENT '格式: txt, epub, mobi, pdf...',
  `book_id`         BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '解析成功后关联的 book.id',
  `parse_status`    CHAR(1)         NOT NULL DEFAULT '1'              COMMENT '解析状态: 1-待处理, 2-处理中, 3-成功, 4-失败',
  `parse_message`   VARCHAR(512)    NULL     DEFAULT NULL             COMMENT '解析结果说明 / 失败原因',
  `create_time`     DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `update_time`     DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_uploader` (`uploader_type`, `uploader_id`),
  KEY `idx_book_id` (`book_id`),
  KEY `idx_file_md5` (`file_md5`),
  KEY `idx_parse_status` (`parse_status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='书籍上传/解析任务表';

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
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_reader_book` (`reader_id`, `book_id`),
  KEY `idx_book_id` (`book_id`),
  KEY `idx_last_read` (`reader_id`, `last_read_time`)
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
  `chapter_id`      BIGINT UNSIGNED NOT NULL                          COMMENT '当前章节 id',
  `chapter_no`      INT UNSIGNED    NOT NULL                          COMMENT '当前章节序号 (冗余, 便于排序展示)',
  `position`        INT UNSIGNED    NOT NULL DEFAULT 0                COMMENT '章内位置 (字符偏移量)',
  `percent`         DECIMAL(5,2)    NOT NULL DEFAULT 0.00             COMMENT '全书进度百分比 0.00-100.00',
  `read_duration`   INT UNSIGNED    NOT NULL DEFAULT 0                COMMENT '累计阅读时长 (秒)',
  `last_read_time`  DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '最后阅读时间',
  `create_time`     DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_reader_book` (`reader_id`, `book_id`),
  KEY `idx_chapter` (`chapter_id`)
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
  `create_time`     DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
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
  `status`        CHAR(1)         NOT NULL DEFAULT '1'              COMMENT '状态: 1-正常, 2-隐藏, 3-审核中',
  `create_time`   DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `update_time`   DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at`    DATETIME(3)     NULL     DEFAULT NULL             COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_book_id` (`book_id`),
  KEY `idx_reader_id` (`reader_id`),
  KEY `idx_deleted_at` (`deleted_at`)
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
  `status`        CHAR(1)         NOT NULL DEFAULT '1'              COMMENT '状态: 1-正常, 2-隐藏, 3-审核中',
  `create_time`   DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `deleted_at`    DATETIME(3)     NULL     DEFAULT NULL             COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_chapter` (`chapter_id`),
  KEY `idx_book_id` (`book_id`),
  KEY `idx_reader_id` (`reader_id`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='章节评论表';

-- ---------------------------------------------------------------------
-- Table: book_parse_rule (章节识别规则)
-- 用于解析 txt 时切分章节, 支持多条规则按优先级匹配
-- scope_type 区分全局默认还是单书覆盖
-- ---------------------------------------------------------------------
DROP TABLE IF EXISTS `book_parse_rule`;
CREATE TABLE `book_parse_rule` (
  `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT          COMMENT '主键 id',
  `rule_name`       VARCHAR(64)     NOT NULL                          COMMENT '规则名称',
  `scope_type`      CHAR(1)         NOT NULL DEFAULT '1'              COMMENT '作用域: 1-全局默认, 2-单书覆盖',
  `book_id`         BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '关联书 id (scope_type=2 时必填)',
  `pattern`         VARCHAR(512)    NOT NULL                          COMMENT '章节标题正则 (Go RE2 语法), 例: ^第[零一二三四五六七八九十百千0-9]+[章回节]\\s.*$',
  `title_group`     INT             NOT NULL DEFAULT 0                COMMENT '正则中作为标题的捕获组序号, 0=整行',
  `min_chapter_len` INT UNSIGNED    NOT NULL DEFAULT 100              COMMENT '章节最小字符数 (过滤误匹配如目录页)',
  `max_chapter_len` INT UNSIGNED    NOT NULL DEFAULT 100000           COMMENT '章节最大字符数 (过大可能是未切分)',
  `priority`        INT             NOT NULL DEFAULT 0                COMMENT '优先级 (越大越先匹配)',
  `description`     VARCHAR(255)    NULL     DEFAULT NULL             COMMENT '说明 / 示例',
  `status`          CHAR(1)         NOT NULL DEFAULT '1'              COMMENT '状态: 1-启用, 2-禁用',
  `create_by`       BIGINT UNSIGNED NULL     DEFAULT NULL             COMMENT '创建人 (存 sys_user.id)',
  `create_time`     DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `update_time`     DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at`      DATETIME(3)     NULL     DEFAULT NULL             COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_scope` (`scope_type`, `book_id`),
  KEY `idx_priority` (`priority`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='章节识别规则表';

-- ---------------------------------------------------------------------
-- Table: content_filter_rule (内容净化规则)
-- 入库时(解析阶段) 或 出库时(读章节时) 应用
-- action: 1-替换  2-拦截整章  3-标记审核
-- ---------------------------------------------------------------------
DROP TABLE IF EXISTS `content_filter_rule`;
CREATE TABLE `content_filter_rule` (
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
  `update_time`     DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at`      DATETIME(3)     NULL     DEFAULT NULL             COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_stage_status` (`apply_stage`, `status`),
  KEY `idx_category` (`category`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='内容净化规则表';

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
  `word_count`     INT UNSIGNED    NOT NULL DEFAULT 0                COMMENT '本次阅读字数',
  `event_date`     DATE            NOT NULL                          COMMENT '事件日期 (聚合维度, 冗余于 event_time)',
  `event_time`     DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '事件时间',
  PRIMARY KEY (`id`),
  KEY `idx_reader_date` (`reader_id`, `event_date`),
  KEY `idx_book_id` (`book_id`),
  KEY `idx_event_date` (`event_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='阅读事件原子表';

-- ---------------------------------------------------------------------
-- Table: reader_read_stat_daily (按日聚合)
-- 凌晨定时任务把昨日 reader_read_event 聚合到这张表
-- 周/月/年/总 = SUM(GROUP BY) 此表, 无需额外建表
--   按周:  WHERE stat_date BETWEEN ... GROUP BY YEARWEEK(stat_date)
--   按月:  WHERE stat_date BETWEEN ... GROUP BY DATE_FORMAT(stat_date,'%Y-%m')
--   按年:  GROUP BY YEAR(stat_date)
--   总计:  SUM 不带 GROUP BY
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
  `create_time`    DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `update_time`    DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_reader_date` (`reader_id`, `stat_date`),
  KEY `idx_stat_date` (`stat_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='读者阅读日统计';

-- ---------------------------------------------------------------------
-- Table: reader_read_stat_book (按读者-书-日聚合, 用于"某本书读了多久")
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
  `create_time`    DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `update_time`    DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_reader_book_date` (`reader_id`, `book_id`, `stat_date`),
  KEY `idx_book_date` (`book_id`, `stat_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='读者-书-日 阅读统计';