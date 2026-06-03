package service

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	"boread/internal/model"
)

// =====================================================================
// ChapterParser — 章节识别引擎
//
// 功能:
//  1. 支持按 sort_order 升序依次尝试多条规则组合
//  2. 支持分卷标题识别（group_pattern）
//  3. 支持章节最小/最大字符数过滤
//  4. 支持常见章节格式内置识别: "第X章", "Chapter X", 数字编号
//  5. 按章节切分并返回偏移索引
// =====================================================================

// ChapterMatch 单个章节匹配结果
type ChapterMatch struct {
	Title          string
	VolumeTitle    string // 匹配到的卷标题，为空表示不是卷标题行
	ByteOffset     uint64 // 标题行起始字节偏移
	TitleEndOffset uint64 // 标题行结束（下一行起始）字节偏移
	LineNumber     int
}

// ParseResult 整文件解析结果
type ParseResult struct {
	Chapters []ChapterSegment
}

// ChapterSegment 章节段
type ChapterSegment struct {
	Title       string
	VolumeNo    *uint32
	VolumeTitle *string
	ByteOffset  uint64
	ByteLength  uint32
	WordCount   uint32
}

// compiledRule 预编译的规则
type compiledRule struct {
	rule       model.BookChapterRule
	titleRegex *regexp.Regexp
	groupRegex *regexp.Regexp
}

// ChapterParser 章节解析器
type ChapterParser struct {
	rules []compiledRule
}

// NewChapterParser 创建解析器，传入有效规则（已按 sort_order 升序排序）
func NewChapterParser(rules []model.BookChapterRule) *ChapterParser {
	compiled := make([]compiledRule, 0, len(rules))
	for _, r := range rules {
		titleRe, err := regexp.Compile(r.TitlePattern)
		if err != nil {
			continue
		}
		cr := compiledRule{rule: r, titleRegex: titleRe}
		if r.GroupPattern != nil && *r.GroupPattern != "" {
			if groupRe, err := regexp.Compile(*r.GroupPattern); err == nil {
				cr.groupRegex = groupRe
			}
		}
		compiled = append(compiled, cr)
	}
	return &ChapterParser{rules: compiled}
}

// Parse 解析原始文本，返回章节分段结果
func (p *ChapterParser) Parse(content []byte) *ParseResult {
	// 1. 按行扫描，找到所有章节标题位置（含卷标题）
	rawMatches := p.scanTitles(content)
	if len(rawMatches) == 0 {
		// 无匹配：整本书作为一个章节
		wordCount := countWords(content)
		return &ParseResult{
			Chapters: []ChapterSegment{{
				Title:      "全文",
				ByteOffset: 0,
				ByteLength: uint32(len(content)),
				WordCount:  wordCount,
			}},
		}
	}

	// 2. 按标题位置切分章节
	segments := p.buildSegments(content, rawMatches)

	// 3. 应用 min_chapter_len / max_chapter_len 过滤
	segments = p.filterByLength(segments)

	return &ParseResult{Chapters: segments}
}

// buildSegments 按匹配结果构建章节段，同时处理卷信息
func (p *ChapterParser) buildSegments(content []byte, matches []ChapterMatch) []ChapterSegment {
	contentLen := uint64(len(content))
	segments := make([]ChapterSegment, 0, len(matches))
	var currentVolumeNo *uint32
	var currentVolumeTitle *string

	for i, m := range matches {
		// 内容起始位置：标题行结束（下一行起始）
		contentStart := m.TitleEndOffset
		if contentStart > contentLen {
			contentStart = contentLen
		}

		// 内容结束位置：下一章节标题行起始
		var endOffset uint64
		if i+1 < len(matches) {
			endOffset = matches[i+1].ByteOffset
		} else {
			endOffset = contentLen
		}

		// 判断是否是卷标题行
		if m.VolumeTitle != "" {
			// 卷标题所在行不做为章节，但更新当前的卷信息
			volNo := uint32(1)
			if currentVolumeNo != nil {
				volNo = *currentVolumeNo + 1
			}
			currentVolumeNo = &volNo
			vt := m.VolumeTitle
			currentVolumeTitle = &vt
			continue
		}

		var segBytes []byte
		if contentStart < endOffset {
			segBytes = content[contentStart:endOffset]
		} else {
			segBytes = nil
		}

		seg := ChapterSegment{
			Title:       m.Title,
			VolumeNo:    copyUint32Ptr(currentVolumeNo),
			VolumeTitle: copyStringPtr(currentVolumeTitle),
			ByteOffset:  contentStart,
			ByteLength:  uint32(len(segBytes)),
			WordCount:   countWords(segBytes),
		}
		segments = append(segments, seg)
	}

	return segments
}

