package main

import (
	"context"
	"errors"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/aalpern/luminosity"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type LuminosityApiServer struct {
	CatalogPath []string
	Addr        string

	server   *echo.Echo
	catalogs map[string]string
}

func (s *LuminosityApiServer) CommandInitialize(cmd *cobra.Command) {
	cmd.Flags().StringArrayVarP(&s.CatalogPath, "catalog-path", "p", []string{},
		"List of directory roots for finding catalogs")
	cmd.MarkFlagRequired("catalog-paths")

	cmd.Flags().StringVarP(&s.Addr, "addr", "a", ":8000",
		"Listening address for Luminosity server")
}

var (
	errorNoCatalogPath = errors.New("Catalog path not set")
)

func (s *LuminosityApiServer) Start(ctx context.Context) error {
	if err := s.loadCatalogs(); err != nil {
		return err
	}

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/api/catalog", s.getCatalogList)
	e.GET("/api/catalog/:name", s.getCatalog)
	e.GET("/api/catalog/:name/stats", s.getCatalogStats)
	e.GET("/api/catalog/:name/sunburst", s.getCatalogSunburst)

	s.server = e
	log.WithFields(log.Fields{
		"action": "api_server",
		"status": "start",
		"addr":   s.Addr,
	}).Info()
	return s.server.Start(s.Addr)
}

func (s *LuminosityApiServer) loadCatalogs() error {
	if len(s.CatalogPath) == 0 {
		return errorNoCatalogPath
	}

	catalogs := make(map[string]string)
	paths := luminosity.FindCatalogs(s.CatalogPath...)
	for _, path := range paths {
		name := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
		catalogs[name] = path
	}
	s.catalogs = catalogs

	return nil
}

func (s *LuminosityApiServer) Stop() error {
	log.WithFields(log.Fields{
		"action": "api_server",
		"status": "stop",
	}).Info()
	if s.server != nil {
		return s.server.Shutdown(context.Background())
	}
	return nil
}

func (s *LuminosityApiServer) Kill() error {
	log.WithFields(log.Fields{
		"action": "api_server",
		"status": "kill",
	}).Info()
	if s.server != nil {
		return s.server.Close()
	}
	return nil
}

func (s *LuminosityApiServer) loadCatalog(name string) (*luminosity.Catalog, error) {
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

func (s *LuminosityApiServer) getCatalogList(ctx echo.Context) error {
	names := []string{}
	for n, _ := range s.catalogs {
		names = append(names, n)
	}
	return ctx.JSON(http.StatusOK, names)
}

func (s *LuminosityApiServer) getCatalog(ctx echo.Context) error {
	name := ctx.Param("name")
	if cat, err := s.loadCatalog(name); err != nil {
		return err
	} else {
		defer cat.Close()
		return ctx.JSON(http.StatusOK, cat)
	}
}

func (s *LuminosityApiServer) getCatalogStats(ctx echo.Context) error {
	name := ctx.Param("name")
	if cat, err := s.loadCatalog(name); err != nil {
		return err
	} else {
		defer cat.Close()
		stats, _ := cat.GetStats()
		return ctx.JSON(http.StatusOK, stats)
	}
}

func (s *LuminosityApiServer) getCatalogSunburst(ctx echo.Context) error {
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
