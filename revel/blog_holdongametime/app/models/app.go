package models

import (
	"fmt"
	"hash/fnv"

	"encoding/json"

	"github.com/google/uuid"
	"gorm.io/gorm"
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
	id    uuid.UUID
	value string
}

func NewTag(value string) *Tag {
	return &Tag{
		id:    uuid.New(),
		value: value,
	}
}

type JsonMapping struct {
	gorm.Model
	Slug string `redis:"slug"`
	File string `redis:"file"`
}

type Game struct {
	Id          uuid.UUID    `redis:"id"`
	Title       string       `redis:"title"`
	Slug        string       `redis:"slug"`
	Platform    PlatformType `redis:"platform"`
	Publisher   string       `redis:"publisher"`
	Creator     string       `redis:"creator"`
	ReleaseDate string       `redis:"release"`
	SteamLink   string       `redis:"steam"`
	Status      GameStatus   `redis:"status"`
}

type GormGame struct {
	gorm.Model
	Id          string `gorm:"primaryKey"`
	Title       string
	Slug        string
	Platform    PlatformType
	Publisher   string
	Creator     string
	ReleaseDate string
	GamesDbLink string
	Status      GameStatus `gorm:"index"`
	Images      string
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

func NewGormGame(title string, slug string, platform PlatformType, publisher string, creator string, releaseDate string, gamesDbLink string, status GameStatus, imgs string) *GormGame {
	id := GenerateId(slug, "game")
	return &GormGame{
		Id:          id,
		Title:       title,
		Slug:        slug,
		Platform:    platform,
		Publisher:   publisher,
		Creator:     creator,
		ReleaseDate: releaseDate,
		GamesDbLink: gamesDbLink,
		Status:      status,
		Images:      imgs,
	}
}

type PostJson struct {
	Title       string
	Slug        string
	PublishDate string
	Index       int
}

type Post struct {
	Id          uuid.UUID `redis:"id"`
	Title       string    `redis:"title"`
	Slug        string    `redis:"slug"`
	GameId      string    `redis:"gameId"`
	ContentPath string    `redis:"contentPath"`
	PublishDate string    `redis:"publishDate"`
	Index       int       `redis:"postIndex"`
}

type GormPost struct {
	gorm.Model
	Id          string `gorm:"primaryKey"`
	Title       string
	Slug        string
	Game        *GormGame `gorm:"foreignKey:Id"`
	ContentPath string
	PublishDate string
}

func NewPost(title string, gameId string, slug string, contentPath string, publishDate string, index int) *Post {
	return &Post{
		Id:          uuid.New(),
		Title:       title,
		Slug:        slug,
		GameId:      gameId,
		ContentPath: contentPath,
		PublishDate: publishDate,
		Index:       index,
	}
}

func NewGormPost(title string, game *GormGame, slug string, contentPath string, publishDate string) *GormPost {
	id := GenerateId(slug, "post")
	return &GormPost{
		Id:          id,
		Title:       title,
		Slug:        slug,
		Game:        game,
		ContentPath: contentPath,
		PublishDate: publishDate,
	}
}

func (p *Post) Value() string {
	return fmt.Sprintf("%s %s %s %s", p.Title, p.Slug, p.GameId, p.ContentPath)
}

func GenerateId(slug string, prefix string) string {
	sum := generateHashSum(slug)
	return fmt.Sprintf("%s-%v", prefix, sum)
}

func generateHashSum(title string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(title))
	return h.Sum32()
}
