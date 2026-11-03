package models

// Account is the test struct of the cache tests, each cache approach reads its columns
// Account 是缓存测试用的模型，每种缓存方式都读取它的列
type Account struct {
	ID       uint   // ID of the account // 账户 ID
	Username string // Login name of the account // 登录名
	Password string // Password hash of the account // 密码哈希
}

// TableName gives the table name, the cache uses it as the index of the columns
// TableName 返回表名，缓存用它作为键
func (*Account) TableName() string {
	return "accounts"
}
