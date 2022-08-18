package gametime

import (
	"encoding/json"
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
	Notes       string `json:"notes"`
}

type Status struct {
	Id   string `json:"uid"`
	Name string `json:"name"`
}

var (
	Unknown    Status = Status{Name: "unknownStatus"}
	Wishlist   Status = Status{Name: "wishlistStatus"}
	Installed  Status = Status{Name: "installedStatus"}
	PlayedSome Status = Status{Name: "playedSomeStatus"}
	PlayedMost Status = Status{Name: "playedMostStatus"}
	Completed  Status = Status{Name: "completedStatus"}
	Reviewed   Status = Status{Name: "reviewedStatus"}
	WontReview Status = Status{Name: "wontReviewStatus"}
)

func ToStatus(str string) Status {
	switch str {
	case Unknown.Name:
		return Unknown
	case Wishlist.Name:
		return Wishlist
	case Installed.Name:
		return Installed
	case PlayedSome.Name:
		return PlayedSome
	case PlayedMost.Name:
		return PlayedMost
	case Completed.Name:
		return Completed
	case Reviewed.Name:
		return Reviewed
	case WontReview.Name:
		return WontReview
	default:
		return Unknown
	}
}

func ToDetails(str string) (d Details) {
	json.Unmarshal([]byte(str), &d)
	return d
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

type Issue struct {
	Id     string `json:"uid"`
	Status string `json:"status"`
}
