package handler

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

func (h *Handler) TestHome(c echo.Context) (err error) {
	return c.JSON(http.StatusOK,map[string]string{"message": "Hello World"})
}
func (h *Handler) TestStream(c echo.Context) (err error) {
	type (
		Geolocation struct {
			Altitude  float64
			Latitude  float64
			Longitude float64
		}
	)
	var (
		locations = []Geolocation{
			{-97, 37.819929, -122.478255},
			{1899, 39.096849, -120.032351},
			{2619, 37.865101, -119.538329},
			{42, 33.812092, -117.918974},
			{15, 37.77493, -122.419416},
		}
	)

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Response().WriteHeader(http.StatusOK)

	enc := json.NewEncoder(c.Response())
	for _, l := range locations {
		if err := enc.Encode(l); err != nil {
			return err
		}
		c.Response().Flush()
		time.Sleep(500 * time.Millisecond)
	}
	return nil
}
