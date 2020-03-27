package livereload

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/omeid/livereload"
)

type (
	LiveReloadConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper middleware.Skipper

		Name string
		Dir  string

		watcher *fsnotify.Watcher
	}
)

var (
	DefaultLiveReloadConfig = LiveReloadConfig{
		Skipper: middleware.DefaultSkipper,
		Name:    os.Args[0],
		Dir:     "assets",
	}
	upgrader = websocket.Upgrader{}
)

func LiveReload() echo.MiddlewareFunc {
	return LiveReloadWithConfig(DefaultLiveReloadConfig)
}

func LiveReloadWithConfig(config LiveReloadConfig) echo.MiddlewareFunc {
	lrs := livereload.New(config.Name)

	var err error
	config.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		panic("livereload: " + err.Error())
	}
	go func() {
		for {
			select {
			case event := <-config.watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					if rel, err := filepath.Rel(config.Dir, event.Name); err == nil {
						if !strings.HasPrefix(filepath.Base(rel), ".") {
							lrs.Reload("/"+filepath.ToSlash(rel), filepath.Ext(event.Name) == ".css")
						}
					}
				}
			}
		}
	}()

	err = config.watcher.Add(config.Dir)
	if err != nil {
		panic("livereload: " + err.Error())
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			if config.Skipper(c) {
				return next(c)
			}

			p := c.Path()
			if p == "/livereload.js" {
				livereload.LivereloadScript(c.Response(), c.Request())
			}
			if p == "/livereload" {
				lrs.ServeHTTP(c.Response(), c.Request())
			}
			return next(c)
		}
	}
}
