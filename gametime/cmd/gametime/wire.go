//go:build wireinject
// +build wireinject

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

// Here you should add any Components that are necessary for running the services
var Component = wire.NewSet(
	// Bind the interface definition to the struct implementation
	app.NewServer,
	app.NewRunner,
	app.NewRouter,
	app.NewHandlerWrapper,

	config.LoadConfig,
	db.NewDgraph,
	handlers.NewApi,
	handlers.NewApp,
	handlers.NewPost,
	handlers.NewTool,
	handlers.NewHandler,

	logrus.New,
	// bind utils
	wire.Bind(new(gametime.Logger), new(*logrus.Logger)),

	// bind services
	wire.Bind(new(app.Router), new(*chi.Mux)),
	wire.Bind(new(app.ApiHandler), new(*handlers.Api)),
	wire.Bind(new(app.AppHandler), new(*handlers.App)),
	wire.Bind(new(app.PostHandler), new(*handlers.Post)),
	wire.Bind(new(app.ToolHandler), new(*handlers.Tool)),
)

func ProvideRunner() *app.Runner {
	panic(wire.Build(Component))
}

// if you want to add another service (WHICH YOU SHOULD DO WITH CAUTION),
// add another type struct to the services package, and add append the return type here
func ProvideServer() *app.Server {
	panic(wire.Build(Component))
}
