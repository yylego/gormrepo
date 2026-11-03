package gormclass_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yylego/gormrepo/gormclass"
	"github.com/yylego/gormrepo/gormclass/internal/models"
)

func TestUseTable(t *testing.T) {
	one, tableName, cls := gormclass.UseTable(&models.Account{})
	require.Equal(t, "accounts", tableName)
	require.Equal(t, one.TableName(), tableName)
	// Table columns come with the table name, which a join needs
	// 表列带表名，联表需要这种列
	require.Equal(t, "accounts.username", cls.Username.Name())
}

func TestUmcTable(t *testing.T) {
	one, tableName, cls := gormclass.UmcTable(&models.Example{})
	require.Equal(t, "examples", tableName)
	require.Equal(t, one.TableName(), tableName)
	require.Equal(t, "examples.name", cls.Name.Name())
}
