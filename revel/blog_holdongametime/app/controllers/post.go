package controllers

import (
	"github.com/revel/revel"
)

type Post struct {
	GormController
}

const author = "Me"

func (c Post) Show(slug string) revel.Result {

	p := c.GetPost(slug)

	rev := c.BuildReviewSkeleton(slug)
	return c.Render(rev, p, author)
}