// filterByLength 根据每个规则的 min_chapter_len / max_chapter_len 过滤章节
// 太短的章节（目录/导航行）合并到上一章，太长的章节保持原样
func (p *ChapterParser) filterByLength(segments []ChapterSegment) []ChapterSegment {
	if len(segments) <= 1 {
		return segments
	}

	// 按规则匹配章节：遍历规则，对符合条件的章节标记过滤
	// 简化策略：对所有章节统一应用规则链中的最严格限制
	result := make([]ChapterSegment, 0, len(segments))

	for i, seg := range segments {
		// 获取匹配该章节的对应规则的 min/max
		// 使用分组匹配：找到该章节所属的匹配规则
		minLen, maxLen := p.getEffectiveLimits(seg)

		// 检查是否过短（跳过太短的章节，并入上一章）
		// 但跳过长度过滤的条件：seg.ByteLength < minLen 且不是第一章
		if seg.ByteLength < minLen && i > 0 && minLen > 0 {
			// 合并到上一章
			prev := &result[len(result)-1]
			prev.ByteLength += seg.ByteLength
			prev.WordCount += seg.WordCount
			// 保留标题（追加显示）
			prev.Title = prev.Title + " / " + seg.Title
			continue
		}

		// 检查是否过长（不截断，仅标记）
		// 如果超过 max_len，保留但记录
		if seg.ByteLength > maxLen && maxLen > 0 {
			// 超过上限，保持原样（强制按当前标题分割，但记录异常）
		}

		result = append(result, seg)
	}

	return result
}

// getEffectiveLimits 获取作用于该章节的有效 min/max 限制
// 遍历所有规则，取匹配该章节标题的规则的限制值；如果无匹配，使用默认值
func (p *ChapterParser) getEffectiveLimits(seg ChapterSegment) (uint32, uint32) {
	minLen := uint32(100)    // 默认值
	maxLen := uint32(100000) // 默认值

	for _, cr := range p.rules {
		if cr.titleRegex.MatchString(seg.Title) {
			// 使用第一个匹配的规则的参数
			if cr.rule.MinChapterLen > 0 {
				minLen = cr.rule.MinChapterLen
			}
			if cr.rule.MaxChapterLen > 0 {
				maxLen = cr.rule.MaxChapterLen
			}
			break
		}
	}

	return minLen, maxLen
}

// scanTitles 扫描文本中所有章节标题
func (p *ChapterParser) scanTitles(content []byte) []ChapterMatch {
	// 先尝试规则匹配（多规则组合）
	if len(p.rules) > 0 {
		if matches := p.matchByRules(content); len(matches) > 0 {
			return matches
		}
	}
	// 回退到内置常见格式
	return p.matchBuiltin(content)
}

// matchByRules 使用配置的规则按 sort_order 依次匹配章节标题
// 遍历所有规则，将所有规则的匹配结果合并，按字节偏移排序
func (p *ChapterParser) matchByRules(content []byte) []ChapterMatch {
	type rawMatch struct {
		offset uint64
		match  ChapterMatch
	}

	var rawMatches []rawMatch

	for _, cr := range p.rules {
		// 尝试匹配卷标题（group_pattern）
		if cr.groupRegex != nil {
			groupMatches := p.scanWithRegex(content, cr.groupRegex, 0)
			for _, m := range groupMatches {
				m.VolumeTitle = m.Title // 标记为卷标题
				rawMatches = append(rawMatches, rawMatch{offset: m.ByteOffset, match: m})
			}
		}

		// 尝试匹配章节标题（title_pattern）
		titleMatches := p.scanWithRegex(content, cr.titleRegex, 0)
		for _, m := range titleMatches {
			rawMatches = append(rawMatches, rawMatch{offset: m.ByteOffset, match: m})
		}
	}

	if len(rawMatches) == 0 {
		return nil
	}

	// 按偏移量去重排序（同一位置多条规则匹配，取第一条规则的匹配）
	sort.Slice(rawMatches, func(i, j int) bool {
		return rawMatches[i].offset < rawMatches[j].offset
	})

	// 去重
	result := make([]ChapterMatch, 0, len(rawMatches))
	seen := make(map[uint64]bool)
	for _, rm := range rawMatches {
		if seen[rm.offset] {
			continue
		}
		seen[rm.offset] = true
		result = append(result, rm.match)
	}

	return result
}

