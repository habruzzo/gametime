package handlers

import (
	"gametime"
	"net/http"
	"strings"
)

const (
	GameEditor = "Game.Editor"
	GameSave   = "Game.Save"

	JSONGAME     = "/json/game"
	JSONGAMESAVE = "/json/game/save"
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
			a.log.Error("error getting posts for game editor")
			w.WriteHeader(500)
			return
		}
		a.log.Info(game.ExecuteTemplate(w, "GameEditor.html", Data{headerTitles: headerTitles{Title: "Game Editor", PageTitle: "game editor"}, Games: games}))
		return
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
		a.log.Info("saving game", len(s))
		games := a.dgraph.GetGames()
		for i, v := range games {
			newStatus := s[v.Id]
			if v.Status.Name != newStatus.Name {
				gs := a.dgraph.GetStatusByName(newStatus.Name)
				v.Status = gs
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
	c.log.Info(r.PostForm)
	sMap := make(map[string]string)
	for k, v := range r.PostForm {
		c.log.Info(k, v)
		id := strings.TrimPrefix(k, "status_")
		sMap[id] = v[0]
	}
	return sMap
}

func (c *Tool) convertStatuses(formMap map[string]string) map[string]gametime.Status {
	sMap := make(map[string]gametime.Status)
	for k, v := range formMap {
		c.log.Info(k, v)
		s := gametime.ToStatus(v)
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
}
