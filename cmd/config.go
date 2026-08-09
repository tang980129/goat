package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tang980129/goat/internal/config"
)

// configCmd 代表 config 命令组。
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "管理数据库连接配置",
	Long:  "添加、查看、修改或删除数据库连接配置。",
}

// addCmd 代表 config add 命令。
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "新增连接配置",
	Long:  "通过交互式引导添加一条数据库连接配置。",
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := config.NewStore(cfgFile)
		if err != nil {
			return err
		}
		svc := config.NewService(store)
		return svc.Add()
	},
}

// listCmd 代表 config list 命令。
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "列出所有配置",
	Long:  "以表格形式展示所有已保存的连接配置，默认连接标记★。",
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := config.NewStore(cfgFile)
		if err != nil {
			return err
		}
		svc := config.NewService(store)
		return svc.List()
	},
}

// removeCmd 代表 config remove 命令。
var removeCmd = &cobra.Command{
	Use:   "remove <别名>",
	Short: "删除指定配置",
	Long:  "根据别名删除对应的连接配置。",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := config.NewStore(cfgFile)
		if err != nil {
			return err
		}
		svc := config.NewService(store)
		return svc.Remove(args[0])
	},
}

// editCmd 代表 config edit 命令。
var editCmd = &cobra.Command{
	Use:   "edit <别名>",
	Short: "修改已有配置",
	Long:  "通过交互式方式修改指定别名的连接配置，可直接回车保留原值。",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := config.NewStore(cfgFile)
		if err != nil {
			return err
		}
		svc := config.NewService(store)
		return svc.Edit(args[0])
	},
}

func init() {
	// 将子命令注册到 config 命令
	configCmd.AddCommand(addCmd)
	configCmd.AddCommand(listCmd)
	configCmd.AddCommand(removeCmd)
	configCmd.AddCommand(editCmd)

	// 将 config 注册到根命令
	rootCmd.AddCommand(configCmd)
}
