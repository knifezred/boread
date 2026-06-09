package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	JWT      JWTConfig      `yaml:"jwt"`
	Log      LogConfig      `yaml:"log"`
	CORS     CORSConfig     `yaml:"cors"`
	Meta     MetaConfig     `yaml:"meta"`
}

type CORSConfig struct {
	AllowedOrigins []string `yaml:"allowed_origins"` // 为空则允许所有源 (*)
}

type MetaConfig struct {
	Rules []MetaExtractRule `yaml:"rules"`
}

type MetaExtractRule struct {
	Name        string `yaml:"name"`
	Pattern     string `yaml:"pattern"`
	TitleGroup  string `yaml:"titleGroup"`
	AuthorGroup string `yaml:"authorGroup"`
	Source      string `yaml:"source"`
	Priority    int    `yaml:"priority"`
}

type ServerConfig struct {
	Port int    `yaml:"port"`
	Mode string `yaml:"mode"`
}

type DatabaseConfig struct {
	Driver       string `yaml:"driver"`
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
	DBName       string `yaml:"dbname"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
	MaxOpenConns int    `yaml:"max_open_conns"`
}

type JWTConfig struct {
	Secret string `yaml:"secret"`
	Expire int    `yaml:"expire"`
}

type LogConfig struct {
	Level string `yaml:"level"`
	File  string `yaml:"file"`
}

var Cfg *Config

func Load(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	Cfg = &Config{}
	return yaml.Unmarshal(data, Cfg)
}

// Exists 检查配置文件是否存在
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// Save 将当前配置序列化写入文件，自动创建父目录
func Save(path string) error {
	data, err := yaml.Marshal(Cfg)
	if err != nil {
		return err
	}
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
