package controllers

import (
	"github.com/revel/revel"
)

type App struct {
	GormController
}

func (c App) Index() revel.Result {
	return c.Render()
}

func (c App) About() revel.Result {
	return c.Render()
}

func (c App) Posts() revel.Result {
	var author = "Me"
	p := c.GetPostList()
	return c.Render(p, author)
}

func (c App) Format() revel.Result {
	return c.Render()
}

func (c App) Backlog() revel.Result {
	return c.Render()
}

func (c App) Contact() revel.Result {
	return c.Render()
}
