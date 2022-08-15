package db

import (
	"bytes"
	"errors"
	"fmt"
	"gametime"
	"gametime/config"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/dgraph-io/dgo/v200"
	"github.com/dgraph-io/dgo/v200/protos/api"
	"google.golang.org/grpc"
)

var (
	postType   DType = []string{"Post"}
	reviewType DType = []string{"Review"}
	gameType   DType = []string{"Game"}
	authorType DType = []string{"Author"}
	statusType DType = []string{"Status"}
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
	Uid        string      `json:"uid,omitempty"`
	GameTitle  string      `json:"gameTitle,omitempty"`
	Status     DgameStatus `json:"gameStatus"`
	Text       string      `json:"gameDetailsText"`
	DgraphType DType       `json:"dgraph.type"`
}

type DgameStatus struct {
	Uid        string `json:"uid,omitempty"`
	Name       string `json:"statusName"`
	DgraphType DType  `json:"dgraph.type"`
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
	log gametime.Logger
	cfg *config.Config
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

func NewDgraph(log gametime.Logger, cfg *config.Config) *Dgraph {
	d := Dgraph{
		log: log,
		cfg: cfg,
	}
	d.newClient(cfg.Dgraph.Url)
	return &d
}

func (d *Dgraph) Dump() error {
	d.newClient(d.cfg.Dgraph.Url)
	dumpMutation := `mutation {
						export(input: {}) {
						response {
							message
							code
						}
						}
					}`
	target := fmt.Sprintf("http://%s%s", d.cfg.Dgraph.HttpUrl, "/admin")
	b := []byte(dumpMutation)
	req, err := http.NewRequest(http.MethodPost, target, bytes.NewReader(b))
	if err != nil {
		d.log.Error("error creating request for dump", err)
	}
	req.ContentLength = int64(len(b))
	req.Header.Add("Content-Type", "application/graphql")
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		d.log.Error("problem sending message to dgraph!", "error", err)
		return err
	}
	bb, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode > 299 {
		d.log.Error(string(bb))
		return errors.New(string(bb))
	}
	return nil
}
