package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

// Service 提供配置管理的业务操作。
type Service struct {
	store *Store
}

// NewService 创建配置服务实例。
func NewService(store *Store) *Service {
	return &Service{store: store}
}

// Add 交互式添加一条新配置。
func (s *Service) Add() error {
	entry := ConfigEntry{}
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("别名: ")
	entry.Alias, _ = reader.ReadString('\n')
	entry.Alias = strings.TrimSpace(entry.Alias)
	if entry.Alias == "" {
		return fmt.Errorf("别名不能为空")
	}

	fmt.Print("驱动 (默认 mysql): ")
	entry.Driver, _ = reader.ReadString('\n')
	entry.Driver = strings.TrimSpace(entry.Driver)
	if entry.Driver == "" {
		entry.Driver = "mysql"
	}

	fmt.Print("主机 (默认 127.0.0.1): ")
	entry.Host, _ = reader.ReadString('\n')
	entry.Host = strings.TrimSpace(entry.Host)
	if entry.Host == "" {
		entry.Host = "127.0.0.1"
	}

	fmt.Print("端口 (默认 3306): ")
	portStr, _ := reader.ReadString('\n')
	portStr = strings.TrimSpace(portStr)
	if portStr == "" {
		portStr = "3306"
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return fmt.Errorf("端口无效: %w", err)
	}
	entry.Port = port

	fmt.Print("用户名: ")
	entry.User, _ = reader.ReadString('\n')
	entry.User = strings.TrimSpace(entry.User)

	fmt.Print("密码: ")
	password, _ := reader.ReadString('\n')
	entry.Password = EncodePassword(strings.TrimSpace(password))

	fmt.Print("数据库名: ")
	entry.Database, _ = reader.ReadString('\n')
	entry.Database = strings.TrimSpace(entry.Database)

	fmt.Print("设为默认连接? (y/N): ")
	defaultStr, _ := reader.ReadString('\n')
	entry.Default = strings.TrimSpace(strings.ToLower(defaultStr)) == "y"

	configs, err := s.store.Load()
	if err != nil {
		return err
	}

	// 若已存在同名别名则拒绝添加
	for _, c := range configs {
		if c.Alias == entry.Alias {
			return fmt.Errorf("别名 %q 已存在", entry.Alias)
		}
	}

	// 若设为默认，取消其他默认标记
	if entry.Default {
		for i := range configs {
			configs[i].Default = false
		}
	}

	configs = append(configs, entry)
	if err := s.store.Save(configs); err != nil {
		return err
	}

	color.Green("✓ 配置 %q 添加成功", entry.Alias)
	return nil
}

// List 列出所有保存的配置，表格形式，密码脱敏。
func (s *Service) List() error {
	configs, err := s.store.Load()
	if err != nil {
		return err
	}
	if len(configs) == 0 {
		color.Yellow("暂无保存的配置")
		return nil
	}

	// 表头
	color.Cyan("%-15s %-8s %-22s %-10s %-10s %s", "别名", "驱动", "主机:端口", "用户", "数据库", "默认")
	for _, c := range configs {
		defaultMark := ""
		if c.Default {
			defaultMark = color.GreenString("★")
		}
		// 密码永远显示为 ****
		color.White("%-15s %-8s %-22s %-10s %-10s %s",
			c.Alias, c.Driver, fmt.Sprintf("%s:%d", c.Host, c.Port),
			c.User, c.Database, defaultMark)
	}
	return nil
}

// Remove 删除指定别名的配置。
func (s *Service) Remove(alias string) error {
	configs, err := s.store.Load()
	if err != nil {
		return err
	}

	found := false
	newConfigs := make([]ConfigEntry, 0, len(configs)-1)
	for _, c := range configs {
		if c.Alias == alias {
			found = true
			continue
		}
		newConfigs = append(newConfigs, c)
	}
	if !found {
		return fmt.Errorf("别名 %q 不存在", alias)
	}
	if err := s.store.Save(newConfigs); err != nil {
		return err
	}
	color.Green("✓ 配置 %q 已删除", alias)
	return nil
}

// Edit 交互式修改已有配置（可仅修改部分字段）。
func (s *Service) Edit(alias string) error {
	configs, err := s.store.Load()
	if err != nil {
		return err
	}

	var target *ConfigEntry
	for i := range configs {
		if configs[i].Alias == alias {
			target = &configs[i]
			break
		}
	}
	if target == nil {
		return fmt.Errorf("别名 %q 不存在", alias)
	}

	reader := bufio.NewReader(os.Stdin)
	applyChange := func(prompt string, current string, setter func(string)) {
		fmt.Printf("%s [%s]: ", prompt, current)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input != "" {
			setter(input)
		}
	}

	applyChange("主机", target.Host, func(v string) { target.Host = v })
	applyChange(fmt.Sprintf("端口"), strconv.Itoa(target.Port), func(v string) {
		p, err := strconv.Atoi(v)
		if err == nil {
			target.Port = p
		}
	})
	applyChange("用户名", target.User, func(v string) { target.User = v })
	fmt.Printf("密码 [****]: ")
	if pw, _ := reader.ReadString('\n'); strings.TrimSpace(pw) != "" {
		target.Password = EncodePassword(strings.TrimSpace(pw))
	}
	applyChange("数据库", target.Database, func(v string) { target.Database = v })
	fmt.Printf("设为默认? (y/N) [%v]: ", target.Default)
	defStr, _ := reader.ReadString('\n')
	if strings.TrimSpace(strings.ToLower(defStr)) == "y" {
		// 取消其他默认
		for i := range configs {
			configs[i].Default = false
		}
		target.Default = true
	} else if strings.TrimSpace(defStr) != "" {
		target.Default = false
	}

	if err := s.store.Save(configs); err != nil {
		return err
	}
	color.Green("✓ 配置 %q 已更新", alias)
	return nil
}
