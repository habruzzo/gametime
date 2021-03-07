package controllers

import (
	"blog_holdongametime/app/models"

	"github.com/revel/revel"
)

type Post struct {
	GormController
}

const author = "Me"

func (c Post) Show(slug string) revel.Result {

	p, pathString := c.GetPostAndFile(slug)

	rev := models.NewReviewSkeleton(pathString)
	return c.Render(rev, p, author)
}
