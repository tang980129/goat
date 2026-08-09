package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd 是 root 的子命令，打印版本信息。
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "显示版本信息",
	Long:  "显示 goat 的版本号、Git 提交哈希与构建时间。",
	Run: func(cmd *cobra.Command, args []string) {
		// 直接使用根命令中设置的版本模板（已含颜色和格式化）
		fmt.Print(rootCmd.VersionTemplate())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
