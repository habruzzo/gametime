package handlers

import (
	"encoding/json"
	"fmt"
	"gametime"
	"gametime/config"
	"gametime/db"
	"html/template"
	"net/http"
	"os"
	"strings"
)

const (
	APP  = "app"
	POST = "post"
	TOOL = "tool"
)

type headerTitles struct {
	Title     string `html:"title"`
	PageTitle string `html:"pageTitle"`
}

type HandlerMap struct {
	fmap map[string]HandlerFuncWrapper
}

type HandlerFuncWrapper struct {
	method      string
	path        string
	handlerFunc http.HandlerFunc
}

type Handler struct {
	name   string
	log    gametime.Logger
	cfg    *config.Config
	dgraph *db.Dgraph
	tmap   map[string]*template.Template
	hmap   *HandlerMap
}

func NewHandler(l gametime.Logger, cfg *config.Config, db *db.Dgraph) *Handler {
	l.Info("CREATING HANDLER", "handler")
	h := &Handler{
		name:   "handler",
		log:    l,
		cfg:    cfg,
		dgraph: db,
	}
	h.getViews()
	h.log.Info(h.tmap)
	h.getHandlerMap()
	h.log.Info(len(h.hmap.fmap))

	return h
}

func NewHandlerMap() *HandlerMap {
	fmap := make(map[string]HandlerFuncWrapper)
	return &HandlerMap{
		fmap: fmap,
	}
}

func (h *Handler) getHandlerMap() {
	handlers := make(map[string]HandlerFuncWrapper)
	hmap := &HandlerMap{fmap: handlers}
	h.hmap = hmap
}

func (h *Handler) GetHandlerMap() *HandlerMap {
	return h.hmap
}

func (h *HandlerMap) Get() map[string]HandlerFuncWrapper {
	return h.fmap
}

func (h *HandlerMap) Add(m map[string]HandlerFuncWrapper) {
	for i, v := range m {
		h.fmap[i] = v
	}
}

func (h *Handler) getViews() {
	var appFiles []string
	var postFiles []string
	var toolFiles []string
	files, err := os.ReadDir("./views")
	h.log.Info("parsing views")
	if err != nil {
		h.log.Error(err)
	}
	h.log.Info("parsing views")
	for _, file := range files {
		filename := file.Name()
		if strings.HasSuffix(filename, ".html") {
			appFiles = append(appFiles, "./views/"+filename)
			postFiles = append(postFiles, "./views/"+filename)
			toolFiles = append(toolFiles, "./views/"+filename)
		}
		if file.IsDir() && file.Name() == "App" {
			aFiles, err := os.ReadDir("./views/App")
			if err != nil {
				h.log.Error(err)
				continue
			}
			for _, aFile := range aFiles {
				filename := aFile.Name()
				if strings.HasSuffix(filename, ".html") {
					appFiles = append(appFiles, "./views/App/"+filename)
				}
			}
		}
		if file.IsDir() && file.Name() == "Post" {
			pFiles, err := os.ReadDir("./views/Post")
			if err != nil {
				h.log.Error(err)
				continue
			}
			for _, pFile := range pFiles {
				filename := pFile.Name()
				if strings.HasSuffix(filename, ".html") {
					postFiles = append(postFiles, "./views/Post/"+filename)
				}
			}
		}
		if file.IsDir() && file.Name() == "Tool" {
			tFiles, err := os.ReadDir("./views/Tool")
			if err != nil {
				h.log.Error(err)
				continue
			}
			for _, tFile := range tFiles {
				filename := tFile.Name()
				if strings.HasSuffix(filename, ".html") {
					toolFiles = append(toolFiles, "./views/Tool/"+filename)
				}
			}
		}
	}

	tmap := make(map[string]*template.Template)
	at, err := template.New(APP).Funcs(template.FuncMap{
		"url": func(action string) string {
			return h.hmap.GetPath(action)
		}}).ParseFiles(appFiles...)
	h.log.Info("parsed the app files to app template")
	if err != nil || at == nil {
		h.log.Error(err)
		panic(fmt.Sprintf("error parsing app views!! %v", err))
	}
	pt, err := template.New(POST).Funcs(template.FuncMap{
		"url": func(action string) string {
			return h.hmap.GetPath(action)
		},
		"post_url": func(slug string) string {
			path := h.hmap.GetPath("Post.GetPostsByMostRecent")
			return fmt.Sprintf("%s%s", path, slug)
		},
	}).ParseFiles(postFiles...)
	h.log.Info("parsed the files to post template")
	if err != nil || pt == nil {
		h.log.Error(err)
		panic(fmt.Sprintf("error parsing post views!! %v", err))
	}
	tt, err := template.New(TOOL).Funcs(template.FuncMap{
		"url": func(action string) string {
			return h.hmap.GetPath(action)
		},
		"marshal": func(obj interface{}) string {
			b, _ := json.Marshal(obj)
			return string(b)
		},
	}).ParseFiles(toolFiles...)
	h.log.Info("parsed the files to tool template")
	if err != nil || pt == nil {
		h.log.Error(err)
		panic(fmt.Sprintf("error parsing tool views!! %v", err))
	}
	tmap[APP] = at
	tmap[POST] = pt
	tmap[TOOL] = tt
	h.tmap = tmap
}

func (h *HandlerMap) GetMethod(key string) string {
	if f, ok := h.fmap[key]; ok {
		return f.method
	} else {
		panic(fmt.Sprintf("no handler for that key!! %s", key))
	}
}

func (h *HandlerMap) GetPath(key string) string {
	if f, ok := h.fmap[key]; ok {
		return f.path
	} else {
		panic(fmt.Sprintf("no handler for that key!! %s", key))
	}
}

func (h *HandlerMap) GetFunc(key string) http.HandlerFunc {
	if f, ok := h.fmap[key]; ok {
		return f.handlerFunc
	} else {
		panic(fmt.Sprintf("no handler for that key!! %s", key))
	}
}

func (h *HandlerFuncWrapper) GetFunc() http.HandlerFunc {
	return h.handlerFunc
}

func (h *HandlerFuncWrapper) GetMethod() string {
	return h.method
}

func (h *HandlerFuncWrapper) GetPath() string {
	return h.path
}
