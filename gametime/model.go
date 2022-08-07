package gametime

import (
	"time"
)

type Post struct {
	ReviewId string    `json:"review"`
	Id       string    `json:"uid"`
	Date     time.Time `json:"postDate"`
}

type Review struct {
	Id       string     `json:"uid"`
	Overall  []Question `json:"Overall"`
	Art      Art
	Gameplay Gameplay
	Text     string   `json:"reviewText"`
	Game     Game     `json:"game"`
	Pull     string   `json:"pull"`
	Imgs     []string `json:"img"`
	Author   Author   `json:"author"`
	Slug     string   `json:"slug"`
	Post     Post     `json:"post"`
}

type Game struct {
	Id      string `json:"uid"`
	Title   string `json:"title"`
	Details string `json:"details"`
}

type Art struct {
	Graphics []Question `json:"Graphics"`
	Sound    []Question `json:"Sound"`
	Story    []Question `json:"Story"`
	Themes   []Question `json:"Themes"`
}

type Gameplay struct {
	Mechanics  []Question `json:"Mechanics"`
	Difficulty []Question `json:"Difficulty"`
	Experience []Question `json:"Experience"`
}

type Question struct {
	Id     string `json:"uid"`
	Prompt string `json:"question"`
	Answer string `json:"answer"`
}

type Author struct {
	Id   string `json:"uid"`
	Name string `json:"name"`
}
