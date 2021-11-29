package api

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type IndexApi struct {
}

func (i *IndexApi) Welcome(c echo.Context) error {
	return c.String(http.StatusOK, "Welcome to my mapp")
}
