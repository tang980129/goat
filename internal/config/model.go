package config

import (
	"fmt"
)

// ConfigEntry 描述一条数据库连接配置。
type ConfigEntry struct {
	Alias    string `mapstructure:"alias"`
	Driver   string `mapstructure:"driver"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"` // Base64 编码存储
	Database string `mapstructure:"database"`
	Default  bool   `mapstructure:"default"`
}

// DSN 根据配置生成数据库连接字符串。
// 自动解码 Base64 存储的密码，针对不同驱动拼接对应的 DSN。
func (c *ConfigEntry) DSN() (string, error) {
	plainPass, err := DecodePassword(c.Password)
	if err != nil {
		return "", fmt.Errorf("解码密码失败: %w", err)
	}
	switch c.Driver {
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true",
			c.User, plainPass, c.Host, c.Port, c.Database), nil
	default:
		return "", fmt.Errorf("不支持的驱动类型: %s", c.Driver)
	}
}
