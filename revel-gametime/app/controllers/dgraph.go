package controllers

import (
	"gametime/conf"
	"gametime/db"

	"github.com/revel/revel"
)

type DController struct {
	*revel.Controller
	dgraph *db.Dgraph
	cfg    *conf.Config
}

func NewDController(cfg *conf.Config) *DController {
	return &DController{
		cfg: cfg,
	}
}
