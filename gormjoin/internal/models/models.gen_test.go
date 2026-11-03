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

	// Configure generation options, TableColumns gives the columns with the table name used in joins
	options := gormcngen.NewOptions()

	// Define models to process
	models := []any{
		&Student{},
		&Score{},
		&Course{},
	}

	// Create config and generate
	cfg := gormcngen.NewConfigs(models, options, absPath)
	cfg.WithGoBuildDirective(gormcngen.DirectiveSkipGenerate)
	cfg.Gen()
}
