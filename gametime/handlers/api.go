package handlers

import (
	"encoding/json"
	"gametime"
	"gametime/db"
	"io"
	"net/http"
)

const (
	ApiInsertPost = "Api.InsertPost"
	ApiDump       = "Api.Dump"
	ApiHealth     = "Api.Health"

	INSERTPOST = "/insert"
	DUMP       = "/dump"
	HEALTH     = "/health"
)

type Api struct {
	*Handler
}

func NewApi(h *Handler) *Api {
	h.log.Info("CREATING API", h.name)
	h.name = "api"
	a := &Api{
		Handler: h,
	}

	a.getHandlerMap()
	return a
}

func (a Api) InsertPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a.log.Info("start insert post")
		if !a.checkAuth(r) {
			a.log.Error("error authing request")
			w.WriteHeader(401)
			return
		}
		var p gametime.Review
		var err error
		if p, err = a.bodyodyody(r); err != nil {
			a.log.Error("error parsing body")
			w.WriteHeader(400)
			return
		}
		a.dgraph = db.NewDgraph(a.log, a.cfg)
		if err = a.dgraph.InsertPost(p); err != nil {
			a.log.Error("error inserting post")
			w.WriteHeader(500)
			return
		}

		a.log.Info("finish insert post")
		w.WriteHeader(200)
		return
	}
}

func (a Api) Health() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}
}

func (a Api) Dump() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a.log.Info("start dump")
		if !a.checkAuth(r) {
			a.log.Error("error authing request")
			w.WriteHeader(401)
			return
		}
		if err := a.dgraph.Dump(); err != nil {
			a.log.Error("error dumping db")
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(200)
		a.log.Info("finish dump")
		return
	}
}

func (a Api) checkAuth(r *http.Request) bool {
	authHeader := r.Header.Get("Authorization")
	if authHeader != "cG90YXRv" {
		a.log.Error("no auth", authHeader)
		return false
	}
	return true
}

func (a Api) bodyodyody(r *http.Request) (gametime.Review, error) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		a.log.Error("bad request", "error", err)
		return gametime.Review{}, err
	}
	a.log.Info("req", len(b), string(b), r.Header.Get("Content-Length"), "_")
	var p gametime.Review
	err = json.Unmarshal(b, &p)
	if err != nil {
		a.log.Error("bad request", "error", err)
		return gametime.Review{}, err
	}
	p.Text = string(b)
	return p, nil
}

func (c *Api) getHandlerMap() {
	c.hmap.fmap[ApiInsertPost] = HandlerFuncWrapper{
		method:      "POST",
		path:        INSERTPOST,
		handlerFunc: c.InsertPost(),
	}
	c.hmap.fmap[ApiDump] = HandlerFuncWrapper{
		method:      "POST",
		path:        DUMP,
		handlerFunc: c.Dump(),
	}
	c.hmap.fmap[ApiHealth] = HandlerFuncWrapper{
		method:      "GET",
		path:        HEALTH,
		handlerFunc: c.Health(),
	}
}
