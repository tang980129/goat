package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/tang980129/goat/internal/config"
	"github.com/tang980129/goat/internal/database"
	_ "github.com/tang980129/goat/internal/database/mysql" // 注册 MySQL 驱动
)

var (
	execSQL  string // -e 传入的 SQL 语句
	execFile string // -f 传入的 SQL 文件路径
	execConn string // -c 指定的连接别名
)

// execCmd 是 root 命令的子命令，非交互式执行 SQL。
var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "非交互式执行 SQL 语句或脚本",
	Long: "通过 -e 执行一条 SQL 语句，或通过 -f 执行一个 SQL 脚本文件。\n" +
		"未指定连接时自动使用默认配置。",
	RunE: runExec,
}

func init() {
	execCmd.Flags().StringVarP(&execSQL, "execute", "e", "", "执行的 SQL 语句")
	execCmd.Flags().StringVarP(&execFile, "file", "f", "", "执行的 SQL 脚本文件")
	execCmd.Flags().StringVarP(&execConn, "connection", "c", "", "使用的连接配置别名（为空则使用默认）")
	execCmd.MarkFlagsMutuallyExclusive("execute", "file")

	rootCmd.AddCommand(execCmd)
}

func runExec(cmd *cobra.Command, args []string) error {
	// 加载配置
	store, err := config.NewStore(cfgFile)
	if err != nil {
		return fmt.Errorf("初始化配置存储失败: %w", err)
	}
	configs, err := store.Load()
	if err != nil {
		return fmt.Errorf("读取配置失败: %w", err)
	}

	var entry *config.ConfigEntry

	if execConn != "" {
		// 使用指定的别名
		for i := range configs {
			if configs[i].Alias == execConn {
				entry = &configs[i]
				break
			}
		}
		if entry == nil {
			return fmt.Errorf("别名 %q 不存在", execConn)
		}
	} else {
		// 查找默认配置
		for i := range configs {
			if configs[i].Default {
				entry = &configs[i]
				break
			}
		}
		if entry == nil {
			return fmt.Errorf("未指定连接且没有默认配置，请通过 -c 指定别名或设置默认配置")
		}
	}

	// 获取待执行的 SQL 文本
	var sqlText string
	if execSQL != "" {
		sqlText = execSQL
	} else if execFile != "" {
		data, err := os.ReadFile(execFile)
		if err != nil {
			return fmt.Errorf("读取 SQL 文件失败: %w", err)
		}
		sqlText = string(data)
	} else {
		return fmt.Errorf("必须指定 -e 或 -f")
	}

	sqlText = strings.TrimSpace(sqlText)
	if sqlText == "" {
		return fmt.Errorf("SQL 内容为空")
	}

	// 建立连接
	dsn, err := entry.DSN()
	if err != nil {
		return err
	}
	db, driver, err := database.Open(entry.Driver, dsn)
	if err != nil {
		return fmt.Errorf("连接数据库失败: %w", err)
	}
	defer func() {
		if err := driver.Close(db); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "关闭连接失败: %v\n", err)
		}
	}()

	// 执行并输出结果
	if isQuery(sqlText) {
		rs, err := driver.Query(db, sqlText)
		if err != nil {
			return fmt.Errorf("查询失败: %w", err)
		}
		rs.Print()
	} else {
		res, err := driver.Exec(db, sqlText)
		if err != nil {
			return fmt.Errorf("执行失败: %w", err)
		}
		rows, _ := res.RowsAffected()
		color.Green("执行成功，影响 %d 行", rows)
	}

	return nil
}
