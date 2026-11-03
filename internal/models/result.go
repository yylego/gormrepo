package models

// AccountInfo is a result row with no table, scanned via ScanList / ScanItem with typed aliases
// AccountInfo 是没有表的结果行，通过 ScanList 或 ScanItem 扫描，别名带类型检查
type AccountInfo struct {
	LoginName string
	ShowName  string
}
