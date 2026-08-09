# 🐐 goat

**GOAT (Greatest Of All Time) 通用数据库命令行工具**

基于 Go 语言构建，支持多数据库连接配置管理、交互式 SQL 终端及批量执行。

## 当前版本

v1.0.0 — 完整版发布，包含全部预期功能。

## 已实现命令

```bash
goat --help                          # 显示帮助信息
goat --version                       # 查看版本信息

# 配置管理
goat config add                      # 交互式添加连接配置
goat config list                     # 列出所有配置（表格，密码脱敏）
goat config remove <别名>            # 删除配置
goat config edit <别名>              # 修改已有配置

# 数据库连接与执行
goat connect <别名>                  # 进入交互式 SQL 终端
goat exec -e "SQL语句"               # 执行单条 SQL
goat exec -f <文件路径>              # 执行 SQL 脚本文件
goat exec -c <别名> -e "..."        # 指定连接执行

goat version                         # 显示版本、Commit 与构建时间
```

## 依赖

| 依赖                                                                                                                            | 用途     |
|-------------------------------------------------------------------------------------------------------------------------------|--------|
| [github.com/spf13/cobra](https://link.wtturl.cn/?target=https%3A%2F%2Fgithub.com%2Fspf13%2Fcobra&scene=im&aid=497858&lang=zh) | CLI 框架 |
| [github.com/spf13/viper](https://link.wtturl.cn/?target=https%3A%2F%2Fgithub.com%2Fspf13%2Fviper&scene=im&aid=497858&lang=zh) | 配置管理   |
| [github.com/fatih/color](https://link.wtturl.cn/?target=https%3A%2F%2Fgithub.com%2Ffatih%2Fcolor&scene=im&aid=497858&lang=zh) | 终端彩色输出 |
| [github.com/go-sql-driver/mysql](https://github.com/go-sql-driver/mysql) | MySQL 驱动 |

其他间接依赖由 Go Module 自动管理。

# 构建与安装

## 使用 Makefile（推荐）

```
make build              # 构建当前平台可执行文件
make build-windows      # 交叉编译 Windows .exe
make clean              # 清理
```

## 手动构建

### Windows PowerShell

```
$commit = if (Get-Command git -ErrorAction SilentlyContinue) { git rev-parse --short HEAD } else { "dev" }
$buildDate = (Get-Date).ToUniversalTime().ToString("yyyy-MM-ddTHH:mm:ssZ")
go build -ldflags "-X 'github.com/tang980129/goat/cmd.version=1.0.0' -X 'github.com/tang980129/goat/cmd.commit=$commit' -X 'github.com/tang980129/goat/cmd.buildDate=$buildDate'" -o goat.exe .
```

### macOS / Linux

```
go build -ldflags " \
  -X 'github.com/tang980129/goat/cmd.version=1.0.0' \
  -X 'github.com/tang980129/goat/cmd.commit=$(git rev-parse --short HEAD)' \
  -X 'github.com/tang980129/goat/cmd.buildDate=$(date -u +'%Y-%m-%dT%H:%M:%SZ')'" \
  -o goat .
```

# 配置文件

默认位置 `$HOME/.goat.yaml`，可通过 `--config` 全局标志覆盖。

配置示例：

```
configs:
  - alias: local
    driver: mysql
    host: 127.0.0.1
    port: 3306
    user: root
    password: cm9vdA==  # Base64 编码
    database: test
    default: true
```

# 项目结构

```
goat/
├── cmd/              # CLI 命令定义
├── internal/
│   ├── config/       # 配置管理模块
│   └── database/     # 数据库抽象层及 MySQL 驱动
├── main.go
├── Makefile
└── go.mod
```