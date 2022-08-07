// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"gametime"
	"gametime/app"
	"gametime/config"
	"gametime/db"
	"gametime/handlers"
	"github.com/go-chi/chi"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
)

// Injectors from wire.go:

func ProvideRunner() *app.Runner {
	logger := logrus.New()
	runner := app.NewRunner(logger)
	return runner
}

// if you want to add another service (WHICH YOU SHOULD DO WITH CAUTION),
// add another type struct to the services package, and add append the return type here
func ProvideServer() *app.Server {
	logger := logrus.New()
	configConfig := config.LoadConfig()
	dgraph := db.NewDgraph(logger, configConfig)
	handler := handlers.NewHandler(logger, configConfig, dgraph)
	api := handlers.NewApi(handler)
	handlersApp := handlers.NewApp(handler)
	post := handlers.NewPost(handler)
	handlerWrapper := app.NewHandlerWrapper(logger, api, handlersApp, post)
	mux := app.NewRouter(logger, handlerWrapper)
	server := app.NewServer(logger, mux, configConfig)
	return server
}

// wire.go:

// Here you should add any Components that are necessary for running the services
var Component = wire.NewSet(app.NewServer, app.NewRunner, app.NewRouter, app.NewHandlerWrapper, config.LoadConfig, db.NewDgraph, handlers.NewApi, handlers.NewApp, handlers.NewPost, handlers.NewHandler, logrus.New, wire.Bind(new(gametime.Logger), new(*logrus.Logger)), wire.Bind(new(app.Router), new(*chi.Mux)), wire.Bind(new(app.ApiHandler), new(*handlers.Api)), wire.Bind(new(app.AppHandler), new(*handlers.App)), wire.Bind(new(app.PostHandler), new(*handlers.Post)))