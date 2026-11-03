package models

import (
	"testing"

	"github.com/yylego/gormcngen"
	"github.com/yylego/osexistpath/osmustexist"
	"github.com/yylego/runpath/runtestpath"
)

//go:generate go test -run ^TestGenerateResult$ -tags=gormcngen_generate
func TestGenerateResult(t *testing.T) {
	// Get absolute path to target file (result.gen.go)
	absPath := osmustexist.FILE(runtestpath.SrcPath(t))
	t.Log(absPath)

	// Configure generation options, the result rows have no table so TableColumns is not needed
	options := gormcngen.NewOptions().
		WithTableColumns(false)

	// Define result rows to process
	models := []any{
		&StudentClass{},
	}

	// Create config and generate
	cfg := gormcngen.NewConfigs(models, options, absPath)
	cfg.WithGoBuildDirective(gormcngen.DirectiveSkipGenerate)
	cfg.Gen()
}
