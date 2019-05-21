package main

import (
	"github.com/labstack/gommon/log"
	"github.com/metallurgical/go-echo-boilerplate/config"
	"github.com/metallurgical/go-echo-boilerplate/routes"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Load config files
	config.AppNew()

	// Define API wrapper
	api := echo.New()
	api.Use(middleware.Logger())
	api.Use(middleware.Recover())
	// CORS middleware for API endpoint.
	api.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))
	routes.DefineApiRoute(api)

	// Define WEB wrapper
	web := echo.New()
	web.Use(middleware.Logger())
	web.Use(middleware.Recover())
	routes.DefineWebRoute(web)

	server := echo.New()
	server.Any("/*", func(c echo.Context) (err error) {
		req := c.Request()
		res := c.Response()

		if req.URL.Path[:4] == "/api" {
			api.ServeHTTP(res, req)
		} else {
			web.ServeHTTP(res, req)
		}
		return
	})

	// Start server to listen to port 1200
	server.Logger.Fatal(server.Start(":1200"))
}