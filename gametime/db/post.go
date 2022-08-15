package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gametime"
	"github.com/dgraph-io/dgo/v200/protos/api"
	"time"
)

func (d *Dgraph) dPostToPost(dp Dpost) gametime.Post {
	return gametime.Post{
		Id:       dp.Uid,
		ReviewId: dp.ReviewId,
		Date:     dp.Date,
	}
}

func (d *Dgraph) postToDpost(p gametime.Post) Dpost {
	return Dpost{
		Uid:        p.Id,
		ReviewId:   p.ReviewId,
		Date:       p.Date,
		DgraphType: postType,
	}
}

func (d *Dgraph) GetPostByReview(review gametime.Review) (gametime.Post, error) {
	d.log.Info("get post by review (app)")
	dr := d.reviewToDreview(review)
	dp, err := d.getPostByReview(dr)
	if err != nil {
		return gametime.Post{ReviewId: review.Id}, err
	}
	return gametime.Post{
		Id:       dp.Uid,
		Date:     dp.Date,
		ReviewId: review.Id,
	}, nil
}

func (d *Dgraph) getPostByReview(review Dreview) (Dpost, error) {
	d.newClient(d.cfg.Dgraph.Url)
	d.log.Info("get post by review")
	queryFmt := `{postByReview(func: eq(review, %s)) {
		uid
		review 
		postDate
		dgraph.type
	}}`
	query := fmt.Sprintf(queryFmt, review.Uid)
	resp, err := d.Dgraph.NewReadOnlyTxn().Query(context.Background(), query)
	if err != nil {
		d.log.Error("error getting post")
		return Dpost{ReviewId: review.Uid, DgraphType: postType}, err
	}
	type posts struct {
		Dposts []Dpost `json:"postByReview"`
	}
	var p posts
	err = json.Unmarshal(resp.Json, &p)
	if err != nil {
		d.log.Error("error unmarshalling post", "error", err)
		return Dpost{ReviewId: review.Uid, DgraphType: postType}, err
	}
	if len(p.Dposts) < 1 {
		d.log.Error("error getting posts from db")
		return Dpost{ReviewId: review.Uid, DgraphType: postType}, errors.New("post not in db")
	}
	if p.Dposts[0].ReviewId == "" {
		d.log.Error("no review id on the post that was just retrieved")
		p.Dposts[0].ReviewId = review.Uid
	}
	p.Dposts[0].DgraphType = postType
	return p.Dposts[0], nil
}

func (d *Dgraph) GetPostsByMostRecent() ([]gametime.Post, error) {
	d.log.Info("get posts by recent (app)")
	dPosts, err := d.getPostsByMostRecent()
	if err != nil {
		d.log.Error("error getting posts recent", "error", err)
		return []gametime.Post{}, err
	}
	var posts []gametime.Post
	for _, v := range dPosts {
		posts = append(posts, d.dPostToPost(v))
	}
	return posts, nil
}

func (d *Dgraph) getPostsByMostRecent() ([]Dpost, error) {
	d.newClient(d.cfg.Dgraph.Url)
	d.log.Info("get posts by recent")
	query := `{postByRecent(func: type("Post"),orderdesc: postDate) {
		review
		uid
		postDate
	}}`
	resp, err := d.Dgraph.NewReadOnlyTxn().Query(context.Background(), query)
	if err != nil {
		d.log.Error("error getting post")
		return []Dpost{}, err
	}
	type posts struct {
		Dposts []Dpost `json:"postByRecent"`
	}
	var p posts
	err = json.Unmarshal(resp.Json, &p)
	if err != nil {
		d.log.Error("error unmarshalling posts", "error", err)
		return []Dpost{}, err
	}
	if len(p.Dposts) < 1 {
		d.log.Error("error getting posts from")
		return []Dpost{}, errors.New("posts not in db")
	}
	return p.Dposts, nil
}

// reviews and posts are 1:1
func (d *Dgraph) InsertPost(review gametime.Review) error {
	d.log.Info("insert post (app)", "name", review.Author.Name)

	//check if author already exists
	da := d.getAuthorByName(review.Author.Name)
	// if it doesnt exist, create author
	if da.Uid == "" {
		d.log.Error("author doesnt exist yet, creating", "author", da.AuthorName)
		if err := d.insertAuthor(da); err != nil {
			d.log.Error("couldnt create author", "author", da.AuthorName)
			return err
		}
		da = d.getAuthorByName(review.Author.Name)
	}
	//check if game already exists
	dg := d.getGameByTitle(review.Game.Title)
	// if it doesnt exist, create game
	if dg.Uid == "" {
		d.log.Error("game doesnt exist yet, creating", "game", dg.GameTitle)
		if err := d.insertGame(dg); err != nil {
			d.log.Error("couldnt create game", "game", dg.GameTitle)
			return err
		}
		dg = d.getGameByTitle(review.Game.Title)
	}
	//check if review already exists
	d.log.Info("does review exist?", "slug", review.Slug)
	dr, err := d.getReviewBySlug(review.Slug)
	if err != nil {
		d.log.Error("review doesnt exist yet, creating", "slug", dr.Slug)
		dr.Author = da
		dr.Game = dg
		dr.ReviewText = review.Text
		if err := d.insertReview(dr); err != nil {
			d.log.Error("couldnt create review", "slug", dr.Slug)
			return err
		}
		dr, err = d.getReviewBySlug(review.Slug)
		if err != nil {
			d.log.Error("couldnt create review again", "slug", dr.Slug)
			return err
		}
	}
	dp, err := d.getPostByReview(dr)
	if err != nil {
		d.log.Error("error getting post by review")
		dp.Date = time.Now()
		b, err := json.Marshal(dp)
		if err != nil {
			return err
		}
		d.newClient(d.cfg.Dgraph.Url)
		txn := d.Dgraph.NewTxn()
		defer txn.Discard(context.Background())
		r, err := txn.Mutate(context.Background(), &api.Mutation{SetJson: b,
			CommitNow: true})
		if err != nil {
			return err
		}
		fmt.Println(r)
	}
	fmt.Println(dr.Game, dr.Author, dr.Slug, dp.Date, dp.DgraphType)
	fmt.Println(dg, da, dr.Uid, dp.Uid)
	return nil
}

func (d *Dgraph) GetPostBySlug(slug string) (gametime.Post, error) {
	d.log.Info("get post by slug (app)")
	dr, err := d.getReviewBySlug(slug)
	if err != nil {
		d.log.Error("error getting post by slug (reviewbyslug)", "error", err)
		return d.dPostToPost(Dpost{ReviewId: dr.Uid}), err
	}
	dp, err := d.getPostByReview(dr)
	if err != nil {
		d.log.Error("error getting post by slug (postbyreview)", "error", err)
		return d.dPostToPost(dp), err
	}
	return d.dPostToPost(dp), nil
}
