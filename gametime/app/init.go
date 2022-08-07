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
	cmd := exec.Command("find", "../reviews/json", "-type", "f", "-regex", ".*rubric.*")
	b, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("INIT DB FUCKERS")
	log := logrus.New()
	cfg := config.LoadConfig()
	dgraph := db.NewDgraph(log, cfg)
	//from the perspective of root of project
	fn := "./config/schema.dql"
	buf, err := ioutil.ReadFile(fn)
	fmt.Println("ERROR READING SCHEMA??", err)
	schema := string(buf)

	op := api.Operation{Schema: schema}
	err = dgraph.Alter(context.TODO(), &op)
	fmt.Println("ERROR ALTERING SCHEMA??", err)
	for _, v := range strings.Split(string(b), "\n") {
		if v == "" {
			continue
		}
		postBuf, err := ioutil.ReadFile(v)
		fmt.Println("ERROR READING JSON??", err, "slug", v)

		var p gametime.Review
		err = json.Unmarshal(postBuf, &p)
		fmt.Println("ERROR UNMARSHALLING JSON??", err)
		p.Text = string(postBuf)
		err = dgraph.InsertPost(p)
		fmt.Println("ERROR INSERTING POST?", err)
	}
}
