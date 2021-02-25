package controllers

import "github.com/revel/revel"

func init() {
	revel.InterceptMethod((*GormController).SetDB, revel.BEFORE)
}
