package app

import (
	"gametime"
	"gametime/handlers"
	"net/http"

	"github.com/go-chi/chi"
)

type ApiHandler interface {
	GetHandlerMap() *handlers.HandlerMap
}

type AppHandler interface {
	GetHandlerMap() *handlers.HandlerMap
}

type PostHandler interface {
	GetHandlerMap() *handlers.HandlerMap
}

type HandlerWrapper struct {
	MasterMap *handlers.HandlerMap
}

func NewHandlerWrapper(log gametime.Logger, api ApiHandler, app AppHandler, post PostHandler) *HandlerWrapper {
	log.Info("CREATING HANDLER WRAPPER")

	m := handlers.NewHandlerMap()
	m.Add(api.GetHandlerMap().Get())
	log.Info("api online")
	m.Add(app.GetHandlerMap().Get())
	log.Info("app online")
	m.Add(post.GetHandlerMap().Get())
	log.Info("post online")
	return &HandlerWrapper{
		MasterMap: m,
	}
}

func NewRouter(log gametime.Logger, hw *HandlerWrapper) *chi.Mux {
	r := chi.NewRouter()
	r.Handle("/public/*", http.StripPrefix("/public", http.FileServer(http.Dir("./public"))))
	for i, v := range hw.MasterMap.Get() {
		r.MethodFunc(v.GetMethod(), v.GetPath(), v.GetFunc())
		log.Info("registered path", "path", v.GetPath(), "key", i)
	}
	return r
}
