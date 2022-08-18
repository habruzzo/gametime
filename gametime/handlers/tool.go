package handlers

import (
	"encoding/json"
	"gametime"
	"net/http"
	"os"
	"strings"
)

const (
	GameEditor = "Game.Editor"
	GameSave   = "Game.Save"
	GameWrite  = "Game.Write"

	JSONGAME      = "/json/game"
	JSONGAMESAVE  = "/json/game/save"
	JSONGAMEWRITE = "/json/game/write"
)

type Tool struct {
	*Handler
}

func NewTool(h *Handler) *Tool {
	h.log.Info("CREATING TOOL", h.name)
	h.name = "tool"
	a := &Tool{
		Handler: h,
	}
	a.getHandlerMap()
	return a
}

func (a *Tool) GameEditor() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type Data struct {
			headerTitles
			Games []gametime.Game
		}
		a.log.Info("start game editor")
		if a.tmap == nil {
			a.log.Error("error rendering templates for game editor")
			w.WriteHeader(500)
			return
		}
		game := a.tmap[TOOL].Lookup("GameEditor.html")
		if game == nil {
			a.log.Error("error rendering templates for game 2")
			w.WriteHeader(500)
			return
		}
		games := a.dgraph.GetGames()
		if games == nil {
			a.log.Error("error getting games for game editor")
			w.WriteHeader(500)
			return
		}
		a.log.Info(game.ExecuteTemplate(w, "GameEditor.html", Data{headerTitles: headerTitles{Title: "Game Editor", PageTitle: "game editor"}, Games: games}))
		return
	}
}

func (a *Tool) GameWrite() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		games := a.dgraph.GetGames()
		if games == nil {
			a.log.Error("error getting games for writing")
			w.WriteHeader(500)
			return
		}
		b, err := json.Marshal(games)
		if err != nil {
			a.log.Error("error marshaling for writing")
			w.WriteHeader(500)
			return
		}
		err = os.WriteFile("../reviews/games/json/dump-1.json", b, os.ModePerm)
		if err != nil {
			a.log.Error("error writing")
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(200)
	}
}

func (a *Tool) GameSave() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}
		s := a.convertStatuses(a.collectStatuses(r))
		d := a.convertDetails(a.collectDetails(r))
		a.log.Info("saving game", len(s))
		games := a.dgraph.GetGames()
		set := false
		for i, v := range games {
			newStatus := s[v.Id]
			newDetails := d[v.Id]
			if v.Status.Name != newStatus.Name {
				gs := a.dgraph.GetStatusByName(newStatus.Name)
				v.Status = gs
				set = true
			}
			if v.Details != newDetails {
				v.Details = newDetails
				set = true
			}
			if set {
				games[i] = v
				err := a.dgraph.UpdateGame(games[i])
				if err != nil {
					a.log.Error("error updating game", "error", err, "title", v.Title)
					w.WriteHeader(500)
					w.Write([]byte(err.Error()))
					return
				}
			}

		}
		a.log.Info("leaving saving games", len(s))
		http.Redirect(w, r, "/json/game", 302)
	}
}

func (c *Tool) collectStatuses(r *http.Request) map[string]string {
	sMap := make(map[string]string)
	for k, v := range r.PostForm {
		id := strings.TrimPrefix(k, "status_")
		sMap[id] = v[0]
	}
	return sMap
}

func (c *Tool) convertStatuses(formMap map[string]string) map[string]gametime.Status {
	sMap := make(map[string]gametime.Status)
	for k, v := range formMap {
		s := gametime.ToStatus(v)
		sMap[k] = s
	}
	return sMap
}

func (c *Tool) collectDetails(r *http.Request) map[string]string {
	sMap := make(map[string]string)
	for k, v := range r.PostForm {
		id := strings.TrimPrefix(k, "details_")
		sMap[id] = v[0]
	}
	return sMap
}

func (c *Tool) convertDetails(formMap map[string]string) map[string]gametime.Details {
	sMap := make(map[string]gametime.Details)
	for k, v := range formMap {
		s := gametime.ToDetails(v)
		sMap[k] = s
	}
	return sMap
}

func (c *Tool) getHandlerMap() {
	c.hmap.fmap[GameEditor] = HandlerFuncWrapper{
		method:      "GET",
		path:        JSONGAME,
		handlerFunc: c.GameEditor(),
	}
	c.hmap.fmap[GameSave] = HandlerFuncWrapper{
		method:      "POST",
		path:        JSONGAMESAVE,
		handlerFunc: c.GameSave(),
	}
	c.hmap.fmap[GameWrite] = HandlerFuncWrapper{
		method:      "POST",
		path:        JSONGAMEWRITE,
		handlerFunc: c.GameWrite(),
	}
}
