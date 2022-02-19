package routes

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/kurvaid/bbb-interface/internal/config"
	"github.com/kurvaid/bbb-interface/internal/handlers"
	"github.com/kurvaid/bbb-interface/internal/middlewares"
)

func SetupRoutes(app *fiber.App, conf *config.Model, hCl *http.Client) {
	// Built-in fiber middlewares
	app.Use(recover.New())
	// Use log file only in production
	switch conf.EnvIsProd {
	case true:
		fConf := logger.Config{
			Format:     "[${time}] ${status} | ${method} - ${latency} - ${ip} | ${path}\n",
			TimeFormat: "02-Jan-2006 15:04:05",
			Output:     conf.LogFile,
		}
		app.Use(logger.New(fConf))
	case false:
		app.Use(logger.New())
	}

	// This app's endpoints
	app.Post("/create",
		middlewares.Auth(conf),
		handlers.CreateMeeting(conf, hCl),
	)
	app.Post("/join",
		middlewares.Auth(conf),
		handlers.JoinMeeting(conf),
	)
	app.Post("/end",
		middlewares.Auth(conf),
		handlers.EndMeeting(conf, hCl),
	)
	app.Post("/is_run",
		middlewares.Auth(conf),
		handlers.IsRunning(conf, hCl),
	)
	app.Get("/callback/destroy", handlers.CallbackOnDestroy(conf, hCl))

	// Custom middlewares AFTER endpoints
	app.Use(handlers.DefaultRouteNotFound)
}
