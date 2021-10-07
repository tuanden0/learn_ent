package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func initRoutes(e *echo.Echo, srv Service) {

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	g := e.Group("/v1/items/")
	// Not check authen
	g.Add(http.MethodGet, "ping", srv.Ping)

	// Check authen
	g.Use(middleware.JWTWithConfig(configJWT))

	g.Add(http.MethodPost, "", srv.Create)
	g.Add(http.MethodGet, ":id", srv.Retrieve)
	g.Add(http.MethodPatch, ":id", srv.Update)
	g.Add(http.MethodDelete, ":id", srv.Delete)
}
