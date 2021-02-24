package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

const JsonPath = "/opt/gametime/reviews/"

type Question struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type ReviewSection struct {
	Title     string     `json:"title"`
	Questions []Question `json:"questions"`
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

type ReviewSkeleton struct {
	Overall      ReviewSection `json:"overall"`
	ArtSkeleton  ArtSection    `json:"art"`
	GameSkeleton GameSection   `json:"game"`
	Pull         string        `json:"pull"`
}

func NewReviewSkeleton(path string) *ReviewSkeleton {
	var r ReviewSkeleton
	fullPath := JsonPath + path
	data, err := ioutil.ReadFile(fullPath)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &r)
	if err != nil {
		panic(err)
	}
	return &r
}

func (r ReviewSkeleton) BuildPostToTemplate() {
	//r := NewReviewSkeleton(p.contentPath)

}

func (r ReviewSkeleton) Value() string {
	return fmt.Sprintf("%s", r.Pull)
}
