# 🐐 goat

GOAT (Greatest Of All Time) 通用数据库命令行工具

基于 Go 语言构建，提供多数据库连接配置管理、交互式 SQL 终端及批量执行能力。

## 当前版本

v0.1.0 - 项目骨架搭建，支持版本展示与全局配置标志。

## 已实现命令

```
# 显示帮助信息（含 ASCII 山羊 Logo）
goat --help

# 查看版本、Commit 与构建时间
goat --version
```

## 依赖

表格

| 依赖 | 版本 | 用途 |
| --- | --- | --- |
| [github.com/spf13/cobra](https://link.wtturl.cn/?target=https%3A%2F%2Fgithub.com%2Fspf13%2Fcobra&scene=im&aid=497858&lang=zh) | v1.10.2 | CLI 框架 |
| [github.com/spf13/viper](https://link.wtturl.cn/?target=https%3A%2F%2Fgithub.com%2Fspf13%2Fviper&scene=im&aid=497858&lang=zh) | v1.21.0 | 配置管理 |
| [github.com/fatih/color](https://link.wtturl.cn/?target=https%3A%2F%2Fgithub.com%2Ffatih%2Fcolor&scene=im&aid=497858&lang=zh) | v1.19.0 | 终端彩色输出 |

其他间接依赖由 Go Module 自动管理。

## 构建与安装

### Windows (PowerShell)

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