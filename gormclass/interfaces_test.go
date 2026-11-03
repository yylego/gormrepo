package gormclass_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yylego/gormcnm"
	"github.com/yylego/gormrepo/gormclass"
	"github.com/yylego/gormrepo/gormclass/internal/models"
)

func TestColumnsInterface(t *testing.T) {
	var x gormclass.ColumnsInterface[*models.AccountColumns] = &models.Account{}
	c := x.Columns()
	require.Equal(t, "username", c.Username.Name())
	require.Equal(t, "nickname", c.Nickname.Name())
	require.Equal(t, "password", c.Password.Name())
}

func TestColumnsWithTableNameInterface(t *testing.T) {
	var x gormclass.ColumnsWithTableNameInterface[*models.AccountColumns] = &models.Account{}
	require.Equal(t, "accounts", x.TableName())
	c := x.Columns()
	require.Equal(t, "username", c.Username.Name())
}

func TestTableColumnsInterface(t *testing.T) {
	var x gormclass.TableColumnsInterface[*models.AccountColumns] = &models.Account{}
	c := x.TableColumns(gormcnm.NewTableDecoration("accounts"))
	require.Equal(t, "accounts.username", c.Username.Name())
}

func TestTableColumnsWithTableNameInterface(t *testing.T) {
	var x gormclass.TableColumnsWithTableNameInterface[*models.AccountColumns] = &models.Account{}
	require.Equal(t, "accounts", x.TableName())
	c := x.TableColumns(gormcnm.NewTableDecoration(x.TableName()))
	require.Equal(t, "accounts.username", c.Username.Name())
	require.Equal(t, "accounts.nickname", c.Nickname.Name())
}
