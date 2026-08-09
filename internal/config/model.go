// Package config 提供数据库连接配置的数据模型与持久化管理。
package config

// ConfigEntry 表示一条数据库连接配置。
type ConfigEntry struct {
	Alias    string `mapstructure:"alias"`  // 唯一别名
	Driver   string `mapstructure:"driver"` // 数据库驱动（如 mysql）
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"` // Base64 编码存储
	Database string `mapstructure:"database"`
	Default  bool   `mapstructure:"default"` // 是否为默认连接
}
