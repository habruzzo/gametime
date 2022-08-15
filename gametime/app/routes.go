package app

import (
	"gametime"
	"gametime/handlers"
	"net/http"

	"github.com/go-chi/chi"
)

type Handler interface {
	GetHandlerMap() *handlers.HandlerMap
}

type ApiHandler Handler

type AppHandler Handler

type PostHandler Handler

type ToolHandler Handler

type HandlerWrapper struct {
	MasterMap *handlers.HandlerMap
}

func NewHandlerWrapper(log gametime.Logger, api ApiHandler, app AppHandler, post PostHandler, tool ToolHandler) *HandlerWrapper {
	log.Info("CREATING HANDLER WRAPPER")

	m := handlers.NewHandlerMap()
	m.Add(api.GetHandlerMap().Get())
	log.Info("api online")
	m.Add(app.GetHandlerMap().Get())
	log.Info("app online")
	m.Add(post.GetHandlerMap().Get())
	log.Info("post online")
	m.Add(tool.GetHandlerMap().Get())
	log.Info("tool online")
	return &HandlerWrapper{
		MasterMap: m,
	}
}

func NewRouter(log gametime.Logger, hw *HandlerWrapper) *chi.Mux {
	r := chi.NewRouter()
	for i, v := range hw.MasterMap.Get() {
		r.MethodFunc(v.GetMethod(), v.GetPath(), v.GetFunc())
		log.Info("registered path", "path", v.GetPath(), "key", i)
	}
	r.Handle("/*", http.FileServer(http.Dir("./public")))

	return r
}
