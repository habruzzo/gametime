package controllers

import (
	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	c.Log.Info("index")
	return c.Render()
}

func (c App) About() revel.Result {
	c.Log.Info("about")
	return c.Render()
}

func (c App) Format() revel.Result {
	c.Log.Info("format")
	return c.Render()
}

func (c App) Backlog() revel.Result {
	c.Log.Info("backlog")
	return c.Render()
}

func (c App) Contact() revel.Result {
	c.Log.Info("contact")
	return c.Render()
}
