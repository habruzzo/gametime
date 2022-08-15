package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type details struct {
	Text string `json:"text"`
	Date string `json:"installed"`
}

type game struct {
	Details details `json:"details"`
	Status  status  `json:"status"`
	Title   string  `json:"title"`
}

type status struct {
	Name string `json:"name"`
}

func main() {
	//read tsv files
	b, err := os.ReadFile("../reviews/games/csv/steam_licenses_080722.csv")
	if err != nil {
		panic("panic reading file")
	}
	lines := strings.Split(string(b), "\n")
	//create game list

	games := []game{}
	for _, v := range lines {
		fmt.Println(v)
		tmp := strings.Split(v, " \t")
		games = append(games, game{Details: details{Date: tmp[0]}, Title: tmp[1]})
	}

	//dump to json
	gb, err := json.Marshal(games)
	if err != nil {
		panic("panic marshalling")
	}
	err = os.WriteFile("../reviews/games/json/steam_licenses_080722.json", gb, os.ModePerm)
	if err != nil {
		panic("panic writing")
	}
	fmt.Println("done loading, wrote to file")
	return
}
