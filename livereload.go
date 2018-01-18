package livereload

import (
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
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
		Dir:     ".",
	}
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
					lrs.Reload(event.Name, filepath.Ext(event.Name) == ".css")
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
				return
			}
			if p == "/livereload" {
				lrs.ServeHTTP(c.Response(), c.Request())
				return
			}
			return next(c)
		}
	}
}
