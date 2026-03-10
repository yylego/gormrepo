package example14

import "github.com/yylego/gormrepo/internal/examples/example14/internal/models"

type Guest struct {
	models.Guest
	Orders []*models.Order
}
