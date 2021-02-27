package models

import (
	"fmt"

	"github.com/jinzhu/gorm"

	"encoding/json"
	"hash/fnv"
	"time"
)

type GameStatus string

const (
	Unknown   GameStatus = "Unknown"
	Acquired             = "Acquired"
	Started              = "Started"
	Completed            = "Completed"
	Reviewed             = "Reviewed"
	Suggested            = "Suggested"
	Published            = "Published"
	Discarded            = "Discarded"
)

type BugStatus string

const (
	Logged BugStatus = "Logged"
	Fixed            = "Fixed"
)

type PlatformType string

const (
	PC          PlatformType = "PC"
	GameBoy                  = "GameBoy"
	PlayStation              = "PlayStation"
	Xbox                     = "Xbox"
	NintendoDS               = "Nintendo DS"
	Wii                      = "Wii"
	Switch                   = "Switch"
	Mobile                   = "Mobile"
	Other                    = "Other"
)

type Tag struct {
	gorm.Model
	id    uint32
	value string
}

func NewTag(value string) *Tag {
	return &Tag{
		id:    generateId(value),
		value: value,
	}
}

type Game struct {
	gorm.Model
	//Id          uint32
	Title       string
	Slug        string
	Platform    PlatformType
	Publisher   string
	Creator     string
	ReleaseDate string
	SteamLink   string
	Status      GameStatus
}

func (g Game) Value() ([]byte, error) {
	return json.Marshal(g)
}

func NewGame(title string, slug string, platform PlatformType, publisher string, creator string, releaseDate string, steamLink string, status GameStatus) *Game {
	return &Game{
		//Id:          generateId(title),
		Title:       title,
		Slug:        slug,
		Platform:    platform,
		Publisher:   publisher,
		Creator:     creator,
		ReleaseDate: releaseDate,
		SteamLink:   steamLink,
		Status:      status,
	}
}

type PostJson struct {
	Title string
	Slug  string
}

type Post struct {
	gorm.Model
	//Id          uint32
	Title       string
	Slug        string
	GameId      uint
	ContentPath string
	Rating      int
	PublishDate time.Time
}

func NewPost(title string, gameId uint, slug string, contentPath string, rating int, publishDate time.Time) *Post {
	return &Post{
		//Id:          generateId(title),
		Title:       title,
		Slug:        slug,
		GameId:      gameId,
		ContentPath: contentPath,
		Rating:      rating,
		PublishDate: publishDate,
	}
}

func (p *Post) Value() string {
	return fmt.Sprintf("%s %s %s %s %s %s", p.Title, p.Slug, p.GameId, p.Status, p.Rating, p.ContentPath)
}

func generateId(title string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(title))
	return h.Sum32()
}
