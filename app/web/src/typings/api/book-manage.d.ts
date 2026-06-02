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

    /** book chapter */
    type BookChapter = Common.CommonRecord<{
      bookId: number;
      fileId: number;
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
      { bookId: number | null; fileId: number | null; chapterNo: number | null } & CommonSearchParams
    >;

    /** chapter content response */
    type ChapterContentResponse = BookChapter & {
      content: string;
    };

    // ==================== Book Chapter Rule ====================

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
      scopeType: string;
      bookId: number | null;
      pattern: string;
      titleGroup: number;
      minChapterLen: number;
      maxChapterLen: number;
      priority: number;
      description: string | null;
      status: string;
    }>;

    /** chapter rule request */
    type ChapterRuleRequest = {
      ruleName: string;
      scopeType: string;
      bookId?: number | null;
      pattern: string;
      titleGroup?: number;
      minChapterLen?: number;
      maxChapterLen?: number;
      priority?: number;
      description?: string | null;
      status?: string;
    };

    /** chapter rule search params */
    type ChapterRuleSearchParams = CommonType.RecordNullable<
      { ruleName: string; scopeType: string; bookId: number | null; status: string } & CommonSearchParams
    >;

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
  }
}
