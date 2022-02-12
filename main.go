package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/kurvaid/bbb-interface/internal/config"
	"github.com/kurvaid/bbb-interface/internal/handlers"
	"github.com/kurvaid/bbb-interface/internal/logger"
	"github.com/kurvaid/bbb-interface/internal/routes"
)

func main() {
	f, err := os.ReadFile("app-config.yml")
	if err != nil {
		log.Fatalln("failed to read config file:", err)
	}

	var appConfig config.Model
	app, err := setup(&appConfig, bytes.NewReader(f))
	if err != nil {
		log.Fatalln("failed setup the app:", err)
	}

	// init custom app logger
	appConfig.LogFile, err = os.OpenFile(appConfig.LogDir+"log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0770)
	if err != nil {
		log.Fatalln("failed to open|create log file:", err)
	}

	cl := &http.Client{}
	routes.SetupRoutes(app, &appConfig, cl)

	logger.InfL.Printf("listening on %s:%v\n", appConfig.Host, appConfig.PortNum)
	logger.ErrL.Fatalln(app.Listen(fmt.Sprintf("%s:%v", appConfig.Host, appConfig.PortNum)))
}

// setup prepare everything that necessary before starting this app.
func setup(conf *config.Model, fBuf io.Reader) (*fiber.App, error) {
	// init and load the config file.
	newConf, err := config.NewConfig(fBuf)
	if err != nil {
		return nil, fmt.Errorf("failed to load config file: %v\n", err)
	}
	*conf = *newConf
	if err := conf.Sanitization(); err != nil {
		return nil, fmt.Errorf("failed sanitizing config: %v\n", err)
	}
	conf.SanitizationLog()
	if err := conf.BBB.Sanitization(); err != nil {
		return nil, fmt.Errorf("failed sanitizing BBB config: %v\n", err)
	}

	// Init internal logging.
	if err := logger.InitLogger(conf); err != nil {
		return nil, fmt.Errorf("failed to init internal logging: %v\n", err)
	}

	// if app in production use hostname from Nginx instead.
	var proxyHeader string
	if conf.EnvIsProd {
		proxyHeader = "X-Real-Ip"
	}

	app := fiber.New(fiber.Config{
		DisableStartupMessage: conf.EnvIsProd,
		ErrorHandler:          handlers.DefaultError,
		ProxyHeader:           proxyHeader,
	})

	return app, nil
}
