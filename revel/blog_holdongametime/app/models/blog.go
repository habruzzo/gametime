package models

import (
	"io/ioutil"

	"gorm.io/gorm"
)

const JsonPath = "/opt/gametime/reviews/json/"

type Question struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type ReviewSection struct {
	Title     string     `json:"title"`
	Questions []Question `gorm:"foreignKey:ID"`
}

type ArtSection struct {
	Title    string        `json:"title"`
	Graphics ReviewSection `json:"graphics"`
	Sound    ReviewSection `json:"sound"`
	Story    ReviewSection `json:"story"`
	Themes   ReviewSection `json:"themes"`
}

type GameSection struct {
	Title      string        `json:"title"`
	Mechanics  ReviewSection `json:"mechanics"`
	Difficulty ReviewSection `json:"difficulty"`
	Experience ReviewSection `json:"experience"`
}

type OverallSection struct {
	Title   string        `json:"title"`
	Overall ReviewSection `json:"overall"`
}

type ReviewSkeleton struct {
	OverallSkeleton OverallSection `json:"overall"`
	ArtSkeleton     ArtSection     `json:"art"`
	GameSkeleton    GameSection    `json:"game"`
	Pull            string         `json:"pull"`
	Pics            []string       `json:"imgs"`
}

type Review struct {
	gorm.Model
	Id        string
	Post      *GormPost `gorm:"foreignKey:Id"`
	Path      string
	JsonBytes []byte
}

func NewReview(p *GormPost) *Review {
	id := GenerateId(p.Slug, "review")
	fullPath := JsonPath + p.ContentPath
	data, err := ioutil.ReadFile(fullPath)
	if err != nil {
		panic(err)
	}
	return &Review{
		Id:        id,
		Post:      p,
		Path:      p.ContentPath,
		JsonBytes: data,
	}

}
