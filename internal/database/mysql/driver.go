// Package mysql 实现 MySQL 数据库的 Driver 接口。
package mysql

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql" // 导入 MySQL 驱动，注册到 database/sql
	"github.com/tang980129/goat/internal/database"
)

// Driver 是 MySQL 的 database.Driver 实现。
type Driver struct{}

// init 自动注册 MySQL 驱动到 database 包。
func init() {
	database.Register("mysql", func() database.Driver {
		return &Driver{}
	})
}

// Open 使用 MySQL 驱动打开连接。
func (d *Driver) Open(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("打开 MySQL 连接失败: %w", err)
	}
	if err = db.Ping(); err != nil {
		// 关闭无效连接，并将关闭错误作为诊断信息输出
		if closeErr := db.Close(); closeErr != nil {
			if _, ferr := fmt.Fprintf(os.Stderr, "关闭无效连接时出错: %v\n", closeErr); ferr != nil {
				// 如果连 stderr 都写不了，放弃诊断输出
			}
		}
		return nil, fmt.Errorf("MySQL 连接测试失败: %w", err)
	}
	return db, nil
}

// Close 关闭数据库连接。
func (d *Driver) Close(db *sql.DB) error {
	if err := db.Close(); err != nil {
		return fmt.Errorf("关闭数据库连接失败: %w", err)
	}
	return nil
}

// Query 执行查询并转换为统一结果集。
func (d *Driver) Query(db *sql.DB, query string) (*database.ResultSet, error) {
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("查询执行失败: %w", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			if _, ferr := fmt.Fprintf(os.Stderr, "关闭结果集时出错: %v\n", closeErr); ferr != nil {
				// stderr 写入失败，无法记录
			}
		}
	}()

	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("获取列信息失败: %w", err)
	}

	result := &database.ResultSet{
		Columns: columns,
		Rows:    make([][]interface{}, 0),
	}

	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, fmt.Errorf("扫描行数据失败: %w", err)
		}

		row := make([]interface{}, len(columns))
		for i, v := range values {
			switch val := v.(type) {
			case []byte:
				row[i] = string(val)
			default:
				row[i] = val
			}
		}
		result.Rows = append(result.Rows, row)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("行迭代错误: %w", err)
	}

	return result, nil
}

// Exec 执行非查询语句。
func (d *Driver) Exec(db *sql.DB, query string) (sql.Result, error) {
	res, err := db.Exec(query)
	if err != nil {
		return nil, fmt.Errorf("执行失败: %w", err)
	}
	return res, nil
}

// Tables 查询当前数据库中的所有表名。
func (d *Driver) Tables(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		return nil, fmt.Errorf("获取表列表失败: %w", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			if _, ferr := fmt.Fprintf(os.Stderr, "关闭结果集时出错: %v\n", closeErr); ferr != nil {
				// stderr 写入失败，无法记录
			}
		}
	}()

	var tables []string
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return nil, fmt.Errorf("扫描表名失败: %w", err)
		}
		tables = append(tables, table)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("行迭代错误: %w", err)
	}
	return tables, nil
}
