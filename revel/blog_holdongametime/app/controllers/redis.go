package controllers

import (
	"blog_holdongametime/app/models"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/gomodule/redigo/redis"
	"github.com/revel/revel"
)

// type: revel controller with `redis.Conn`
type RedController struct {
	*revel.Controller
}

func getPost(postKey string) models.Post {
	conn := Pool.Get()
	v, err := redis.Values(conn.Do("HGETALL", postKey))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	var post models.Post
	if err := redis.ScanStruct(v, &post); err != nil {
		fmt.Println(err)
		panic(err)
	}
	conn.Close()
	return post
}

func (r RedController) GetPostAndFile(slug string) (models.Post, string) {
	conn := Pool.Get()
	pathString, err := redis.String(conn.Do("HGET", redis.Args{}.Add("slug:files").Add(slug)))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	postKey, err := redis.String(conn.Do("HGET", redis.Args{}.Add("slug:posts").Add(slug)))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	post := getPost(postKey)
	conn.Close()
	return post, pathString
}

func (r RedController) GetPostList() []*models.Post {
	var postList []*models.Post
	conn := Pool.Get()
	postIds, err := redis.Strings(conn.Do("HGETALL", redis.Args{}.Add("posts-ordered")))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	for _, v := range postIds {
		post := getPost(v)
		postList = append(postList, &post)
	}
	conn.Close()
	return postList
}

func LoadSlugs() []JsonMapping {
	data, err := ioutil.ReadFile(jsonSlugPath)
	if err != nil {
		panic(err)
	}
	var j []JsonMapping
	err = json.Unmarshal(data, &j)
	if err != nil {
		panic(err)
	}
	return j
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

func LoadPosts() []*models.Post {
	data, err := ioutil.ReadFile(jsonPostPath)
	conn := Pool.Get()
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
		slugKey := fmt.Sprintf("slug-game-%s", v.Slug)
		gId, err := redis.String(conn.Do("GET", redis.Args{}.Add(slugKey)))
		if err != nil {
			fmt.Println("FATAL", err)
			panic(err)
		}
		slugKey = fmt.Sprintf("slug-%s", v.Slug)
		path, err := redis.String(conn.Do("GET", redis.Args{}.Add(slugKey)))
		if err != nil {
			fmt.Println("FATAL", err)
			panic(err)
		}
		pClean = append(pClean, models.NewPost(v.Title, gId, v.Slug, path, v.PublishDate, v.Index))
	}
	conn.Close()
	return pClean
}

func InitRedisDB(conn redis.Conn) {
	for _, v := range LoadSlugs() {
		args := []string{"slugfiles", v.Slug, v.File}
		_, err := conn.Do("HSET", args)
		if err != nil {
			fmt.Println(args)
			fmt.Println("problem!")
			panic(err)
		}
		b, err := redis.Bool(conn.Do("EXISTS", redis.Args{}.Add(SlugFiles).Add(v.Slug)))
		if err != nil {
			panic(err)
		}
		if !b {
			panic(err)
		}
	}

	for _, v := range LoadGames() {
		gameKey := fmt.Sprintf("game-%v", v.Id)
		_, err := conn.Do("ZADD", redis.Args{}.Add("games:status").Add(v.Status).Add(gameKey))
		if err != nil {
			panic(err)
		}
		_, err = conn.Do("HSET", redis.Args{}.Add(gameKey).AddFlat(v))
		if err != nil {
			panic(err)
		}
		_, err = conn.Do("HSET", redis.Args{}.Add("slug:games").Add(v.Slug).Add(gameKey))
		if err != nil {
			panic(err)
		}
		b, err := redis.Bool(conn.Do("EXISTS", redis.Args{}.Add(gameKey)))
		if err != nil {
			panic(err)
		}
		if !b {
			panic(fmt.Sprintf("%s not exists", gameKey))
		}
		b, err = redis.Bool(conn.Do("HEXISTS", redis.Args{}.Add("slug:games").Add(v.Slug)))
		if err != nil {
			panic(err)
		}
		if !b {
			panic(fmt.Sprintf("%s not exists in slug:files", v.Slug))
		}
	}

	for _, v := range LoadPosts() {
		postKey := fmt.Sprintf("post-%v", v.Id)
		_, err := conn.Do("ZADD", redis.Args{}.Add("post-ordered").Add(v.Index).Add(postKey))
		if err != nil {
			panic(err)
		}
		_, err = conn.Do("HSET", redis.Args{}.Add(postKey).AddFlat(v))
		if err != nil {
			panic(err)
		}
		_, err = conn.Do("HSET", redis.Args{}.Add("slug:posts").Add(v.Slug).Add(postKey))
		if err != nil {
			panic(err)
		}
		b, err := redis.Bool(conn.Do("EXISTS", postKey))
		if err != nil {
			panic(err)
		}
		if !b {
			panic(fmt.Sprintf("%s not exists", postKey))
		}
		b, err = redis.Bool(conn.Do("HEXISTS", redis.Args{}.Add("slug:posts").Add(v.Slug)))
		if err != nil {
			panic(err)
		}
		if !b {
			panic(fmt.Sprintf("%s not exists in slug:posts", v.Slug))
		}
	}
}
