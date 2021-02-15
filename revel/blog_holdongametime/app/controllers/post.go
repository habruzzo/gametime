package controllers

import (
	"github.com/revel/revel"
)

type Post struct {
	*revel.Controller
}

func (c Post) Show(slug string) revel.Result {
	//c.BuildBlogPost()
	//jsonMapping := loadSlugs()
	//file := jsonMapping[slug]
	//r := models.NewReviewSkeleton(file)

	return c.Render()
}
