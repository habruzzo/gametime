package models

import (
	"strings"
	"time"

	"github.com/google/uuid"
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
	Game        *Game
	Status      PublishStatus
	ContentPath string
	Tags        []Tag
	Rating      int
	PublishDate time.Time
}

func NewPost(title string, game *Game, status PublishStatus, contentPath string, tags []Tag, rating int, publishDate time.Time) *Post {
	return &Post{
		Id:          uuid.New(),
		Title:       title,
		Slug:        ToSlug(title),
		Game:        game,
		Status:      status,
		ContentPath: contentPath,
		Tags:        tags,
		Rating:      rating,
		PublishDate: publishDate,
	}
}

func ToSlug(title string) string {
	slug := strings.ToLower(strings.ReplaceAll(title, " ", "_"))
	if words := strings.Count(title, " "); words > 3 {
		slug = slug[:3]
	}
	return slug
}
