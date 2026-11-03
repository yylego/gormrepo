package gormclass_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yylego/gormrepo/gormclass"
	"github.com/yylego/gormrepo/gormclass/internal/models"
)

func TestUsePlain(t *testing.T) {
	one, cls := gormclass.UsePlain(&models.Account{})
	require.Equal(t, "accounts", one.TableName())
	// Plain columns come without the table name
	// 普通列不带表名
	require.Equal(t, "username", cls.Username.Name())
}

func TestUmcPlain(t *testing.T) {
	one, cls := gormclass.UmcPlain(&models.Example{})
	require.Equal(t, "examples", one.TableName())
	require.Equal(t, "name", cls.Name.Name())
}
