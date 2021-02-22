package models

import (
	"github.com/google/uuid"

	"encoding/json"
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

type PublishStatus string

const (
	Draft   PublishStatus = "Draft"
	Publish               = "Publish"
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
	id    uuid.UUID
	value string
}

func NewTag(value string) *Tag {
	return &Tag{
		id:    uuid.New(),
		value: value,
	}
}

type Game struct {
	Id          uuid.UUID
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
		Id:          uuid.New(),
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
	Title  string
	Slug   string
	Status PublishStatus
	Rating int
}

type Post struct {
	Id          uuid.UUID
	Title       string
	Slug        string
	GameId      uuid.UUID
	Status      PublishStatus
	ContentPath string
	Rating      int
	PublishDate time.Time
}

func NewPost(title string, gameId uuid.UUID, slug string, status PublishStatus, contentPath string, rating int, publishDate time.Time) *Post {
	return &Post{
		Id:          uuid.New(),
		Title:       title,
		Slug:        slug,
		GameId:      gameId,
		Status:      status,
		ContentPath: contentPath,
		Rating:      rating,
		PublishDate: publishDate,
	}
}
