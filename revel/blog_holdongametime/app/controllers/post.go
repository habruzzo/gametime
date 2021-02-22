package controllers

import (
	"blog_holdongametime/app/models"

	"github.com/revel/revel"
)

type Post struct {
	GormController
}

func (c Post) Show(slug string) revel.Result {
	//c.BuildBlogPost()
	//jsonMapping := loadSlugs()
	var p models.Post
	if err := c.DB.Create(&p).Error; err != nil {
		panic(err)
	}
	return c.RenderJSON(p)
	//err := c.rgorp.Txn.SelectOne(p, c.rgorp.Db.SqlStatementBuilder.Select("*").From("post").Where("slug=?", slug))
	//if err != nil {
	//	panic(err)
	//}
	//r := models.NewReviewSkeleton(p.ContentPath)
	//return c.Render(r, p)
}
