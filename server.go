package main

import (
	"movie-festival/helpers"
	"movie-festival/models"
	"movie-festival/routes"
	"net/http"
	"os"

	"github.com/labstack/echo/v4/middleware"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	errLoadingEnvFile := godotenv.Load()
	if errLoadingEnvFile != nil {
		helpers.HandleError("error loading the .env file", errLoadingEnvFile)
	}

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowCredentials: true,
		AllowMethods:     []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))

	models.ModelsDb = models.Connect(1)
	routes.Build(e)
	e.Logger.Fatal(e.Start(":" + os.Getenv("APP_PORT")))

}
