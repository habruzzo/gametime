package app

import (
	"context"
	"encoding/json"
	"fmt"
	"gametime"
	"gametime/config"
	"gametime/db"
	"io/ioutil"
	"os/exec"
	"strings"

	"github.com/dgraph-io/dgo/v200/protos/api"
	"github.com/sirupsen/logrus"
)

var (
	// AppVersion revel app version (ldflags)
	AppVersion string

	// BuildTime revel app build-time (ldflags)
	BuildTime string
)

func init() {
	// revel.DevMod and revel.RunMode work here
	// // Use this script to check for dev mode and set dev/prod startup scripts here!
	// if revel.DevMode == true {
	// 	// Dev mode
	// }
	b := getReviewFiles()
	g := getGameFiles()
	fmt.Println("INIT DB FUCKERS")
	cfg := config.LoadConfig()
	log := logrus.New()
	dgraph := db.NewDgraph(log, cfg)
	//from the perspective of root of project
	setupSchema(dgraph)
	loadStatuses(dgraph)
	loadGames(string(g), dgraph)
	loadReviews(string(b), dgraph)
}

func setupSchema(dgraph *db.Dgraph) {
	fn := "./config/schema.dql"
	buf, err := ioutil.ReadFile(fn)
	fmt.Println("ERROR READING SCHEMA??", err)
	schema := string(buf)

	op := api.Operation{Schema: schema}
	err = dgraph.Alter(context.TODO(), &op)
	fmt.Println("ERROR ALTERING SCHEMA??", err)
}

func loadReviews(b string, d *db.Dgraph) {
	for _, v := range strings.Split(b, "\n") {
		loadReview(v, d)
	}
}

func loadReview(v string, d *db.Dgraph) {
	if v == "" {
		return
	}
	postBuf, err := ioutil.ReadFile(v)
	fmt.Println("ERROR READING JSON??", err, "slug", v)

	var p gametime.Review
	err = json.Unmarshal(postBuf, &p)
	fmt.Println("ERROR UNMARSHALLING JSON??", err)
	p.Text = string(postBuf)
	err = d.InsertPost(p)
	fmt.Println("ERROR INSERTING POST?", err)
}

func getReviewFiles() []byte {
	cmd := exec.Command("find", "../reviews/json", "-type", "f", "-regex", ".*rubric.*")
	b, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		return []byte{}
	}
	return b
}

func loadGames(b string, d *db.Dgraph) {
	for _, v := range strings.Split(b, "\n") {
		loadGameFile(v, d)
	}
}

func loadGameFile(v string, d *db.Dgraph) {
	if v == "" {
		return
	}
	gameBuf, err := ioutil.ReadFile(v)
	fmt.Println("ERROR READING JSON??", err, "file", v)

	var p []gametime.Game
	err = json.Unmarshal(gameBuf, &p)
	fmt.Println("ERROR UNMARSHALLING JSON??", err)
	fmt.Println(len(p), p[0])
	for _, v := range p {
		err = d.InsertGame(v)
		fmt.Println("ERROR INSERTING GAME?", err)
	}

}

func getGameFiles() []byte {
	cmd := exec.Command("find", "../reviews/games/json", "-type", "f")
	b, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		return []byte{}
	}
	return b
}

func loadStatuses(d *db.Dgraph) {
	d.InsertStatus(gametime.Unknown)
	d.InsertStatus(gametime.Wishlist)
	d.InsertStatus(gametime.Installed)
	d.InsertStatus(gametime.PlayedSome)
	d.InsertStatus(gametime.PlayedMost)
	d.InsertStatus(gametime.WontReview)
	d.InsertStatus(gametime.Completed)
	d.InsertStatus(gametime.Reviewed)
}
