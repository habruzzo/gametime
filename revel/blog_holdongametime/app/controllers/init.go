package controllers

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/revel/revel"
)

// HOLDEN's
const (
	jsonSlugPath string = "/opt/gametime/reviews/json/review_map.json"
	jsonGamePath string = "/opt/gametime/reviews/json/game_list.json"
	jsonPostPath string = "/opt/gametime/reviews/json/post_list.json"
	SlugFiles    string = "slug:files"
	SlugPosts    string = "slug:posts"
	SlugGames    string = "slug:games"
	GamesStatus  string = "games:status"
	PostsOrdered string = "posts:ordered"
	GameKey      string = "game:"
	PostKey      string = "post:"
)

var Pool *redis.Pool

type JsonMapping struct {
	Slug string `redis:"slug"`
	File string `redis:"file"`
}

func CheckDBLoaded(conn redis.Conn) bool {
	b, err := redis.Bool(conn.Do("EXISTS", PostsOrdered))
	if err != nil {
		panic(err)
	}
	return b
}

func CheckDB() {
	var err error
	// init db
	conn := Pool.Get()
	if err != nil {
		fmt.Println("FATAL", err)
		panic(err)
	}
	if loaded := CheckDBLoaded(conn); !loaded {
		InitRedisDB(conn)
	}
	fmt.Println("LOADED DB")
	conn.Close()
}

func startPool() {
	Pool = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", "localhost:7001") },
	}
}

func init() {
	//if redis, uncomment this
	//revel.OnAppStart(startPool)
	//revel.OnAppStart(CheckDB)
	//if gorm, uncomment this
	revel.OnAppStart(InitGormDB)
	revel.InterceptMethod((*GormController).SetDB, revel.BEFORE)
}
