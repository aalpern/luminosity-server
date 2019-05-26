package main

import (
	_ "expvar"
	"net/http"
	"time"

	"github.com/aalpern/go-metrics-charts"
	"github.com/aalpern/luminosity"
	"github.com/jawher/mow.cli"
	"github.com/rcrowley/go-metrics"
	"github.com/rcrowley/go-metrics/exp"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	CacheDir    string
	CatalogDirs []string
	Addr        string
}

func CmdStart(app *cli.Cli) {
	app.Command("start", "Launch the server", func(cmd *cli.Cmd) {
		var cfg Config

		cmd.Spec = "[--catalog-dir] [--cache-dir] [--addr]"

		cmd.StringsOptPtr(&cfg.CatalogDirs, "c catalog-dir", nil,
			"Paths to directories containing catalog files")
		cmd.StringOptPtr(&cfg.CacheDir, "cache-dir", ".",
			"Location for temporary storage")
		cmd.StringOptPtr(&cfg.Addr, "a addr", ":8080",
			"Listening address")

		cmd.Action = func() { server(&cfg) }

	})
}

// This is all just temporary code right now
func server(cfg *Config) {
	// Log config
	log.WithFields(log.Fields{
		"action":           "server",
		"status":           "start",
		"cfg_addr":         cfg.Addr,
		"cfg_cache_dir":    cfg.CacheDir,
		"cfg_catalog_dirs": cfg.CatalogDirs,
	}).Info("Server starting")

	// Fetch list of catalogs
	catalogPaths := luminosity.FindCatalogs(cfg.CatalogDirs...)
	log.WithFields(log.Fields{
		"action": "find_catalogs",
		"status": "ok",
		"count":  len(catalogPaths),
	}).Info("Found catalogs")

	for _, cat := range catalogPaths {
		log.WithFields(log.Fields{
			"action":  "find_catalogs",
			"catalog": cat,
		}).Debug()
	}

	// Set up system metrics
	r := metrics.NewRegistry()
	metrics.RegisterDebugGCStats(r)
	metrics.RegisterRuntimeMemStats(r)
	go metrics.CaptureDebugGCStats(r, time.Second*5)
	go metrics.CaptureRuntimeMemStats(r, time.Second*5)

	exp.Exp(r)
	metricscharts.Register()

	http.ListenAndServe(cfg.Addr, nil)

	log.WithFields(log.Fields{
		"action": "server",
		"status": "done",
	}).Info("Done")
}
