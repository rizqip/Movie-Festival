package routes

import (
	"log"
	"movie-festival/controllers"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var prefix = "/api/v1"
var Auth = middleware.JWT([]byte(os.Getenv("JWT_SECRET")))

func implContains(sl []string, name string) bool {
	// iterate over the array and compare given string to each element
	for _, value := range sl {
		if strings.Contains(name, value) {
			return true
		}
	}
	return false
}

// Build returns the application routes
func Build(e *echo.Echo) {
	e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
		Skipper: func(c echo.Context) bool {
			// Skip middleware if path is equal 'login'
			log.Println(c.Request().URL.Path)
			arrPath := []string{
				/* API Unused JWT */
				prefix + "/access/register",
				prefix + "/access/login",
			}

			var skip bool = true
			skip = implContains(arrPath, strings.ToLower(c.Request().URL.Path))
			return skip
		},
	}))


	//e.Use(echoMidleware.JWTWithConfig(DefaultJWTConfig))
	RouteGeneralApi(e)
}

func RouteGeneralApi(e *echo.Echo) {
	r := e.Group(prefix)
	/* Access */
	r.POST("/access/register", controllers.Register)
	r.POST("/access/login", controllers.Login)
	r.POST("/access/logout", controllers.Logout)

	/* Admin Only - Movie */
	r.POST("/admin/movie", controllers.CreateMovie)
	r.PUT("/admin/movie/:MovieId", controllers.UpdateMovie)
	/* Admin Only - Vote */
	r.GET("/admin/vote", controllers.MostVoted)
	/* Admin Only - View */
	r.GET("/admin/view", controllers.MostViewed)

	/* All User - Movie */
	r.GET("/general/movie", controllers.ListMovie)
	r.GET("/general/movie/:MovieId", controllers.DetailMovie)
	/* All User - View */
	r.POST("/general/view/:MovieId", controllers.ViewMovieRecord)
	r.GET("/general/view", controllers.RecordViewedMovie)
	/* All User - Vote */
	r.POST("/general/vote", controllers.Vote)
}