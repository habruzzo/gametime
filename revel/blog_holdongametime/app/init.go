package app

import (
	"blog_holdongametime/app/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	gorp "github.com/revel/modules/orm/gorp/app"
	"github.com/revel/revel"
)

// HOLDEN's
const jsonSlugPath = "/opt/gametime/reviews/review_map.json"
const jsonGamePath = "/opt/gametime/reviews/game_list.json"
const jsonPostPath = "/opt/gametime/reviews/post_list.json"

var jsonMapping map[string]string

var (
	// AppVersion revel app version (ldflags)
	AppVersion string

	// BuildTime revel app build-time (ldflags)
	BuildTime string
)

// HOLDEN's
type JsonMapping struct {
	Slug string
	File string
}

// HOLDEN's
func LoadSlugs() {
	data, err := ioutil.ReadFile(jsonSlugPath)
	if err != nil {
		panic(err)
	}
	var j []JsonMapping
	err = json.Unmarshal(data, &j)
	if err != nil {
		panic(err)
	}
	for _, v := range j {
		jsonMapping[v.Slug] = v.File
		fmt.Println(jsonMapping[v.Slug])
	}
}

func LoadGames() []*models.Game {
	data, err := ioutil.ReadFile(jsonGamePath)
	if err != nil {
		panic(err)
	}
	var g []models.Game
	err = json.Unmarshal(data, &g)
	if err != nil {
		panic(err)
	}
	var gClean []*models.Game
	for _, v := range g {
		gClean = append(gClean, models.NewGame(v.Title, v.Slug, v.Platform, v.Publisher, v.Creator, v.ReleaseDate, v.SteamLink, v.Status))
	}
	return gClean
}

func MatchPostToGame(p models.PostJson, games []*models.Game) *models.Game {
	for _, v := range games {
		if p.Slug == v.Slug {
			return v
		}
	}
	panic("No matching game")
}

func LoadPosts(games []*models.Game) []*models.Post {
	data, err := ioutil.ReadFile(jsonPostPath)
	if err != nil {
		panic(err)
	}
	var p []models.PostJson
	err = json.Unmarshal(data, &p)
	if err != nil {
		panic(err)
	}
	var pClean []*models.Post
	for _, v := range p {
		g := MatchPostToGame(v, games)
		pClean = append(pClean, models.NewPost(v.Title, g, v.Status, jsonMapping[v.Slug], nil, v.Rating, time.Now()))
	}
	return pClean
}

func InitDB() {
	LoadSlugs()
	games := LoadGames()
	posts := LoadPosts(games)
	gorp.Db.SetDbInit(func(dbGorp *gorp.DbGorp) error {
		// Register tables
		gorp.Db.Map.AddTableWithName(models.Game{}, "game")
		gorp.Db.Map.AddTableWithName(models.Post{}, "post")
		gorp.Db.Map.CreateTables()
		for _, v := range games {
			if err := dbGorp.Insert(v); err != nil {
				panic(err)
			}
		}
		for _, v := range posts {
			if err := dbGorp.Insert(v); err != nil {
				panic(err)
			}
		}
		return nil
	})
}

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

	// HOLDEN's
	jsonMapping = make(map[string]string)
	revel.OnAppStart(InitDB)

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

//func ExampleStartupScript() {
//	// revel.DevMod and revel.RunMode work here
//	// Use this script to check for dev mode and set dev/prod startup scripts here!
//	if revel.DevMode == true {
//		// Dev mode
//	}
//}
