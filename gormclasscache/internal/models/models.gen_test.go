package models

import (
	"testing"

	"github.com/yylego/gormcngen"
	"github.com/yylego/osexistpath/osmustexist"
	"github.com/yylego/runpath/runtestpath"
)

//go:generate go test -run ^TestGenerateModels$ -tags=gormcngen_generate
func TestGenerateModels(t *testing.T) {
	// Get absolute path to target file (models.gen.go)
	absPath := osmustexist.FILE(runtestpath.SrcPath(t))
	t.Log(absPath)

	// Configure generation options, without embedded operations to keep the columns struct plain
	options := gormcngen.NewOptions().
		WithEmbedColumnOperations(false)

	// Define models to process
	models := []any{
		&Account{},
	}

	// Create config and generate
	cfg := gormcngen.NewConfigs(models, options, absPath)
	cfg.WithGoBuildDirective(gormcngen.DirectiveSkipGenerate)
	cfg.Gen()
}
