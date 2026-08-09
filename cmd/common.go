package cmd

import "strings"

// isQuery 判断 SQL 语句是否为查询（返回结果集的语句）。
// 支持 SELECT、SHOW、DESC、DESCRIBE、EXPLAIN、WITH 等。
func isQuery(sql string) bool {
	t := strings.TrimSpace(sql)
	if len(t) < 4 {
		return false
	}
	prefix := strings.ToUpper(t[:6]) // 取前6个字符，足够覆盖常见关键字

	switch {
	case strings.HasPrefix(prefix, "SELECT"),
		strings.HasPrefix(prefix, "SHOW"),
		strings.HasPrefix(prefix, "DESC"),
		strings.HasPrefix(prefix, "EXPLAIN"),
		strings.HasPrefix(prefix, "WITH"):
		return true
	default:
		return false
	}
}
