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
	Id      string  `json:"uid"`
	Title   string  `json:"title"`
	Status  Status  `json:"status"`
	Details Details `json:"details"`
}

type Details struct {
	Developer   string `json:"developer"`
	InstallDate string `json:"installed"`
}

type Status struct {
	Id   string `json:"uid"`
	Name string `json:"name"`
}

var (
	Unknown   Status = Status{Name: "unknownStatus"}
	Wishlist  Status = Status{Name: "wishlistStatus"}
	Installed Status = Status{Name: "installedStatus"}
	Played    Status = Status{Name: "playedStatus"}
	Completed Status = Status{Name: "completedStatus"}
	Reviewed  Status = Status{Name: "reviewedStatus"}
)

func ToStatus(str string) Status {
	switch str {
	case Unknown.Name:
		return Unknown
	case Wishlist.Name:
		return Wishlist
	case Installed.Name:
		return Installed
	case Played.Name:
		return Played
	case Completed.Name:
		return Completed
	case Reviewed.Name:
		return Reviewed
	default:
		return Unknown
	}
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
