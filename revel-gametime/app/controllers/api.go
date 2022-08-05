package controllers

import (
	"encoding/json"
	"gametime/conf"
	"gametime/db"
	"io"

	"github.com/revel/revel"
)

type Api struct {
	DController
}

func (a Api) InsertPost() revel.Result {
	a.Log.Info("start insert post")
	if !a.checkAuth() {
		return a.Result
	}
	var p conf.Review
	var err error
	if p, err = a.bodyodyody(); err != nil {
		return a.Result
	}
	a.dgraph = db.NewDgraph(a.Log, Config)
	if err = a.dgraph.InsertPost(p); err != nil {
		a.DController.Response.SetStatus(500)
		return a.Result
	}
	a.DController.Response.SetStatus(200)

	a.Log.Info("finish insert post")
	return a.Result
}

func (a Api) Health() revel.Result {
	a.Response.SetStatus(200)
	return a.Result
}

func (a Api) Dump() revel.Result {
	a.Log.Info("start dump")
	if !a.checkAuth() {
		return a.Result
	}
	if err := a.dgraph.Dump(); err != nil {
		a.DController.Response.SetStatus(500)
		return a.Result
	}
	a.DController.Response.SetStatus(200)
	a.Log.Info("finish dump")
	return a.Result
}

func (a Api) checkAuth() bool {
	req := a.DController.Request
	authHeader := req.GetHttpHeader("Authorization")
	if authHeader != "cG90YXRv" {
		a.DController.Response.SetStatus(401)
		a.Log.Error("no auth", authHeader)
		return false
	}
	return true
}

func (a Api) bodyodyody() (conf.Review, error) {
	req := a.DController.Controller.Request
	b, err := io.ReadAll(req.GetBody())
	if err != nil {
		a.DController.Controller.Response.SetStatus(400)
		a.Log.Error("bad request", err.Error(), err)
		return conf.Review{}, err
	}
	a.Log.Info("req", len(b), string(b), req.GetHttpHeader("Content-Length"), "_")
	var p conf.Review
	err = json.Unmarshal(b, &p)
	if err != nil {
		a.DController.Controller.Response.SetStatus(400)
		a.Log.Error("bad request", err, "_")
		return conf.Review{}, err
	}
	p.Text = string(b)
	return p, nil
}
