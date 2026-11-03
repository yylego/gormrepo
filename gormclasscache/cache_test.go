// Package gormclasscache_test tests the three cache approaches UmcV1 / UmcV2 / UmcV3, each with its own cache backend
// The test struct and its generated columns are in internal/models
//
// gormclasscache_test 测试 UmcV1 / UmcV2 / UmcV3 三种缓存方式，各用各的缓存后端
// 测试模型和生成的列定义在 internal/models 里
package gormclasscache_test

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yylego/gormrepo/gormclasscache"
	"github.com/yylego/gormrepo/gormclasscache/internal/models"
	"github.com/yylego/mutexmap"
	"github.com/yylego/mutexmap/cachemap"
)

// TestUmcV1 tests UmcV1 with a cachemap backend: the second use gets the cached columns
// TestUmcV1 测试用 cachemap 做后端的 UmcV1：第二次调用拿到的是缓存里的列
func TestUmcV1(t *testing.T) {
	cache := cachemap.NewMap[string, any](1)

	one, cls := gormclasscache.UmcV1(&models.Account{}, cache)
	require.Equal(t, "accounts", one.TableName())
	require.Equal(t, "username", cls.Username.Name())

	_, cached := gormclasscache.UmcV1(&models.Account{}, cache)
	require.Same(t, cls, cached)
}

// TestUmcV2 tests UmcV2 with a mutexmap backend: the second use gets the cached columns
// TestUmcV2 测试用 mutexmap 做后端的 UmcV2：第二次调用拿到的是缓存里的列
func TestUmcV2(t *testing.T) {
	cache := mutexmap.NewMap[string, any](1)

	one, cls := gormclasscache.UmcV2(&models.Account{}, cache)
	require.Equal(t, "accounts", one.TableName())
	require.Equal(t, "username", cls.Username.Name())

	_, cached := gormclasscache.UmcV2(&models.Account{}, cache)
	require.Same(t, cls, cached)
}

// TestUmcV3 tests UmcV3 with a sync.Map backend: the second use gets the cached columns
// TestUmcV3 测试用 sync.Map 做后端的 UmcV3：第二次调用拿到的是缓存里的列
func TestUmcV3(t *testing.T) {
	cache := &sync.Map{}

	one, cls := gormclasscache.UmcV3(&models.Account{}, cache)
	require.Equal(t, "accounts", one.TableName())
	require.Equal(t, "username", cls.Username.Name())

	_, cached := gormclasscache.UmcV3(&models.Account{}, cache)
	require.Same(t, cls, cached)
}
