package controllers

import (
	"blog_holdongametime/app/models"

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
	var p []models.Post
	var author = "Me"
	result := c.DB.Order("publish_date desc").Find(&p)
	if result.Error != nil {
		panic(result.Error)
	}
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
