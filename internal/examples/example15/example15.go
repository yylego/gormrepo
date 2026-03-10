package example15

import (
	"github.com/yylego/gormrepo/internal/examples/example15/internal/models"
)

type User struct {
	models.User
	Orders []*Order
}

type Order struct {
	models.Order
	Products []*models.Product
}
