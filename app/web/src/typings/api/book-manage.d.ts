declare namespace Api {
  /**
   * namespace BookManage
   *
   * backend api module: "bookManage" — all book-related types
   */
  namespace BookManage {
    type CommonSearchParams = Pick<
      Common.PaginatingCommonParams,
      "current" | "size"
    >;

    // ==================== Book Category ====================

    /** book category */
    type BookCategory = Common.CommonRecord<{
      parentId: number;
      categoryName: string;
      categoryCode: string;
      description: string | null;
      sortOrder: number;
      isHot: boolean;
      children?: BookCategory[] | null;
    }>;

    /** category search params */
    type CategorySearchParams = CommonType.RecordNullable<
      Pick<BookCategory, "categoryName" | "categoryCode" | "parentId" | "status" | "isHot"> & CommonSearchParams
    >;

    /** category list */
    type BookCategoryList = Common.PaginatingQueryRecord<BookCategory>;

    /** hot category item */
    type HotCategoryItem = {
      id: number;
      categoryName: string;
      categoryCode: string;
    };

    // ==================== Book Tag ====================

    /** book tag */
    type BookTag = {
      id: number;
      tagName: string;
      description: string | null;
      usageCount: number;
    };

    /** tag search params */
    type TagSearchParams = CommonType.RecordNullable<
      { tagName: string } & CommonSearchParams
    >;

    /** tag list */
    type BookTagList = Common.PaginatingQueryRecord<BookTag>;

    // ==================== Book (Novel) ====================

    /** serial status: 1-ongoing 2-finished 3-stopped */
    type SerialStatus = "1" | "2" | "3";
    /** visibility: 1-public 2-private 3-dept */
    type Visibility = "1" | "2" | "3";
    /** book listing status: 1-listed 2-unlisted 3-reviewing 4-rejected */
    type BookListingStatus = "1" | "2" | "3" | "4";
    /** aggregate status: 1-single 2-merging 3-done */
    type AggregateStatus = "1" | "2" | "3";

    /** book (novel) */
    type Book = Common.CommonRecord<{
      title: string;
      author: string;
      cover: string | null;
      intro: string | null;
      categoryId: number | null;
      categoryName?: string;
      language: string;
      serialStatus: SerialStatus;
      visibility: Visibility;
      primaryFileId: number | null;
      totalChapters: number;
      totalWords: number;
      aggregateStatus: AggregateStatus;
      avgRating: number;
      ratingCount: number;
      deptId: number | null;
      tagIds?: number[];
      tags?: { id: number; tagName: string }[];
      latestChapterTitle?: string;
      latestChapterNo?: number;
    }>;

    /** book search params */
    type BookSearchParams = CommonType.RecordNullable<
      {
        title: string;
        author: string;
        categoryId: number | null;
        status: BookListingStatus;
        visibility: Visibility;
        serialStatus: SerialStatus;
        tagId: number | null;
        minWords: number | null;
        maxWords: number | null;
        updateTimeFrom: string | null;
        updateTimeTo: string | null;
      } & CommonSearchParams
    >;

    /** book list */
    type BookList = Common.PaginatingQueryRecord<Book>;

    /** book create/update request */
    type BookRequest = {
      title: string;
      author: string;
      cover: string | null;
      intro: string | null;
      categoryId: number | null;
      language: string;
      serialStatus: SerialStatus;
      visibility: Visibility;
      tagIds: number[];
    };

    /** book status update request */
    type BookUpdateStatusRequest = {
      status: BookListingStatus;
    };

    // ==================== Book File & Upload ====================

    /** file upload response */
    type FileUploadResponse = {
      uploadId: number;
      originalName: string;
      fileSize: number;
      sourceFormat: string | null;
      suggestedTitle: string;
      suggestedAuthor: string;
      matchedBookId: number | null;
      matchedBookTitle: string;
    };

    type ConfirmImportRequest = {
      uploadId: number;
      title: string;
      author: string;
    };

    type ConfirmImportResponse = {
      uploadId: number;
      bookId: number;
      bookTitle: string;
      bookAuthor: string;
      chapterCount: number;
      parseStatus: string;
      parseMessage: string | null;
    };

    /** scan result */
    type ScanResult = {
      uploadId: number;
      originalName: string;
      fileSize: number;
      parseStatus: string;
      parseMessage: string | null;
      bookId: number | null;
      chapterCount: number | null;
    };

    /** scan all response */
    type ScanAllResponse = {
      results: ScanResult[];
      success: number;
      failed: number;
    };

    type ScanPathResponse = {
      total: number;
      imported: number;
      failed: number;
      results: ScanResult[];
    };

    /** book upload record */
    type BookUpload = Common.CommonRecord<{
      bookId: number | null;
      originalName: string;
      filePath: string;
      fileSize: number;
      fileMd5: string | null;
      sourceFormat: string | null;
      parseStatus: string;
      parseMessage: string | null;
      chapterCount: number | null;
    }>;

    /** upload search params */
    type UploadSearchParams = CommonType.RecordNullable<
      { originalName: string; parseStatus: string; bookId: number | null } & CommonSearchParams
    >;

    /** upload list */
    type BookUploadList = Common.PaginatingQueryRecord<BookUpload>;

    /** book file record */
    type BookFile = Common.CommonRecord<{
      bookId: number;
      originalName: string;
      sourceType: string;
      sourceFormat: string | null;
      sourceFileUrl: string | null;
      contentPath: string | null;
      contentSize: number;
      contentMd5: string | null;
      contentCharset: string;
      contentVersion: number;
      chapterCount: number;
      isPrimary: boolean;
      fileStatus: string;
      parseMessage: string | null;
    }>;

    /** file search params */
    type FileSearchParams = CommonType.RecordNullable<
      { bookId: number | null; fileStatus: string; sourceType: string } & CommonSearchParams
    >;

    /** file list */
    type BookFileList = Common.PaginatingQueryRecord<BookFile>;

    /** book chapter */
    type BookChapter = Common.CommonRecord<{
      bookId: number;
      fileId: number;
      volumeNo: number | null;
      volumeTitle: string | null;
      chapterNo: number;
      title: string;
      byteOffset: number;
      byteLength: number;
      wordCount: number;
      isVip: boolean;
      status: string;
    }>;

    /** chapter search params */
    type ChapterSearchParams = CommonType.RecordNullable<
      { bookId: number | null; fileId: number | null; chapterNo: number | null; title: string; status: string } & CommonSearchParams
    >;

    /** chapter list */
    type BookChapterList = Common.PaginatingQueryRecord<BookChapter>;

    /** chapter content response */
    type ChapterContentResponse = BookChapter & {
      content: string;
    };

    // ==================== Chapter Management ====================

    /** update chapter title request */
    type ChapterTitleUpdateRequest = {
      title: string;
    };

    /** batch update chapter titles request */
    type ChapterTitleBatchRequest = {
      ids: number[];
      title: string;
    };

    /** batch update chapter status request */
    type ChapterStatusBatchRequest = {
      ids: number[];
      status: string;
    };

    /** merge chapters request */
    type ChapterMergeRequest = {
      bookId: number;
      targetId: number;
      sourceIds: number[];
    };

    /** format chapter numbers request */
    type ChapterFormatRequest = {
      ids: number[];
    };

    /** save chapter content request */
    type ChapterContentSaveRequest = {
      bookId: number;
      content: string;
    };

    /** re-parse chapters request */
    type ReParseRequest = {
      bookId: number;
      ruleId?: number;
    };

    /** re-parse chapters response */
    type ReParseResponse = {
      bookId: number;
      bookTitle: string;
      oldCount: number;
      newCount: number;
      totalWords: number;
    };

    /** book chapter rule */
    type BookChapterRule = Common.CommonRecord<{
      ruleName: string;
      ruleType: string;
      userId: number | null;
      titlePattern: string;
      groupPattern: string | null;
      minChapterLen: number;
      maxChapterLen: number;
      sortOrder: number;
      description: string | null;
      status: string;
    }>;

    /** chapter rule request */
    type ChapterRuleRequest = {
      ruleName: string;
      ruleType: string;
      titlePattern: string;
      groupPattern?: string | null;
      minChapterLen?: number;
      maxChapterLen?: number;
      sortOrder?: number;
      description?: string | null;
      status?: string;
    };

    /** chapter rule search params */
    type ChapterRuleSearchParams = CommonType.RecordNullable<
      { ruleName: string; ruleType: string; userId: number | null; status: string } & CommonSearchParams
    >;

    /** chapter rule list */
    type BookChapterRuleList = Common.PaginatingQueryRecord<BookChapterRule>;

    // ==================== Book Chapter Rule Binding ====================

    /** chapter rule bind request */
    type ChapterRuleBindRequest = {
      bookId: number;
      ruleId: number;
    };

    /** chapter rule bind response */
    type ChapterRuleBindResponse = {
      id: number;
      bookId: number;
      readerId: number;
      ruleId: number;
      ruleName: string;
      createTime: string;
    };

    // ==================== Book Content Filter Rule ====================

    /** book content filter rule */
    type BookContentFilterRule = Common.CommonRecord<{
      ruleName: string;
      matchType: string;
      pattern: string;
      action: string;
      replacement: string;
      applyStage: string;
      category: string | null;
      severity: string;
      description: string | null;
      status: string;
    }>;

    /** filter rule request */
    type FilterRuleRequest = {
      ruleName: string;
      matchType: string;
      pattern: string;
      action: string;
      replacement?: string;
      applyStage: string;
      category?: string | null;
      severity?: string;
      description?: string | null;
      status?: string;
    };

    /** filter rule search params */
    type FilterRuleSearchParams = CommonType.RecordNullable<
      { ruleName: string; applyStage: string; category: string; status: string } & CommonSearchParams
    >;

    /** filter rule list */
    type BookContentFilterRuleList = Common.PaginatingQueryRecord<BookContentFilterRule>;

    // ==================== Book Card Component Types ====================

    /** BookFilter 筛选条件值 */
    interface BookFilterParams {
      categoryId: number;
      serialStatus: string;
      wordCount: string;
      tagId: string;
      updateTime: string;
      title: string;
      sortBy: string;
      sortOrder: string;
    }

    // ==================== Reader Bookshelf ====================

    /** 书架条目 */
    type BookshelfItem = {
      id: number;
      readerId: number;
      bookId: number;
      groupName: string;
      isTop: boolean;
      lastReadTime: string | null;
      addTime: string;
      createTime: string;
      bookTitle: string;
      bookAuthor: string;
      bookCover: string | null;
      totalChapters: number;
      totalWords: number;
      chapterId: number | null;
      chapterNo: number | null;
      position: number | null;
      readPercent: number;
      readDuration: number;
    };

    /** 书架搜索参数 */
    type BookshelfSearchParams = CommonType.RecordNullable<
      {
        groupName: string;
        keyword: string;
      } & CommonSearchParams
    >;

    /** 书架列表 */
    type BookshelfItemList = Common.PaginatingQueryRecord<BookshelfItem>;

    /** 添加到书架请求 */
    type BookshelfRequest = {
      bookId: number;
      groupName?: string;
    };

    /** 更新书架请求 */
    type BookshelfUpdateRequest = {
      groupName?: string;
      isTop?: boolean;
    };

    /** 书架分组项 */
    type BookshelfGroupItem = {
      groupName: string;
      bookCount: number;
    };

    // ==================== Reader Read Progress ====================

    /** 上报阅读进度请求 */
    type ReadProgressRequest = {
      fileId?: number | null;
      chapterId: number;
      chapterNo: number;
      position: number;
      percent: number;
      readDuration: number;
    };

    /** 阅读进度响应 */
    type ReadProgressResponse = Common.CommonRecord<{
      readerId: number;
      bookId: number;
      fileId: number | null;
      chapterId: number;
      chapterNo: number;
      position: number;
      percent: number;
      readDuration: number;
      lastReadTime: string;
    }>;

    // ==================== Reader Read Event ====================

    /** 上报阅读事件请求 */
    type ReadEventRequest = {
      bookId: number;
      chapterId: number;
      sessionId?: string;
      durationSec: number;
      wordCount?: number;
      deviceType?: string;
    };
  }
}
