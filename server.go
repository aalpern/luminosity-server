package main

import (
	"context"
	"errors"
	"path/filepath"
	"strings"

	"github.com/aalpern/luminosity"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type LuminosityServer struct {
	AssetsDir   string
	CatalogPath []string
	Addr        string

	server   *echo.Echo
	catalogs map[string]string
}

func (s *LuminosityServer) CommandInitialize(cmd *cobra.Command) {
	cmd.Flags().StringVar(&s.AssetsDir, "assets-dir", "static",
		"Path to the static assets directory for the Luminosity web UI")
	cmd.Flags().StringArrayVarP(&s.CatalogPath, "catalog-path", "p", []string{},
		"List of directory roots for finding catalogs")
	cmd.MarkFlagRequired("catalog-paths")
	cmd.Flags().StringVarP(&s.Addr, "addr", "a", ":8000",
		"Listening address for Luminosity server")
}

var (
	errorNoCatalogPath = errors.New("Catalog path not set")
)

func (s *LuminosityServer) Start(ctx context.Context) error {
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
	e.Static("/*", s.AssetsDir)

	s.server = e
	log.WithFields(log.Fields{
		"action":     "api_server",
		"status":     "start",
		"addr":       s.Addr,
		"static_dir": s.AssetsDir,
	}).Info()
	return s.server.Start(s.Addr)
}

func (s *LuminosityServer) loadCatalogs() error {
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

func (s *LuminosityServer) Stop() error {
	log.WithFields(log.Fields{
		"action": "api_server",
		"status": "stop",
	}).Info()
	if s.server != nil {
		return s.server.Shutdown(context.Background())
	}
	return nil
}

func (s *LuminosityServer) Kill() error {
	log.WithFields(log.Fields{
		"action": "api_server",
		"status": "kill",
	}).Info()
	if s.server != nil {
		return s.server.Close()
	}
	return nil
}
