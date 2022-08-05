package controllers

import (
	"gametime/conf"
	"sync"

	"github.com/revel/revel"
)

var (
	// HOLDEN my own variable

	mu     sync.RWMutex
	Config conf.Config
)

func init() {
	revel.OnAppStart(InitConfig)
}

func InitConfig() {
	mu.Lock()
	defer mu.Unlock()
	Config = conf.LoadConfig()
}
