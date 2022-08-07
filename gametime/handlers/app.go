package handlers

import (
	"net/http"
)

const (
	AppIndex   = "App.Index"
	AppAbout   = "App.About"
	AppBacklog = "App.Backlog"
	AppContact = "App.Contact"
	AppFormat  = "App.Format"

	INDEX   = "/"
	ABOUT   = "/about"
	BACKLOG = "/backlog"
	CONTACT = "/contact"
	FORMAT  = "/format"
)

type App struct {
	*Handler
}

func NewApp(h *Handler) *App {
	h.log.Info("CREATING APP", h.name)
	h.name = "app"
	a := &App{
		Handler: h,
	}
	a.getHandlerMap()
	return a
}

func (c *App) Index() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.log.Info("index")
		if c.tmap == nil || len(c.tmap) < 1 {
			c.log.Error("error rendering templates for index 1")
			w.WriteHeader(500)
			return
		}
		index := c.tmap[APP].Lookup("Index.html")
		if index == nil {
			c.log.Error("error rendering templates for index 2")
			w.WriteHeader(500)
			return
		}

		index.ExecuteTemplate(w, "Index.html", headerTitles{Title: "Home", PageTitle: "home"})
		return
	}
}

func (c *App) About() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.log.Info("about")
		if c.tmap == nil {
			c.log.Error("error rendering templates for about")
			w.WriteHeader(500)
			return
		}
		about := c.tmap[APP].Lookup("About.html")
		if about == nil {
			c.log.Error("error rendering templates for about 2")
			w.WriteHeader(500)
			return
		}
		about.ExecuteTemplate(w, "About.html", headerTitles{Title: "About", PageTitle: "about holdon (aka me)"})
		return
	}
}

func (c *App) Format() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.log.Info("format")
		if c.tmap == nil {
			c.log.Error("error rendering templates for format")
			w.WriteHeader(500)
			return
		}
		format := c.tmap[APP].Lookup("Format.html")
		if format == nil {
			c.log.Error("error rendering templates for format 2")
			w.WriteHeader(500)
			return
		}
		format.ExecuteTemplate(w, "Format.html", headerTitles{Title: "Review Format", PageTitle: "format"})
		return
	}
}

func (c *App) Backlog() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.log.Info("backlog")
		if c.tmap == nil {
			c.log.Error("error rendering templates for backlog")
			w.WriteHeader(500)
			return
		}
		backlog := c.tmap[APP].Lookup("Backlog.html")
		if backlog == nil {
			c.log.Error("error rendering templates for backlog 2")
			w.WriteHeader(500)
			return
		}

		backlog.ExecuteTemplate(w, "Backlog.html", headerTitles{Title: "Backlog", PageTitle: "backlog"})
		return
	}
}

func (c *App) Contact() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.log.Info("contact")
		if c.tmap == nil {
			c.log.Error("error rendering templates for contact")
			w.WriteHeader(500)
			return
		}
		contact := c.tmap[APP].Lookup("Contact.html")
		if contact == nil {
			c.log.Error("error rendering templates for contact 2")
			w.WriteHeader(500)
			return
		}
		contact.ExecuteTemplate(w, "Contact.html", headerTitles{Title: "Contact", PageTitle: "faq + contact"})
		return
	}
}

func (c *App) getHandlerMap() {
	c.hmap.fmap[AppIndex] = HandlerFuncWrapper{
		method:      "GET",
		path:        INDEX,
		handlerFunc: c.About(),
	}
	c.hmap.fmap[AppAbout] = HandlerFuncWrapper{
		method:      "GET",
		path:        ABOUT,
		handlerFunc: c.About(),
	}
	c.hmap.fmap[AppBacklog] = HandlerFuncWrapper{
		method:      "GET",
		path:        BACKLOG,
		handlerFunc: c.Backlog(),
	}
	c.hmap.fmap[AppContact] = HandlerFuncWrapper{
		method:      "GET",
		path:        CONTACT,
		handlerFunc: c.Contact(),
	}
	c.hmap.fmap[AppFormat] = HandlerFuncWrapper{
		method:      "GET",
		path:        FORMAT,
		handlerFunc: c.Format(),
	}
}
