package service

import (
	"bufio"
	"bytes"
	"path/filepath"
	"regexp"
	"strings"

	"boread/pkg/config"
)

// MetaExtractor 元数据提取器
type MetaExtractor struct {
	filenameRules []*compiledMetaRule
	contentRules  []*compiledMetaRule
}

type compiledMetaRule struct {
	re          *regexp.Regexp
	titleGroup  int
	authorGroup int
}

// DefaultMetaRules 内置默认规则
func DefaultMetaRules() []config.MetaExtractRule {
	return []config.MetaExtractRule{
		{
			Name:        "书名号+作者标记",
			Pattern:     `《(?P<title>[^》]+)》.*?(?:作者[：:]\s*(?P<author>[^\s\[\]【】]+))`,
			TitleGroup:  "title",
			AuthorGroup: "author",
			Source:      "filename",
			Priority:    0,
		},
		{
			Name:        "书名号+后缀全名",
			Pattern:     `《(?P<title>[^》]+)》(?P<author>.+)`,
			TitleGroup:  "title",
			AuthorGroup: "author",
			Source:      "filename",
			Priority:    1,
		},
		{
			Name:        "下划线书名_作者",
			Pattern:     `^(?P<title>[^_]+)_(?P<author>.+)$`,
			TitleGroup:  "title",
			AuthorGroup: "author",
			Source:      "filename",
			Priority:    2,
		},
		{
			Name:       "文件名全文作为书名",
			Pattern:    `(?P<title>.+)`,
			TitleGroup: "title",
			Source:     "filename",
			Priority:   3,
		},
		{
			Name:        "内容-书名标记",
			Pattern:     `(?:^|[^\S\n])(?:书名[：:])\s*(?P<title>[^\n]{1,60})`,
			TitleGroup:  "title",
			AuthorGroup: "author",
			Source:      "content",
			Priority:    0,
		},
		{
			Name:        "内容-作者标记",
			Pattern:     `(?:^|[^\S\n])(?:作者[：:])\s*(?P<author>[^\n]{1,60})`,
			AuthorGroup: "author",
			Source:      "content",
			Priority:    1,
		},
	}
}

// NewMetaExtractor 创建提取器
// 优先使用配置中的规则，没有配置则使用默认规则
func NewMetaExtractor(rules ...config.MetaExtractRule) *MetaExtractor {
	if len(rules) == 0 {
		rules = DefaultMetaRules()
		if config.Cfg != nil && len(config.Cfg.Meta.Rules) > 0 {
			rules = config.Cfg.Meta.Rules
		}
	}

	ex := &MetaExtractor{}
	for _, r := range rules {
		re, err := regexp.Compile(r.Pattern)
		if err != nil {
			continue
		}
		cr := &compiledMetaRule{re: re}
		names := re.SubexpNames()
		for i, name := range names {
			if name == r.TitleGroup {
				cr.titleGroup = i
			}
			if name == r.AuthorGroup {
				cr.authorGroup = i
			}
		}
		switch r.Source {
		case "content":
			ex.contentRules = append(ex.contentRules, cr)
		default:
			ex.filenameRules = append(ex.filenameRules, cr)
		}
	}
	return ex
}

// Extract 从文件名和内容提取书名和作者
// 策略：先按优先级匹配文件名规则 → 匹配到书名后再匹配内容规则补充作者 → 匹配不到书名再从内容规则提取
func (ex *MetaExtractor) Extract(data []byte, filename string) (title, author string) {
	data = tryDecodeToUTF8(data)
	base := strings.TrimSuffix(filename, filepath.Ext(filename))

	// 1. 文件名规则提取
	for _, r := range ex.filenameRules {
		matches := r.re.FindStringSubmatch(base)
		if matches == nil {
			continue
		}
		if r.titleGroup > 0 && r.titleGroup < len(matches) {
			if v := strings.TrimSpace(matches[r.titleGroup]); v != "" {
				title = v
			}
		}
		if r.authorGroup > 0 && r.authorGroup < len(matches) {
			if v := strings.TrimSpace(matches[r.authorGroup]); v != "" {
				author = v
			}
		}
		if title != "" {
			break
		}
	}

	// 2. 内容规则补充作者（如果文件名没提取到作者）
	if author == "" {
		author = ex.extractFromContent(data, "author")
	}
	// 3. 内容规则补充书名（如果文件名没提取到书名）
	if title == "" {
		title = ex.extractFromContent(data, "title")
	}

	if title == "" {
		title = "未命名"
	}
	return
}

// extractFromContent 从内容前 N 行提取指定字段
func (ex *MetaExtractor) extractFromContent(data []byte, field string) string {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Buffer(make([]byte, 0, 64*1024), 64*1024)
	lineNum := 0
	var contentBuf strings.Builder
	for scanner.Scan() {
		line := scanner.Text()
		contentBuf.WriteString(line)
		contentBuf.WriteByte('\n')
		lineNum++
		if lineNum > 10 {
			break
		}
	}
	head := contentBuf.String()

	for _, r := range ex.contentRules {
		isTitle := r.titleGroup > 0 && field == "title"
		isAuthor := r.authorGroup > 0 && field == "author"
		if !isTitle && !isAuthor {
			continue
		}
		matches := r.re.FindStringSubmatch(head)
		if matches == nil {
			continue
		}
		idx := r.titleGroup
		if field == "author" {
			idx = r.authorGroup
		}
		if idx > 0 && idx < len(matches) {
			if v := strings.TrimSpace(matches[idx]); v != "" {
				return v
			}
		}
	}
	return ""
}

// extractMetaFromContent 兼容旧接口，使用默认提取器
func extractMetaFromContent(data []byte, filename string) (title, author string) {
	return defaultExtractor.Extract(data, filename)
}

// defaultExtractor 包级默认提取器
var defaultExtractor = NewMetaExtractor()
