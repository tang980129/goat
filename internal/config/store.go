package config

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const defaultConfigName = ".goat"

// Store 负责配置文件的持久化读写。
type Store struct {
	vpr  *viper.Viper
	path string // 配置文件完整路径
}

// NewStore 初始化配置存储，若文件不存在则创建目录和空文件。
func NewStore(cfgPath string) (*Store, error) {
	v := viper.New()
	v.SetConfigType("yaml")

	if cfgPath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("获取用户目录失败: %w", err)
		}
		cfgPath = filepath.Join(home, defaultConfigName+".yaml")
	}

	v.SetConfigFile(cfgPath)
	if err := ensureConfigFile(cfgPath); err != nil {
		return nil, err
	}

	return &Store{vpr: v, path: cfgPath}, nil
}

// ensureConfigFile 确保配置文件所在目录及文件存在。
func ensureConfigFile(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建配置目录失败: %w", err)
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		f, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("创建配置文件失败: %w", err)
		}
		// 关闭文件并检查错误
		if err := f.Close(); err != nil {
			return fmt.Errorf("关闭配置文件失败: %w", err)
		}
	}
	return nil
}

// Load 从文件加载全部配置项。
func (s *Store) Load() ([]ConfigEntry, error) {
	if err := s.vpr.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			return []ConfigEntry{}, nil
		}
		return nil, fmt.Errorf("读取配置失败: %w", err)
	}

	var raw struct {
		Configs []ConfigEntry `mapstructure:"configs"`
	}
	if err := s.vpr.Unmarshal(&raw); err != nil {
		return nil, fmt.Errorf("解析配置失败: %w", err)
	}
	return raw.Configs, nil
}

// Save 将配置列表写入文件。
func (s *Store) Save(entries []ConfigEntry) error {
	s.vpr.Set("configs", entries)
	if err := s.vpr.WriteConfig(); err != nil {
		return fmt.Errorf("写入配置失败: %w", err)
	}
	return nil
}

// EncodePassword 对密码进行 Base64 编码。
func EncodePassword(plain string) string {
	return base64.StdEncoding.EncodeToString([]byte(plain))
}

// DecodePassword 对 Base64 编码的密码进行解码。
func DecodePassword(encoded string) (string, error) {
	b, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("密码解码失败: %w", err)
	}
	return string(b), nil
}
