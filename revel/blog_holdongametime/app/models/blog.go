package models

import (
	"encoding/json"
	"io/ioutil"
)

const JsonPath = "/opt/gametime/reviews/"

type Question struct {
	question string
	answer   string
}

type ReviewSection struct {
	title     string
	questions []Question
}

type ArtSection struct {
	graphics ReviewSection
	sound    ReviewSection
	story    ReviewSection
	themes   ReviewSection
}

type GameSection struct {
	mechanics  ReviewSection
	difficulty ReviewSection
	experience ReviewSection
}

type ReviewSkeleton struct {
	overall      ReviewSection
	artSkeleton  ArtSection
	gameSkeleton GameSection
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
