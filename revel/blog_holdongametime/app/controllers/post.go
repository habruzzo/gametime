package controllers

import (
	"github.com/revel/revel"
)

type Post struct {
	*revel.Controller
}

func (c Post) Show(slug string) revel.Result {
	c.BuildBlogPost()
	return c.RenderTemplate("app/views/Blog/BlogPost.html")
}
