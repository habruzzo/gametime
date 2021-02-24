package controllers

import (
	"blog_holdongametime/app/models"

	"github.com/revel/revel"
)

type Post struct {
	GormController
}

var author = "Me"

func (c Post) Show(slug string) revel.Result {
	//c.BuildBlogPost()
	//jsonMapping := loadSlugs()
	var p models.Post
	if err := c.DB.Where("slug=?", slug).Find(&p).Error; err != nil {
		panic(err)
	}
	rev := models.NewReviewSkeleton(p.ContentPath)
	return c.Render(rev, p, author)
}
