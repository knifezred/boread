package model

// SysConfig 系统配置主表
type SysConfig struct {
	BaseModel
	ConfigName  string       `gorm:"column:config_name;size:128;not null" json:"configName"`
	ConfigKey   string       `gorm:"column:config_key;size:128;not null;uniqueIndex" json:"configKey"`
	ConfigValue string       `gorm:"column:config_value;type:text" json:"configValue"`
	ConfigType  string       `gorm:"column:config_type;size:32;not null;default:'string'" json:"configType"`
	ConfigGroup string       `gorm:"column:config_group;size:64;not null;default:'system'" json:"configGroup"`
	Description string       `gorm:"column:description;size:512" json:"description"`
	SortOrder   int          `gorm:"column:sort_order;not null;default:0" json:"sortOrder"`
	Status      EnableStatus `gorm:"column:status;type:char(1);not null;default:'1'" json:"status"`
}

func (SysConfig) TableName() string { return "sys_config" }

// ConfigRuleType 规则类型
type ConfigRuleType string

const (
	RuleTitleMatch      ConfigRuleType = "title_match"
	RuleMetadataExtract ConfigRuleType = "metadata_extract"
)

// MatchSource 匹配来源
type MatchSource string

const (
	MatchSourceFilename MatchSource = "filename"
	MatchSourceContent  MatchSource = "content"
	MatchSourceBoth     MatchSource = "both"
)

// PatternType 匹配方式
type PatternType string

const (
	PatternRegex     PatternType = "regex"
	PatternDelimiter PatternType = "delimiter"
)

// SysConfigRule 配置规则明细
type SysConfigRule struct {
	BaseModel
	ConfigID     uint64         `gorm:"column:config_id;not null;index" json:"configId"`
	RuleName     string         `gorm:"column:rule_name;size:128;not null" json:"ruleName"`
	RuleType     ConfigRuleType `gorm:"column:rule_type;size:32;not null" json:"ruleType"`
	MatchSource  MatchSource    `gorm:"column:match_source;size:32;not null;default:'filename'" json:"matchSource"`
	Pattern      string         `gorm:"column:pattern;size:1024;not null" json:"pattern"`
	PatternType  PatternType    `gorm:"column:pattern_type;size:16;not null;default:'regex'" json:"patternType"`
	Delimiter    string         `gorm:"column:delimiter;size:32" json:"delimiter"`
	Position     int            `gorm:"column:position;default:0" json:"position"`
	TitleGroup   string         `gorm:"column:title_group;size:32;default:'title'" json:"titleGroup"`
	AuthorGroup  string         `gorm:"column:author_group;size:32;default:'author'" json:"authorGroup"`
	FieldMapping JSONMap        `gorm:"column:field_mapping;type:json" json:"fieldMapping"`
	Priority     int            `gorm:"column:priority;not null;default:0" json:"priority"`
	Status       EnableStatus   `gorm:"column:status;type:char(1);not null;default:'1'" json:"status"`
	Description  string         `gorm:"column:description;size:512" json:"description"`
}

func (SysConfigRule) TableName() string { return "sys_config_rule" }

// ChangeType 变更类型
type ChangeType string

const (
	ChangeCreate ChangeType = "create"
	ChangeUpdate ChangeType = "update"
	ChangeDelete ChangeType = "delete"
	ChangeImport ChangeType = "import"
)

// SysConfigHistory 配置历史
type SysConfigHistory struct {
	ID           uint64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ConfigID     uint64     `gorm:"column:config_id;not null;index" json:"configId"`
	FieldName    string     `gorm:"column:field_name;size:64;not null;default:'value'" json:"fieldName"`
	OldValue     string     `gorm:"column:old_value;type:text" json:"oldValue"`
	NewValue     string     `gorm:"column:new_value;type:text" json:"newValue"`
	ChangeType   ChangeType `gorm:"column:change_type;size:16;not null;default:'update'" json:"changeType"`
	ChangeDesc   string     `gorm:"column:change_desc;size:512" json:"changeDesc"`
	Operator     uint64     `gorm:"column:operator" json:"operator"`
	OperatorName string     `gorm:"column:operator_name;size:64" json:"operatorName"`
	CreateTime   string     `gorm:"column:create_time" json:"createTime"`
}

func (SysConfigHistory) TableName() string { return "sys_config_history" }
