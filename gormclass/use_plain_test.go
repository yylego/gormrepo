package gormclass_test

import (
	"testing"

	"github.com/yylego/gormrepo/gormclass"
	"github.com/yylego/neatjson/neatjsons"
)

func TestUsePlainWithAccount(t *testing.T) {
	one, cls := gormclass.UsePlain(&Account{})
	t.Log(one.TableName())
	t.Log(neatjsons.S(cls))
}

func TestUmcPlainWithExample(t *testing.T) {
	one, cls := gormclass.UmcPlain(&Example{})
	t.Log(one.TableName())
	t.Log(neatjsons.S(cls))
}
