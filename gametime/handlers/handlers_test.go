package handlers

import (
	"gametime/config"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestRender(t *testing.T) {
	t.Run("test_render_index", testRenderIndex)
}

func testRenderIndex(t *testing.T) {
	log := logrus.New()
	cfg := config.Config{}

	handler := NewHandler(log, &cfg, nil)
	app := NewApp(handler)
	app.Index()
}
