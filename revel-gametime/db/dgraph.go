package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gametime/conf"
	"time"

	"github.com/dgraph-io/dgo/v200"
	"github.com/dgraph-io/dgo/v200/protos/api"
	"github.com/revel/revel/logger"
	"google.golang.org/grpc"
)

const ()

var (
	emptyDpost         = Dpost{}
	emptyDreview       = Dreview{}
	emptyDgame         = Dgame{}
	emptyDauthor       = Dauthor{}
	postType     DType = []string{"Post"}
	reviewType   DType = []string{"Review"}
	gameType     DType = []string{"Game"}
	authorType   DType = []string{"Author"}
)

type DType []string

type Dpost struct {
	Uid        string    `json:"uid,omitempty"`
	ReviewId   string    `json:"review"`
	Date       time.Time `json:"postDate"`
	DgraphType DType     `json:"dgraph.type"`
}

type Dreview struct {
	Uid        string   `json:"uid,omitempty"`
	ReviewText string   `json:"reviewText,omitempty"`
	Game       Dgame    `json:"game,omitempty"`
	DgraphType DType    `json:"dgraph.type"`
	Slug       string   `json:"slug,omitempty"`
	Author     Dauthor  `json:"author,omitempty"`
	Pull       string   `json:"pull,omitempty"`
	Post       Dpost    `json:"post,omitempty"`
	Img        []string `json:"img,omitempty"`
}

type Dgame struct {
	Uid         string `json:"uid,omitempty"`
	GameTitle   string `json:"gameTitle,omitempty"`
	GameDetails string `json:"gameDetails,omitempty"`
	DgraphType  DType  `json:"dgraph.type"`
}

type Dauthor struct {
	Uid        string `json:"uid,omitempty"`
	AuthorName string `json:"authorName,omitempty"`
	DgraphType DType  `json:"dgraph.type"`
}

type Dbug struct {
}

type Dfaq struct {
}

type Dcomment struct {
}

type Dgraph struct {
	*dgo.Dgraph
	log logger.MultiLogger
	cfg conf.Config
}

func (d *Dgraph) newClient(grpcUrl string) {
	// Dial a gRPC connection. The address to dial to can be configured when
	// setting up the dgraph cluster.
	dg, err := grpc.Dial(grpcUrl, grpc.WithInsecure())
	if err != nil {
		panic("no connection")
	}

	d.Dgraph = dgo.NewDgraphClient(
		api.NewDgraphClient(dg),
	)
}

func NewDgraph(log logger.MultiLogger, cfg conf.Config) *Dgraph {
	d := Dgraph{
		log: log,
		cfg: cfg,
	}
	d.newClient(cfg.Dgraph.Url)
	return &d
}

func (d *Dgraph) GetAuthorByName(name string) conf.Author {
	d.log.Info("get author by name", "(app)")
	da := d.getAuthorByName(name)
	return d.dAuthorToAuthor(da)
}

func (d *Dgraph) getAuthorByName(name string) Dauthor {
	d.newClient(d.cfg.Dgraph.Url)
	d.log.Info("get author by name")
	queryFmt := `{authorByName(func: eq(authorName,"%s")) {
		uid
		authorName
		dgraph.type
	}}`
	if name == "" {
		return d.authorToDAuthor(conf.Author{Name: name})
	}
	query := fmt.Sprintf(queryFmt, name)
	// fmt.Println(query)
	resp, err := d.Dgraph.NewReadOnlyTxn().Query(context.Background(), query)
	if err != nil {
		d.log.Error("error getting author", "error", err)
		return d.authorToDAuthor(conf.Author{Name: name})
	}
	type authors struct {
		Dauthor []Dauthor `json:"authorByName"`
	}
	var dAuthor authors
	err = json.Unmarshal(resp.Json, &dAuthor)
	fmt.Println(string(resp.Json))
	if err != nil {
		d.log.Error("error unmarshalling author", "error", err)
		return d.authorToDAuthor(conf.Author{Name: name})
	}
	if len(dAuthor.Dauthor) < 1 {
		d.log.Error("error no author results in db", "name", name)
		return d.authorToDAuthor(conf.Author{Name: name})
	}
	return dAuthor.Dauthor[0]
}

func (d *Dgraph) insertAuthor(dAuthor Dauthor) error {
	d.newClient(d.cfg.Dgraph.Url)
	d.log.Info("insert author")
	b, err := json.Marshal(dAuthor)
	mu := api.Mutation{
		SetJson:   b,
		CommitNow: true,
	}
	//d.log.Info(string(mu.SetJson))
	resp, err := d.Dgraph.NewTxn().Mutate(context.Background(), &mu)
	if err != nil {
		d.log.Error("error inserting author", "author", dAuthor.AuthorName)
		return err
	}
	d.log.Info("", "uids", resp.Uids)
	return nil
}

func (d *Dgraph) dAuthorToAuthor(da Dauthor) conf.Author {
	d.log.Info("dauthor -> author")
	return conf.Author{
		Id:   da.Uid,
		Name: da.AuthorName,
	}
}

func (d *Dgraph) authorToDAuthor(a conf.Author) Dauthor {
	d.log.Info("author -> dauthor")
	return Dauthor{
		Uid:        a.Id,
		AuthorName: a.Name,
		DgraphType: authorType,
	}
}

