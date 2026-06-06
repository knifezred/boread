package v1

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"boread/internal/code"
	"boread/internal/model"
	"boread/pkg/appsignal"
	"boread/pkg/config"
	"boread/pkg/response"
)

// SetupHandler 系统初始化配置
type SetupHandler struct{}

func NewSetupHandler() *SetupHandler {
	return &SetupHandler{}
}

type setupStatusResponse struct {
	Configured bool `json:"configured"`
}

type SaveDatabaseConfigRequest struct {
	Host     string `json:"host" binding:"required"`
	Port     int    `json:"port" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	DBName   string `json:"dbname" binding:"required"`
}

// Status 返回系统是否已配置数据库
// @Summary   系统初始化状态
// @Tags      setup
// @Produce   json
// @Success  200  {object}  response.Response{data=setupStatusResponse}
// @Router   /api/setup/status [get]
func (h *SetupHandler) Status(c *gin.Context) {
	cfg := config.Cfg
	configured := cfg != nil && cfg.Database.Host != ""
	response.Success(c, &setupStatusResponse{Configured: configured})
}

// SaveConfig 验证并保存数据库配置
// @Summary   保存数据库配置
// @Tags      setup
// @Accept    json
// @Produce   json
// @Param    body  body  SaveDatabaseConfigRequest  true  "数据库连接信息"
// @Success  200   {object}  response.Response
// @Router   /api/setup/database [post]
func (h *SetupHandler) SaveConfig(c *gin.Context) {
	var req SaveDatabaseConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, code.ParamInvalid, "请求参数无效: "+err.Error())
		return
	}

	// 尝试连接测试
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		req.Username, req.Password, req.Host, req.Port, req.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		response.Error(c, code.DBNotConfigured, "数据库连接失败: "+err.Error())
		return
	}

	sqlDB, err := db.DB()
	if err != nil {
		response.Error(c, code.DBNotConfigured, "数据库连接失败: "+err.Error())
		return
	}
	if err := sqlDB.Ping(); err != nil {
		sqlDB.Close()
		response.Error(c, code.DBNotConfigured, "数据库连接失败: "+err.Error())
		return
	}

	// 将数据库配置同步写入 sys_setting 表（upsert）
	setupDBSettings(db, req)
	sqlDB.Close()

	// 构建配置
	cfg := &config.Config{}
	if config.Cfg != nil {
		*cfg = *config.Cfg
	}
	cfg.Database.Host = req.Host
	cfg.Database.Port = req.Port
	cfg.Database.Username = req.Username
	cfg.Database.Password = req.Password
	cfg.Database.DBName = req.DBName
	cfg.Database.Driver = "mysql"
	if cfg.Database.MaxIdleConns == 0 {
		cfg.Database.MaxIdleConns = 10
	}
	if cfg.Database.MaxOpenConns == 0 {
		cfg.Database.MaxOpenConns = 100
	}
	if cfg.Server.Port == 0 {
		cfg.Server.Port = 8080
	}
	if cfg.Server.Mode == "" {
		cfg.Server.Mode = "debug"
	}
	if cfg.JWT.Secret == "" {
		cfg.JWT.Secret = "boread-secret"
	}
	if cfg.JWT.Expire == 0 {
		cfg.JWT.Expire = 7200
	}

	config.Cfg = cfg
	if err := config.Save("configs/config.yaml"); err != nil {
		response.Error(c, code.ServerError, "保存配置文件失败: "+err.Error())
		return
	}

	response.Success(c, nil)

	// 配置保存成功后触发服务自动重启
	appsignal.RequestRestart()
}

// setupDBSettings 将数据库连接配置 upsert 到 sys_setting 表
// category = "database"，keys: host / port / username / password / dbname
func setupDBSettings(db *gorm.DB, req SaveDatabaseConfigRequest) {
	items := []struct{ key, value, valueType string }{
		{"host", req.Host, "string"},
		{"port", strconv.Itoa(req.Port), "number"},
		{"username", req.Username, "string"},
		{"password", req.Password, "string"},
		{"dbname", req.DBName, "string"},
	}
	const category = "database"
	for _, item := range items {
		upsertSetupSetting(db, category, item.key, item.value, item.valueType)
	}
}

// upsertSetupSetting 按 category+key upsert sys_setting 记录
// 存在则更新 value/value_type，不存在则插入新行
func upsertSetupSetting(db *gorm.DB, category, key, value, valueType string) {
	var m model.SysSetting
	err := db.Where("category = ? AND `key` = ?", category, key).First(&m).Error
	if err == nil {
		// 已存在 → 更新
		db.Model(&m).Updates(map[string]any{
			"value":      value,
			"value_type": valueType,
		})
		return
	}
	// 不存在 → 创建
	desc := "setup 初始化写入"
	db.Create(&model.SysSetting{
		Category:    category,
		Key:         key,
		Value:       value,
		ValueType:   valueType,
		Description: &desc,
		Editable:    true,
		IsSystem:    true,
		Status:      model.StatusEnabled,
	})
}
