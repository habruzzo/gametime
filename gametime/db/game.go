package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gametime"
	"github.com/dgraph-io/dgo/v200/protos/api"
)

func (d *Dgraph) InsertGame(g gametime.Game) error {
	foundG := d.getGameByTitle(g.Title)
	if foundG.Uid != "" {
		d.log.Info("game already exists", "title", g.Title)
		return errors.New("game already exists")
	}
	dg := d.gameToDgame(g)
	return d.insertGame(dg)
}

func (d *Dgraph) insertGame(da Dgame) error {
	d.newClient(d.cfg.Dgraph.Url)
	d.log.Info("insert game")
	if da.Status.Name != "" {
		ds := d.getStatusByName(da.Status.Name)
		da.Status = ds
	} else {
		ds := d.getStatusByName(gametime.Unknown.Name)
		da.Status = ds
	}
	b, err := json.Marshal(da)
	mu := api.Mutation{
		SetJson:   b,
		CommitNow: true,
	}
	d.log.Info(string(mu.SetJson))
	resp, err := d.Dgraph.NewTxn().Mutate(context.Background(), &mu)
	if err != nil {
		d.log.Error("error inserting game", "game", da.GameTitle)
		return err
	}
	d.log.Info("", "uids", resp.Uids)
	return nil
}

func (d *Dgraph) dGameToGame(da Dgame) gametime.Game {
	de := gametime.Details{}
	err := json.Unmarshal([]byte(da.Text), &de)
	if err != nil {
		d.log.Error("error populating game details", "error", err, "title", da.GameTitle)
	}
	return gametime.Game{
		Id:      da.Uid,
		Title:   da.GameTitle,
		Details: de,
		Status:  d.dstatusToStatus(da.Status),
	}
}

func (d *Dgraph) gameToDgame(a gametime.Game) Dgame {
	b, err := json.Marshal(a.Details)
	if err != nil {
		d.log.Error("error flattening game details", "error", err, "title", a.Title)
	}
	return Dgame{
		Uid:        a.Id,
		GameTitle:  a.Title,
		Status:     d.statusToDstatus(a.Status),
		Text:       string(b),
		DgraphType: gameType,
	}
}

func (d *Dgraph) GetGames() []gametime.Game {
	g := []gametime.Game{}
	dg := d.getGames()
	for _, v := range dg {
		g = append(g, d.dGameToGame(v))
	}
	return g
}

func (d *Dgraph) getGames() []Dgame {
	d.newClient(d.cfg.Dgraph.Url)
	d.log.Info("get games")
	query := `{games(func: has(gameTitle)) {
		uid
		gameTitle
		gameDetailsText
		gameStatus {
			uid
			statusName
			dgraph.type
		}
		dgraph.type
	}}`
	//fmt.Println(query)
	resp, err := d.Dgraph.NewReadOnlyTxn().Query(context.Background(), query)
	if err != nil {
		d.log.Error("error getting games", "error", err)
		return []Dgame{}
	}
	type games struct {
		Dgame []Dgame `json:"games"`
	}
	var dgame games
	err = json.Unmarshal(resp.Json, &dgame)
	if err != nil {
		d.log.Error("error unmarshalling game", "error", err)
		return []Dgame{}
	}
	if len(dgame.Dgame) < 1 {
		d.log.Error("error no game results in db for get games")
		return []Dgame{}
	}
	return dgame.Dgame
}

//func (d *Dgraph) GetGamesForBacklog() []gametime.Game {
//	d.log.Info("get games for backlog (app)")
//
//	dg := d.getGamesForBacklog()
//
//}

//func (d *Dgraph) getGamesForBacklog() []Dgame {
//	d.log.Info("get games for backlog")
//
//}

func (d *Dgraph) GetGameByTitle(Title string) gametime.Game {
	d.log.Info("get game by title (app)")

	da := d.getGameByTitle(Title)
	return d.dGameToGame(da)
}

func (d *Dgraph) getGameByTitle(title string) Dgame {
	d.newClient(d.cfg.Dgraph.Url)
	d.log.Info("get game by title")
	queryFmt := `{gameByTitle(func: eq(gameTitle,"%s")) {
		uid
		gameTitle
		gameDetailsText
		gameStatus {
			uid
			statusName
			dgraph.type
		}
		dgraph.type
	}}`
	query := fmt.Sprintf(queryFmt, title)
	//fmt.Println(query)
	resp, err := d.Dgraph.NewReadOnlyTxn().Query(context.Background(), query)
	if err != nil {
		d.log.Error("error getting game", "error", err)
		return d.gameToDgame(gametime.Game{Title: title})
	}
	type games struct {
		Dgame []Dgame `json:"gameByTitle"`
	}
	var dgame games
	err = json.Unmarshal(resp.Json, &dgame)
	if err != nil {
		d.log.Error("error unmarshalling game", "error", err)
		return d.gameToDgame(gametime.Game{Title: title})
	}
	if len(dgame.Dgame) < 1 {
		d.log.Error("error no game results in db", "title", title)
		return d.gameToDgame(gametime.Game{Title: title})
	}
	return dgame.Dgame[0]
}

func (d *Dgraph) UpdateGame(ug gametime.Game) error {
	d.newClient(d.cfg.Dgraph.Url)
	d.log.Info("update game")
	dg := d.gameToDgame(ug)
	return d.updateGame(dg)
}

func (d *Dgraph) updateGame(dg Dgame) error {
	b, err := json.Marshal(dg)
	if err != nil {
		d.log.Error("error marshalling game", "error", err, "title", dg.GameTitle)
		return err
	}
	mu := api.Mutation{
		SetJson:   b,
		CommitNow: true,
	}
	d.log.Info(string(mu.SetJson))
	resp, err := d.NewTxn().Mutate(context.Background(), &mu)
	if err != nil {
		d.log.Error("error updating game", "error", err, "title", dg.GameTitle)
		return err
	}
	d.log.Info("", "uids", resp.Uids)
	return nil
}
