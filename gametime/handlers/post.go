package handlers

import (
	"gametime"
	"gametime/db"
	"net/http"
	"strings"
)

const (
	PostRenderPostFromSlug   = "Post.RenderPostFromSlug"
	PostGetPostsByMostRecent = "Post.GetPostsByMostRecent"

	RENDERPOST = "/posts/{slug}"
	GETPOSTS   = "/posts/"
)

type Post struct {
	*Handler
}

func NewPost(h *Handler) *Post {
	h.log.Info("CREATING POST", h.name)
	h.name = "post"
	p := &Post{
		Handler: h,
	}

	p.getHandlerMap()
	return p
}

type PostReviewBundle struct {
	Post   gametime.Post
	Review gametime.Review
}

func (a Post) RenderPostFromSlug() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := strings.Split(r.URL.Path, "/")[2]
		a.log.Info("start render post", "slug", slug)
		a.dgraph = db.NewDgraph(a.log, a.cfg)
		review, err := a.dgraph.GetReviewBySlug(slug)
		if err != nil {
			a.log.Error("no review found", "error", err)
			w.WriteHeader(404)
			return
		}
		pst, err := a.dgraph.GetPostByReview(review)
		if err != nil {
			a.log.Error("no post found", "error", err)
			w.WriteHeader(404)
			return
		}
		a.log.Info("finish render post and review", "slug", slug)
		prb := PostReviewBundle{
			Post:   pst,
			Review: review,
		}
		type Data struct {
			Prb PostReviewBundle
			headerTitles
		}
		data := Data{
			Prb: prb,
			headerTitles: headerTitles{
				Title:     slug,
				PageTitle: review.Game.Title,
			},
		}
		if a.tmap == nil {
			a.log.Error("error rendering post from slug view")
			w.WriteHeader(500)
			return
		}
		post := a.tmap[POST].Lookup("RenderPostFromSlug.html")
		if err := post.ExecuteTemplate(w, "RenderPostFromSlug.html", data); err != nil {
			a.log.Error(err)
		}
		return
	}
}

func (d Post) GetPostsByMostRecent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d.dgraph = db.NewDgraph(d.log, d.cfg)
		d.log.Info("start render post recent")

		posts, err := d.dgraph.GetPostsByMostRecent()
		if len(posts) < 1 || err != nil {
			d.log.Error("no posts found", "error", err)
			w.WriteHeader(404)
			return
		}
		var prb []PostReviewBundle
		for _, v := range posts {
			prb = append(prb, PostReviewBundle{Post: v, Review: d.dgraph.GetReviewById(v.ReviewId)})
		}
		if d.tmap == nil {
			d.log.Error("error rendering post from recent posts view")
			w.WriteHeader(500)
			return
		}
		type Data struct {
			Prb []PostReviewBundle
			headerTitles
		}
		data := Data{
			Prb: prb,
			headerTitles: headerTitles{
				Title:     "Recent Posts",
				PageTitle: "recent posts",
			},
		}
		recent := d.tmap[POST].Lookup("GetPostsByMostRecent.html")
		if err := recent.ExecuteTemplate(w, "GetPostsByMostRecent.html", data); err != nil {
			d.log.Error(err)
		}
		return
	}
}

func (c *Post) getHandlerMap() {
	c.hmap.fmap[PostRenderPostFromSlug] = HandlerFuncWrapper{
		method:      "GET",
		path:        RENDERPOST,
		handlerFunc: c.RenderPostFromSlug(),
	}
	c.hmap.fmap[PostGetPostsByMostRecent] = HandlerFuncWrapper{
		method:      "GET",
		path:        GETPOSTS,
		handlerFunc: c.GetPostsByMostRecent(),
	}
}
