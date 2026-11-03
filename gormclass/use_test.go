package gormclass_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/yylego/done"
	"github.com/yylego/gormrepo/gormclass"
	"github.com/yylego/gormrepo/gormclass/internal/models"
	"github.com/yylego/rese"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestUse(t *testing.T) {
	dsn := fmt.Sprintf("file:db-%s?mode=memory&cache=shared", uuid.New().String())
	db := rese.P1(gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}))
	defer rese.F0(rese.P1(db.DB()).Close)

	done.Done(db.AutoMigrate(&models.Example{}))
	require.NoError(t, db.Create(&models.Example{Name: "aaa", Age: 1}).Error)
	require.NoError(t, db.Create(&models.Example{Name: "bbb", Age: 2}).Error)

	// Use gives the object and its columns at once, then the find uses both
	// Use 一次拿到模型和它的列，查询里两个都用上
	one, cls := gormclass.Use(&models.Example{})
	require.NoError(t, db.Where(cls.Name.Eq("bbb")).First(one).Error)
	require.Equal(t, 2, one.Age)
}

func TestUmc(t *testing.T) {
	account, cls := gormclass.Umc(&models.Account{})
	require.Equal(t, "accounts", account.TableName())
	require.Equal(t, "username", cls.Username.Name())
}

func TestCls(t *testing.T) {
	cls := gormclass.Cls(&models.Account{})
	require.Equal(t, "username", cls.Username.Name())
}

func TestUsc(t *testing.T) {
	examples, cls := gormclass.Usc(&models.Example{})
	require.NotNil(t, examples)
	require.Empty(t, examples)
	require.Equal(t, "name", cls.Name.Name())
}

func TestMsc(t *testing.T) {
	one, examples, cls := gormclass.Msc(&models.Example{})
	require.Equal(t, "examples", one.TableName())
	require.Empty(t, examples)
	require.Equal(t, "name", cls.Name.Name())
}

func TestNsc(t *testing.T) {
	one, examples, cls := gormclass.Nsc(&models.Example{}, 32)
	require.Equal(t, "examples", one.TableName())
	require.Empty(t, examples)
	require.Equal(t, 32, cap(examples))
	require.Equal(t, "name", cls.Name.Name())
}

func TestOne(t *testing.T) {
	one := gormclass.One(&models.Account{})
	require.Equal(t, reflect.Pointer, reflect.TypeOf(one).Kind())
}

func TestUms(t *testing.T) {
	examples := gormclass.Ums(&models.Example{})
	require.NotNil(t, examples) // a blank slice, not nil, so it encodes as [] in JSON
	require.Empty(t, examples)
}

func TestUss(t *testing.T) {
	examples := gormclass.Uss[*models.Example]()
	require.NotNil(t, examples)
	require.Empty(t, examples)
}

func TestUsn(t *testing.T) {
	examples := gormclass.Usn[*models.Example](100)
	require.Empty(t, examples)
	require.Equal(t, 100, cap(examples))
}