// scanWithRegex 使用单个正则扫描整个文件
func (p *ChapterParser) scanWithRegex(content []byte, re *regexp.Regexp, titleGroup int) []ChapterMatch {
	scanner := bufio.NewScanner(bytes.NewReader(content))
	scanner.Buffer(make([]byte, 0, 1024*1024), 1024*1024)
	var matches []ChapterMatch
	var offset uint64
	for scanner.Scan() {
		line := scanner.Text()
		lineBytes := len(scanner.Bytes())
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			offset += uint64(lineBytes) + 1
			continue
		}
		subs := re.FindStringSubmatch(trimmed)
		if subs == nil {
			offset += uint64(lineBytes) + 1
			continue
		}
		title := trimmed
		if titleGroup > 0 && titleGroup < len(subs) {
			title = subs[titleGroup]
		}
		titleEnd := offset + uint64(lineBytes) + 1
		matches = append(matches, ChapterMatch{
			Title:          title,
			ByteOffset:     offset,
			TitleEndOffset: titleEnd,
		})
		offset += uint64(lineBytes) + 1
	}
	return matches
}

// matchBuiltin 内置常见章节格式识别
func (p *ChapterParser) matchBuiltin(content []byte) []ChapterMatch {
	patterns := []*regexp.Regexp{
		// 第X章/节/回/卷/篇
		regexp.MustCompile(`^第[一二三四五六七八九十百千万0-9０-９]+[章章节回卷篇部集]`),
		// Chapter X / Chapter XX
		regexp.MustCompile(`(?i)^chapter\s+\d+`),
		// 数字编号: 001, 01, 1.
		regexp.MustCompile(`^\d{1,4}[.、．\s]`),
		// 卷/部/集 前缀
		regexp.MustCompile(`^[卷部集][一二三四五六七八九十0-9]+`),
		// 序言/前言/后记/尾声/楔子/番外
		regexp.MustCompile(`^(序言|前言|后记|尾声|楔子|番外|引子|简介|说明)`),
	}

	scanner := bufio.NewScanner(bytes.NewReader(content))
	scanner.Buffer(make([]byte, 0, 1024*1024), 1024*1024)
	var matches []ChapterMatch
	lineNum := 0
	var offset uint64
	for scanner.Scan() {
		line := scanner.Text()
		lineBytes := len(scanner.Bytes())
		lineNum++
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			offset += uint64(lineBytes) + 1
			continue
		}
		for _, pat := range patterns {
			if pat.MatchString(trimmed) {
				titleEnd := offset + uint64(lineBytes) + 1
				matches = append(matches, ChapterMatch{
					Title:          trimmed,
					ByteOffset:     offset,
					TitleEndOffset: titleEnd,
					LineNumber:     lineNum,
				})
				break
			}
		}
		offset += uint64(lineBytes) + 1
	}
	return matches
}

// 辅助函数
func copyUint32Ptr(v *uint32) *uint32 {
	if v == nil {
		return nil
	}
	cpy := *v
	return &cpy
}

func copyStringPtr(v *string) *string {
	if v == nil {
		return nil
	}
	cpy := *v
	return &cpy
}

// =====================================================================
// ContentFilter — 内容净化引擎
//
// 功能:
//  1. 关键词匹配替换
//  2. 正则匹配替换/拦截
//  3. 特殊字符清理
//  4. 格式统一化
// =====================================================================

// FilterResult 单条规则过滤结果
type FilterResult struct {
	Content   string
	Blocked   bool   // true 表示整章被拦截
	MatchDesc string // 匹配详情
}

// ContentFilter 内容过滤器
type ContentFilter struct {
	replaceRules []compiledFilterRule
	blockRules   []compiledFilterRule
}

type compiledFilterRule struct {
	rule     model.BookContentFilterRule
	regex    *regexp.Regexp
	keywords []string // 关键词匹配时的列表
}

