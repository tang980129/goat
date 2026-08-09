package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// 编译时通过 -ldflags 注入的版本信息。
var (
	version   = "dev"
	commit    = "none"
	buildDate = "unknown"
)

// 保存 --config 标志传入的配置文件路径，用于覆盖默认路径。
var cfgFile string

// goat 的根命令。
var rootCmd = &cobra.Command{
	Use:   "goat",
	Short: "GOAT - 通用数据库命令行工具",
	Long: color.CyanString(`🐐 GOAT (Greatest Of All Time) 通用数据库 CLI

    ██████╗  ██████╗  █████╗ ████████╗
   ██╔════╝ ██╔═══██╗██╔══██╗╚══██╔══╝
   ██║  ███╗██║   ██║███████║   ██║
   ██║   ██║██║   ██║██╔══██║   ██║
   ╚██████╔╝╚██████╔╝██║  ██║   ██║
    ╚═════╝  ╚═════╝ ╚═╝  ╚═╝   ╚═╝
`) + "\n" +
		color.YellowString("支持多数据库，提供配置管理、交互式终端与批量执行。"),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			if err := cmd.Help(); err != nil {
				if _, err := fmt.Fprintf(os.Stderr, "%v\n", err); err != nil {
					os.Exit(2)
				}
				os.Exit(1)
			}
		}
	},
	Version: version,
}

// Execute 注册子命令并运行根命令。
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "",
		"配置文件路径 (默认 $HOME/.goat.yaml)")

	if err := viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config")); err != nil {
		if _, err := fmt.Fprintf(os.Stderr, "绑定 config 标志失败: %v\n", err); err != nil {
			os.Exit(2)
		}
		os.Exit(1)
	}

	rootCmd.SetVersionTemplate(
		fmt.Sprintf(
			color.GreenString("🐐 goat version %s")+"\n"+
				color.HiBlackString("   commit: %s")+"\n"+
				color.HiBlackString("   built:  %s")+"\n",
			version, commit, buildDate,
		),
	)
}
