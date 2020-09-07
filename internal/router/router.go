package router

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

// New returns a new echo router
func New(logger *logrus.Logger, config *Config) *echo.Echo {
	r := echo.New()

	r.Use(middleware.Recover())
	return r
}