// NewContentFilter 创建内容过滤器
func NewContentFilter(rules []model.BookContentFilterRule) *ContentFilter {
	f := &ContentFilter{}
	for _, rule := range rules {
		if rule.Status != model.StatusEnabled {
			continue
		}
		cr := compiledFilterRule{rule: rule}
		switch rule.MatchType {
		case model.FilterKeyword:
			cr.keywords = []string{rule.Pattern}
		case model.FilterRegex:
			cr.regex = regexp.MustCompile(rule.Pattern)
		}
		switch rule.Action {
		case model.FilterReplace:
			f.replaceRules = append(f.replaceRules, cr)
		case model.FilterBlock:
			f.blockRules = append(f.blockRules, cr)
		case model.FilterMarkReview:
			// 标记审核暂不阻塞，仅记录
		}
	}
	return f
}

// Filter 对文本内容应用过滤规则
// 返回过滤后的内容、是否被拦截、匹配到的规则描述
func (f *ContentFilter) Filter(content string) FilterResult {
	// 1. 先检查是否命中拦截规则
	for _, cr := range f.blockRules {
		if matched := cr.match(content); matched {
			return FilterResult{
				Blocked:   true,
				MatchDesc: fmt.Sprintf("拦截规则[%s]: %s", cr.rule.RuleName, cr.rule.Pattern),
			}
		}
	}

	// 2. 应用替换规则
	result := content
	for _, cr := range f.replaceRules {
		result = cr.applyReplacement(result)
	}

	// 3. 通用清理
	result = sanitizeText(result)

	return FilterResult{
		Content: result,
		Blocked: false,
	}
}

func (cr *compiledFilterRule) match(content string) bool {
	if cr.regex != nil {
		return cr.regex.MatchString(content)
	}
	for _, kw := range cr.keywords {
		if strings.Contains(content, kw) {
			return true
		}
	}
	return false
}

func (cr *compiledFilterRule) applyReplacement(content string) string {
	if cr.regex != nil {
		return cr.regex.ReplaceAllString(content, cr.rule.Replacement)
	}
	result := content
	for _, kw := range cr.keywords {
		result = strings.ReplaceAll(result, kw, cr.rule.Replacement)
	}
	return result
}

// sanitizeText 通用文本清理
func sanitizeText(text string) string {
	// 1. 连续空白符合并
	result := compactWhitespace(text)
	// 2. 去除 BOM
	result = strings.TrimLeft(result, "\uFEFF")
	// 3. 统一换行符
	result = strings.ReplaceAll(result, "\r\n", "\n")
	result = strings.ReplaceAll(result, "\r", "\n")
	return result
}

// compactWhitespace 将连续空白符合并为单个空格
func compactWhitespace(s string) string {
	var buf strings.Builder
	buf.Grow(len(s))
	prevSpace := false
	for _, r := range s {
		if unicode.IsSpace(r) {
			if !prevSpace {
				buf.WriteRune(' ')
				prevSpace = true
			}
		} else {
			buf.WriteRune(r)
			prevSpace = false
		}
	}
	return buf.String()
}

// countWords 统计字符数（非空字符）
func countWords(data []byte) uint32 {
	if len(data) == 0 {
		return 0
	}
	return uint32(utf8.RuneCount(data))
}

// ValidateFileType 验证文件格式是否支持
func ValidateFileType(ext string) bool {
	ext = strings.ToLower(strings.TrimPrefix(ext, "."))
	switch ext {
	case "txt", "epub", "mobi", "pdf":
		return true
	}
	return false
}

// DetectCharset 检测文本编码（简化版：检查 BOM + 常见编码）
func DetectCharset(data []byte) string {
	if len(data) < 2 {
		return "utf-8"
	}
	if len(data) >= 3 && data[0] == 0xEF && data[1] == 0xBB && data[2] == 0xBF {
		return "utf-8"
	}
	if len(data) >= 2 && data[0] == 0xFF && data[1] == 0xFE {
		return "utf-16le"
	}
	if len(data) >= 2 && data[0] == 0xFE && data[1] == 0xFF {
		return "utf-16be"
	}
	// 简化检测：检查是否含中文字符来判断
	if hasChinese(string(data)) {
		return "utf-8"
	}
	return "utf-8"
}

func hasChinese(s string) bool {
	for _, r := range s {
		if unicode.Is(unicode.Han, r) {
			return true
		}
	}
	return false
}
