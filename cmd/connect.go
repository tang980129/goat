package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/tang980129/goat/internal/config"
	"github.com/tang980129/goat/internal/database"
	_ "github.com/tang980129/goat/internal/database/mysql" // 注册 MySQL 驱动
)

// connectCmd 是 connect 子命令，提供交互式 SQL 终端。
var connectCmd = &cobra.Command{
	Use:   "connect <别名>",
	Short: "连接数据库并进入交互式 SQL 终端",
	Long:  "根据配置别名连接数据库，进入交互式终端。输入 SQL 执行，'exit' 退出。",
	Args:  cobra.ExactArgs(1), // 必须恰好提供一个别名参数
	RunE:  runConnect,
}

func init() {
	rootCmd.AddCommand(connectCmd)
}

// runConnect 实现连接逻辑与交互循环。
func runConnect(cmd *cobra.Command, args []string) error {
	alias := args[0]

	// 加载配置文件
	store, err := config.NewStore(cfgFile)
	if err != nil {
		return fmt.Errorf("初始化配置存储失败: %w", err)
	}
	configs, err := store.Load()
	if err != nil {
		return fmt.Errorf("读取配置失败: %w", err)
	}

	// 查找对应别名配置
	var entry *config.ConfigEntry
	for i := range configs {
		if configs[i].Alias == alias {
			entry = &configs[i]
			break
		}
	}
	if entry == nil {
		return fmt.Errorf("别名 %q 不存在，请先用 'goat config add' 添加", alias)
	}

	// 构建 DSN 并打开连接
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

	// 连接成功提示
	if _, err := color.New(color.FgGreen).Printf("已连接到 %s (%s@%s:%d/%s)\n",
		alias, entry.User, entry.Host, entry.Port, entry.Database); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "输出连接信息失败: %v\n", err)
	}
	if _, err := fmt.Fprintln(os.Stdout, "输入 SQL 执行，'exit' 或 'quit' 退出。"); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "输出提示失败: %v\n", err)
	}

	// 交互式 REPL
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(color.CyanString("goat> "))
		if !scanner.Scan() {
			break // EOF
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		lowerLine := strings.ToLower(line)
		if lowerLine == "exit" || lowerLine == "quit" {
			color.Yellow("再见！")
			break
		}

		// 根据语句类型分发到 Query 或 Exec
		if isQuery(line) {
			rs, err := driver.Query(db, line)
			if err != nil {
				color.Red("查询错误: %v", err)
				continue
			}
			rs.Print()
		} else {
			res, err := driver.Exec(db, line)
			if err != nil {
				color.Red("执行错误: %v", err)
				continue
			}
			rows, _ := res.RowsAffected()
			color.Green("执行成功，影响 %d 行", rows)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("输入读取错误: %w", err)
	}
	return nil
}

// isQuery 简单判断 SQL 是否为查询语句（以 SELECT 开头）。
func isQuery(sql string) bool {
	t := strings.TrimSpace(sql)
	return len(t) >= 6 && strings.EqualFold(t[:6], "SELECT")
}
