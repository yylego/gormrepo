package gormclass_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yylego/gormrepo/gormclass"
	"github.com/yylego/neatjson/neatjsons"
)

func TestUseTableWithAccount(t *testing.T) {
	one, tableName, cls := gormclass.UseTable(&Account{})
	t.Log(one.TableName())
	t.Log(tableName)
	require.Equal(t, one.TableName(), tableName)
	t.Log(neatjsons.S(cls))
}

func TestUmcTableWithExample(t *testing.T) {
	one, tableName, cls := gormclass.UmcTable(&Example{})
	t.Log(one.TableName())
	t.Log(tableName)
	require.Equal(t, one.TableName(), tableName)
	t.Log(neatjsons.S(cls))
}
