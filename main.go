package main

import (
	"echo-fullstack/app/router"
	"echo-fullstack/config"
	"echo-fullstack/modules/logger"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/tylerb/graceful"
)

func main() {
	godotenv.Load()

	app := echo.New()

	config.InitDB()
	config.Redis()
	logger.Setup(app)
	Log := logger.Log() //module logger
	defer Log.Sync()

	router.Init(app)
	port := os.Getenv("PORT")
	app.Server.Addr = ":" + port
	logger.Log().Info("Server restarted...at " + port)
	graceful.ListenAndServe(app.Server, 5*time.Second)
}
