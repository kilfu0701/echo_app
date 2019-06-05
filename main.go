package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"

	"github.com/kilfu0701/echo_app/config"
	"github.com/kilfu0701/echo_app/core"
	"github.com/kilfu0701/echo_app/routes"
)

func init() {
	// TODO: Configで設定できるようにする
	// Setting timezone.
	location := "Asia/Tokyo"
	difference := 9 * 60 * 60

	loc, err := time.LoadLocation(location)
	if err != nil {
		loc = time.FixedZone(location, difference)
	}
	time.Local = loc
}

func main() {
	e := echo.New()

	e.Logger.SetLevel(log.INFO)
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	// load env.yaml
	cfg := config.Load("config/env.yaml")

	// extend by AppContext
	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ac, err := core.NewAppContext(c, cfg)
			if err != nil {
				return err
			}
			defer ac.AppDBClient.Disconnect(ac.AppDBCtx)
			return h(ac)
		}
	})

	// auth middleware
	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			log.Printf("===> %v", c.Request().Header.Get("X-User-Identifier"))
			return h(c)
		}
	})

	// load routes
	routes.Load(e)

	// start socket
	socket_file := "/tmp/echo_app.sock"
	os.Remove(socket_file)

	l, err := net.Listen("unix", socket_file)
	if err != nil {
		e.Logger.Fatal(err)
	}

	err = os.Chmod(socket_file, 0777)
	if err != nil {
		e.Logger.Fatal(err)
	}
	e.Listener = l
	// end socket

	go func() {
		if err := e.Start(""); err != nil {
			e.Logger.Info("shutting down the server")
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
