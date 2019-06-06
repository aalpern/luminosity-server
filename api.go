package main

import (
	"net/http"
	"net/url"

	"github.com/aalpern/luminosity"
	"github.com/labstack/echo/v4"
)

type APICatalog struct {
	Data *luminosity.Catalog `json:"data,omitempty"`
	Name string              `json:"name"`
	Href string              `json:"href"`
}

func NewAPICatalog(name string, cat *luminosity.Catalog) *APICatalog {
	return &APICatalog{
		Name: name,
		Href: "/api/catalog/" + url.QueryEscape(name),
		Data: cat,
	}
}

func (s *LuminosityServer) loadCatalog(name string) (*luminosity.Catalog, error) {
	if path, ok := s.catalogs[name]; !ok {
		return nil, echo.NewHTTPError(http.StatusNotFound, "Catalog not found")
	} else {
		if cat, err := luminosity.OpenCatalog(path); err != nil {
			return nil, err
		} else {
			return cat, cat.Load()
		}
	}
}

func (s *LuminosityServer) getCatalogList(ctx echo.Context) error {
	names := []*APICatalog{}
	for n, _ := range s.catalogs {
		names = append(names, NewAPICatalog(n, nil))
	}
	return ctx.JSON(http.StatusOK, names)
}

func (s *LuminosityServer) getCatalog(ctx echo.Context) error {
	name := ctx.Param("name")
	if cat, err := s.loadCatalog(name); err != nil {
		return err
	} else {
		defer cat.Close()
		return ctx.JSON(http.StatusOK, NewAPICatalog(name, cat))
	}
}

func (s *LuminosityServer) getCatalogStats(ctx echo.Context) error {
	name := ctx.Param("name")
	if cat, err := s.loadCatalog(name); err != nil {
		return err
	} else {
		defer cat.Close()
		stats, _ := cat.GetStats()

		return ctx.JSON(http.StatusOK, stats)
	}
}

func (s *LuminosityServer) getCatalogSunburst(ctx echo.Context) error {
	name := ctx.Param("name")
	if cat, err := s.loadCatalog(name); err != nil {
		return err
	} else {
		defer cat.Close()
		if data, err := cat.GetSunburstStats(); err != nil {
			return err
		} else {
			return ctx.JSON(http.StatusOK, data)
		}
	}
}