func (d *Dgraph) GetGameByTitle(Title string) conf.Game {
	d.log.Info("get game by title (app)")

	da := d.getGameByTitle(Title)
	return d.dGameToGame(da)
}

func (d *Dgraph) getGameByTitle(title string) Dgame {
	d.newClient(d.cfg.Dgraph.Url)
	d.log.Info("get game by title")
	queryFmt := `{gameByTitle(func: eq(gameTitle,"%s")) {
		uid
		gameTitle
		gameDetails
		dgraph.type
	}}`
	query := fmt.Sprintf(queryFmt, title)
	//fmt.Println(query)
	resp, err := d.Dgraph.NewReadOnlyTxn().Query(context.Background(), query)
	if err != nil {
		d.log.Error("error getting game", "error", err)
		return d.gameToDgame(conf.Game{Title: title})
	}
	type games struct {
		Dgame []Dgame `json:"gameByTitle"`
	}
	var dgame games
	err = json.Unmarshal(resp.Json, &dgame)
	fmt.Println(string(resp.Json))
	if err != nil {
		d.log.Error("error unmarshalling game", "error", err)
		return d.gameToDgame(conf.Game{Title: title})
	}
	if len(dgame.Dgame) < 1 {
		d.log.Error("error no game results in db", "title", title)
		return d.gameToDgame(conf.Game{Title: title})
	}
	return dgame.Dgame[0]
}

func (d *Dgraph) insertGame(da Dgame) error {
	d.newClient(d.cfg.Dgraph.Url)
	d.log.Info("insert game")
	b, err := json.Marshal(da)
	mu := api.Mutation{
		SetJson:   b,
		CommitNow: true,
	}
	//d.log.Info(string(mu.SetJson))
	resp, err := d.Dgraph.NewTxn().Mutate(context.Background(), &mu)
	if err != nil {
		d.log.Error("error inserting game", "game", da.GameTitle)
		return err
	}
	d.log.Info("", "uids", resp.Uids)
	return nil
}

func (d *Dgraph) dGameToGame(da Dgame) conf.Game {
	d.log.Info("dgame -> gmae")

	return conf.Game{
		Id:      da.Uid,
		Title:   da.GameTitle,
		Details: da.GameDetails,
	}
}

func (d *Dgraph) gameToDgame(a conf.Game) Dgame {
	d.log.Info("game -> dgame")

	return Dgame{
		Uid:         a.Id,
		GameTitle:   a.Title,
		GameDetails: a.Details,
		DgraphType:  gameType,
	}
}

func (d *Dgraph) GetReviewBySlug(slug string) (conf.Review, error) {
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
	query := fmt.Sprintf(queryFmt, slug)
	//fmt.Println(query)
	resp, err := d.Dgraph.NewReadOnlyTxn().Query(context.Background(), query)
	if err != nil {
		d.log.Error("error getting review", "error", err)
		return d.reviewToDreview(conf.Review{Slug: slug}), err
	}
	type reviews struct {
		Dreview []Dreview `json:"reviewBySlug"`
	}
	var dr reviews
	err = json.Unmarshal(resp.Json, &dr)
	//fmt.Println(string(resp.Json))
	if err != nil {
		d.log.Error("error unmarshalling review", "error", err)
		return d.reviewToDreview(conf.Review{Slug: slug}), err
	}
	if len(dr.Dreview) < 1 {
		d.log.Error("error no review results in db", "slug", slug)
		return d.reviewToDreview(conf.Review{Slug: slug}), errors.New("no review in db")
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

func (d *Dgraph) dReviewToReview(da Dreview) conf.Review {
	d.log.Info("dreview -> review")

	text := da.ReviewText
	var p conf.Review
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

func (d *Dgraph) reviewToDreview(a conf.Review) Dreview {
	d.log.Info("review -> dreview")

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

func (d *Dgraph) GetReviewById(id string) conf.Review {
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

func (d *Dgraph) GetPostBySlug(slug string) (conf.Post, error) {
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

func (d *Dgraph) GetPostByReview(review conf.Review) (conf.Post, error) {
	d.log.Info("get post by review (app)")
	dr := d.reviewToDreview(review)
	dp, err := d.getPostByReview(dr)
	if err != nil {
		return conf.Post{ReviewId: review.Id}, err
	}
	return conf.Post{
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

func (d *Dgraph) GetPostsByMostRecent() ([]conf.Post, error) {
	d.log.Info("get posts by recent (app)")
	dPosts, err := d.getPostsByMostRecent()
	if err != nil {
		d.log.Error("error getting posts recent", "error", err)
		return []conf.Post{}, err
	}
	var posts []conf.Post
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

//reviews and posts are 1:1
func (d *Dgraph) InsertPost(review conf.Review) error {
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

func (d *Dgraph) Dump() error {
	//dump dgraph contents to json? or that 3 type format
	return nil
}

func (d *Dgraph) dPostToPost(dp Dpost) conf.Post {
	d.log.Info("dPost -> post")
	return conf.Post{
		Id:       dp.Uid,
		ReviewId: dp.ReviewId,
		Date:     dp.Date,
	}
}

func (d *Dgraph) postToDpost(p conf.Post) Dpost {
	d.log.Info("post -> dPost")
	return Dpost{
		Uid:        p.Id,
		ReviewId:   p.ReviewId,
		Date:       p.Date,
		DgraphType: postType,
	}
}
