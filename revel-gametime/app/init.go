package app

import (
	"context"
	"encoding/json"
	"fmt"
	"gametime/conf"
	"gametime/db"
	"io/ioutil"
	"os/exec"
	"strings"

	_ "github.com/revel/modules"
	"github.com/revel/revel"
	"github.com/revel/revel/logger"

	"github.com/dgraph-io/dgo/v200/protos/api"
)

var (
	// AppVersion revel app version (ldflags)
	AppVersion string

	// BuildTime revel app build-time (ldflags)
	BuildTime string
)

func init() {
	// Filters is the default set of global filters.
	revel.Filters = []revel.Filter{
		revel.PanicFilter,             // Recover from panics and display an error page instead.
		revel.RouterFilter,            // Use the routing table to select the right Action
		revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		revel.ParamsFilter,            // Parse parameters into Controller.Params.
		revel.SessionFilter,           // Restore and write the session cookie.
		revel.FlashFilter,             // Restore and write the flash cookie.
		revel.ValidationFilter,        // Restore kept validation errors and save new ones from cookie.
		revel.I18nFilter,              // Resolve the requested language
		HeaderFilter,                  // Add some security based headers
		revel.InterceptorFilter,       // Run interceptors around the action.
		revel.CompressFilter,          // Compress the result.
		revel.BeforeAfterFilter,       // Call the before and after filter functions
		revel.ActionInvoker,           // Invoke the action.
	}

	// Register startup functions with OnAppStart
	// revel.DevMode and revel.RunMode only work inside of OnAppStart. See Example Startup Script
	// ( order dependent )
	// revel.OnAppStart(ExampleStartupScript)
	revel.OnAppStart(InitDb)
	// revel.OnAppStart(FillCache)
}

// HeaderFilter adds common security headers
// There is a full implementation of a CSRF filter in
// https://github.com/revel/modules/tree/master/csrf
var HeaderFilter = func(c *revel.Controller, fc []revel.Filter) {
	c.Response.Out.Header().Add("X-Frame-Options", "SAMEORIGIN")
	c.Response.Out.Header().Add("X-XSS-Protection", "1; mode=block")
	c.Response.Out.Header().Add("X-Content-Type-Options", "nosniff")
	c.Response.Out.Header().Add("Referrer-Policy", "strict-origin-when-cross-origin")

	fc[0](c, fc[1:]) // Execute the next filter stage.
}

func InitDb() {
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
	log := logger.New()
	cfg := conf.LoadConfig()
	dgraph := db.NewDgraph(log, cfg)
	//from the perspective of root of project
	fn := "./conf/schema.dql"
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

		var p conf.Review
		err = json.Unmarshal(postBuf, &p)
		fmt.Println("ERROR UNMARSHALLING JSON??", err)
		p.Text = string(postBuf)
		err = dgraph.InsertPost(p)
		fmt.Println("ERROR INSERTING POST?", err)
	}
}
