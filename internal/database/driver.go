// Package database 提供统一数据库驱动接口及注册机制。
package database

import (
	"database/sql"
	"fmt"
	"sync"
)

// Driver 抽象了不同数据库的通用操作。
type Driver interface {
	// Open 根据 DSN 建立数据库连接。
	Open(dsn string) (*sql.DB, error)
	// Close 释放数据库连接资源。
	Close(db *sql.DB) error
	// Query 执行查询并返回统一结果集。
	Query(db *sql.DB, query string) (*ResultSet, error)
	// Exec 执行非查询语句（INSERT/UPDATE/DELETE 等）。
	Exec(db *sql.DB, query string) (sql.Result, error)
	// Tables 获取数据库中所有表名。
	Tables(db *sql.DB) ([]string, error)
}

// DriverFactory 是创建 Driver 实例的构造函数类型。
type DriverFactory func() Driver

var (
	driversMu sync.RWMutex
	drivers   = make(map[string]DriverFactory)
)

// Register 注册数据库驱动，供初始化时调用。
func Register(name string, factory DriverFactory) {
	driversMu.Lock()
	defer driversMu.Unlock()
	if _, ok := drivers[name]; ok {
		panic("database driver already registered: " + name)
	}
	drivers[name] = factory
}

// Open 根据驱动名和 DSN 打开数据库连接，返回 *sql.DB 和 Driver 实例。
func Open(driverName, dsn string) (*sql.DB, Driver, error) {
	driversMu.RLock()
	factory, ok := drivers[driverName]
	driversMu.RUnlock()
	if !ok {
		return nil, nil, fmt.Errorf("未注册的数据库驱动: %s", driverName)
	}
	d := factory()
	db, err := d.Open(dsn)
	if err != nil {
		return nil, nil, err
	}
	return db, d, nil
}
