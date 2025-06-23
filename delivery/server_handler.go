package delivery

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	sha1ver   string
	buildTime string
	version   string
)

func (h *handler) getVersion(c echo.Context) error {

	r := h.svc.Version(c.Request().Context())

	return c.JSON(http.StatusOK, r)
}
