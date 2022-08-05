package controllers

import (
	"gametime/conf"
	"gametime/db"

	"github.com/revel/revel"
)

type Post struct {
	DController
}

type PostReviewBundle struct {
	Post   conf.Post
	Review conf.Review
}

func (a Post) RenderPostFromSlug(slug string) revel.Result {
	a.Log.Info("start render post", "slug", slug)
	a.dgraph = db.NewDgraph(a.Log, Config)
	review, err := a.dgraph.GetReviewBySlug(slug)
	if err != nil {
		a.Response.SetStatus(404)
		a.Log.Error("no review found", "error", err)
		return a.Result
	}
	post, err := a.dgraph.GetPostByReview(review)
	if err != nil {
		a.Response.SetStatus(404)
		a.Log.Error("no post found", "error", err)
		return a.Result
	}
	a.DController.Response.SetStatus(200)
	a.Log.Info("finish render post and review", "slug", slug)
	prb := PostReviewBundle{
		Post:   post,
		Review: review,
	}
	return a.Render(prb)
}

func (d Post) GetPostsByMostRecent() revel.Result {
	d.dgraph = db.NewDgraph(d.Log, Config)
	d.Log.Info("start render post recent")

	posts, err := d.dgraph.GetPostsByMostRecent()
	if len(posts) < 1 || err != nil {
		d.Response.SetStatus(500)
		d.Log.Error("no posts found", "error", err)
		return d.Result
	}
	var prb []PostReviewBundle
	for _, v := range posts {
		prb = append(prb, PostReviewBundle{Post: v, Review: d.dgraph.GetReviewById(v.ReviewId)})
	}
	return d.Render(prb)
}
