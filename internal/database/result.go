package database

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

// ResultSet 表示一个查询结果集，包含列名和行数据。
type ResultSet struct {
	Columns []string
	Rows    [][]interface{}
}

// Print 将结果集以表格形式打印到终端。
func (rs *ResultSet) Print() {
	if len(rs.Columns) == 0 {
		if _, err := fmt.Fprintln(os.Stdout, "空结果集"); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "输出失败: %v\n", err)
		}
		return
	}

	// 打印表头（青色加粗）
	for _, col := range rs.Columns {
		_, err := color.New(color.FgCyan, color.Bold).Printf("%-20s", col)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "打印表头失败: %v\n", err)
		}
	}
	if _, err := fmt.Fprintln(os.Stdout); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "输出换行失败: %v\n", err)
	}

	// 打印分隔线
	if _, err := fmt.Fprintln(os.Stdout, strings.Repeat("-", 20*len(rs.Columns))); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "输出分隔线失败: %v\n", err)
	}

	// 打印数据行
	for _, row := range rs.Rows {
		for _, val := range row {
			_, err := fmt.Printf("%-20v", val)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "打印数据失败: %v\n", err)
			}
		}
		if _, err := fmt.Fprintln(os.Stdout); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "输出换行失败: %v\n", err)
		}
	}
	if _, err := fmt.Fprintf(os.Stdout, "(%d 行)\n", len(rs.Rows)); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "输出行数失败: %v\n", err)
	}
}
