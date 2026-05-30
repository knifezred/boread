package service

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"boread/internal/model"
)

// =====================================================================
// ChapterParser — 章节识别引擎
//
// 功能:
//  1. 支持正则规则匹配章节标题 (来自 book_chapter_rule)
//  2. 支持常见章节格式内置识别: "第X章", "Chapter X", 数字编号
//  3. 按章节切分并返回偏移索引
// =====================================================================

// ChapterMatch 单个章节匹配结果
type ChapterMatch struct {
	Title      string
	ByteOffset uint64
	LineNumber int
}

// ParseResult 整文件解析结果
type ParseResult struct {
	Chapters []ChapterSegment
}

// ChapterSegment 章节段
type ChapterSegment struct {
	Title      string
	ByteOffset uint64
	ByteLength uint32
	WordCount  uint32
}

// ChapterParser 章节解析器
type ChapterParser struct {
	rules []model.BookChapterRule
}

// NewChapterParser 创建解析器，传入有效规则（按优先级排序）
func NewChapterParser(rules []model.BookChapterRule) *ChapterParser {
	return &ChapterParser{rules: rules}
}

// Parse 解析原始文本，返回章节分段结果
func (p *ChapterParser) Parse(content []byte) *ParseResult {
	// 1. 按行扫描，找到所有章节标题位置
	matches := p.scanTitles(content)
	if len(matches) == 0 {
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
	segments := make([]ChapterSegment, 0, len(matches))
	for i, m := range matches {
		var endOffset uint64
		if i+1 < len(matches) {
			endOffset = matches[i+1].ByteOffset
		} else {
			endOffset = uint64(len(content))
		}
		seg := content[m.ByteOffset:endOffset]
		segments = append(segments, ChapterSegment{
			Title:      m.Title,
			ByteOffset: m.ByteOffset,
			ByteLength: uint32(len(seg)),
			WordCount:  countWords(seg),
		})
	}
	return &ParseResult{Chapters: segments}
}

// scanTitles 扫描文本中所有章节标题
func (p *ChapterParser) scanTitles(content []byte) []ChapterMatch {
	// 先尝试规则匹配
	if len(p.rules) > 0 {
		if matches := p.matchByRules(content); len(matches) > 0 {
			return matches
		}
	}
	// 回退到内置常见格式
	return p.matchBuiltin(content)
}

// matchByRules 使用配置的规则匹配章节标题
// 按优先级逐个规则尝试，第一个能匹配到内容的规则生效
func (p *ChapterParser) matchByRules(content []byte) []ChapterMatch {
	for _, rule := range p.rules {
		re, err := regexp.Compile(rule.Pattern)
		if err != nil {
			continue
		}
		matches := p.scanWithRegex(content, re, rule.TitleGroup)
		if len(matches) > 0 {
			return matches
		}
	}
	return nil
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
		matches = append(matches, ChapterMatch{
			Title:      title,
			ByteOffset: offset,
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
				matches = append(matches, ChapterMatch{
					Title:      trimmed,
					ByteOffset: offset,
					LineNumber: lineNum,
				})
				break
			}
		}
		offset += uint64(lineBytes) + 1
	}
	return matches
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
	rule        model.BookContentFilterRule
	regex       *regexp.Regexp
	keywords    []string // 关键词匹配时的列表
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
