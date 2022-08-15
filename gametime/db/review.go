package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gametime"
	"github.com/dgraph-io/dgo/v200/protos/api"
)

func (d *Dgraph) GetReviewBySlug(slug string) (gametime.Review, error) {
	d.log.Info("get review by slug (app)")
	da, err := d.getReviewBySlug(slug)
	return d.dReviewToReview(da), err
}

func (d *Dgraph) getReviewBySlug(slug string) (Dreview, error) {
	d.log.Info("dgraphurl", "url", d.cfg.Dgraph.Url)

	d.newClient(d.cfg.Dgraph.Url)
	d.log.Info("get review by slug")
	queryFmt := `{reviewBySlug(func: eq(slug,"%s")) {
		uid
		slug
		reviewText
		game {
			uid
			gameTitle
			gameDetailsText
			gameStatus {
				uid
				statusName
			}
			dgraph.type
		}
		author {
			uid
			authorName
			dgraph.type
		}
		pull
		img
		dgraph.type
	}}`
	query := fmt.Sprintf(queryFmt, slug)
	//fmt.Println(query)
	resp, err := d.Dgraph.NewReadOnlyTxn().Query(context.Background(), query)
	if err != nil {
		d.log.Error("error getting review", "error", err)
		return d.reviewToDreview(gametime.Review{Slug: slug}), err
	}
	type reviews struct {
		Dreview []Dreview `json:"reviewBySlug"`
	}
	var dr reviews
	err = json.Unmarshal(resp.Json, &dr)
	//fmt.Println(string(resp.Json))
	if err != nil {
		d.log.Error("error unmarshalling review", "error", err)
		return d.reviewToDreview(gametime.Review{Slug: slug}), err
	}
	if len(dr.Dreview) < 1 {
		d.log.Error("error no review results in db", "slug", slug)
		return d.reviewToDreview(gametime.Review{Slug: slug}), errors.New("no review in db")
	}
	d.log.Info("get review by slug exit", "review id", dr.Dreview[0].Uid, "slug", dr.Dreview[0].Slug)
	return dr.Dreview[0], nil
}

func (d *Dgraph) insertReview(dr Dreview) error {
	d.newClient(d.cfg.Dgraph.Url)
	d.log.Info("insert review")
	b, err := json.Marshal(dr)
	mu := api.Mutation{
		SetJson:   b,
		CommitNow: true,
	}
	//d.log.Info(string(mu.SetJson))
	resp, err := d.Dgraph.NewTxn().Mutate(context.Background(), &mu)
	if err != nil {
		d.log.Error("error inserting review", "slug", dr.Slug, "error", err)
		return err
	}
	d.log.Info("", "uids", resp.Uids)
	return nil
}

func (d *Dgraph) dReviewToReview(da Dreview) gametime.Review {
	text := da.ReviewText
	var p gametime.Review
	d.log.Info("dreview", "lrn", len([]byte(text)))
	if err := json.Unmarshal([]byte(text), &p); err != nil {
		d.log.Error("error unmarshalling review text", "error", err, "slug", da.Slug, "id", da.Uid)
	}
	p.Text = text
	game := d.dGameToGame(da.Game)
	author := d.dAuthorToAuthor(da.Author)
	p.Id = da.Uid
	p.Game = game
	p.Author = author

	return p
}

func (d *Dgraph) reviewToDreview(a gametime.Review) Dreview {

	dr := Dreview{}
	dr.Uid = a.Id
	dg := d.gameToDgame(a.Game)
	da := d.authorToDAuthor(a.Author)
	dr.Author = da
	dr.Game = dg
	dr.Slug = a.Slug
	dr.Pull = a.Pull
	dr.Img = a.Imgs
	dr.ReviewText = a.Text
	dr.DgraphType = reviewType
	return dr
}

func (d *Dgraph) GetReviewById(id string) gametime.Review {
	d.log.Info("get review by id (app)")

	return d.dReviewToReview(d.getReviewById(id))
}

func (d *Dgraph) getReviewById(id string) Dreview {
	d.newClient(d.cfg.Dgraph.Url)
	d.log.Info("get review by id")
	queryFmt := `{reviewById(func: uid(%s)) {
		uid
		slug
		reviewText
		game {
			uid
			gameTitle
			gameDetails
			dgraph.type
		}
		author {
			uid
			authorName
			dgraph.type
		}
		pull
		img
		dgraph.type
	}}`
	query := fmt.Sprintf(queryFmt, id)
	resp, err := d.Dgraph.NewReadOnlyTxn().Query(context.Background(), query)
	if err != nil {
		d.log.Error("error getting review by id")
		return Dreview{Uid: id, DgraphType: reviewType}
	}
	type reviews struct {
		Dreview []Dreview `json:"reviewById"`
	}
	var r reviews
	err = json.Unmarshal(resp.Json, &r)
	if len(r.Dreview) < 1 || err != nil {
		d.log.Error("error getting review", "error", err)
		return Dreview{Uid: id, DgraphType: reviewType}
	}
	//fmt.Println(r.Dreview[0])
	return r.Dreview[0]
}
