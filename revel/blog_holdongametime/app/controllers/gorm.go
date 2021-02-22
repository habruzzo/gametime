//controllers/gorm.go
package controllers

import (
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"

	"github.com/revel/revel"
)

// type: revel controller with `*gorm.DB`
type GormController struct {
	*revel.Controller
	DB *gorm.DB
}

// it can be used for jobs
var Gdb *gorm.DB

func (c *GormController) SetDB() revel.Result {
	c.DB = Gdb
	return nil
}
