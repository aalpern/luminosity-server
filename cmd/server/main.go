package main

import (
	"os"

	"github.com/jawher/mow.cli"
	log "github.com/sirupsen/logrus"
)

func main() {
	app := cli.App("luminosity-server", "Experimental web UI for Lightroom catalogs")

	app.Spec = "[--verbose]"

	verbose := app.BoolOpt("v verbose", false,
		"Enable debug logging")

	app.Before = func() {
		if *verbose {
			log.SetLevel(log.DebugLevel)
		}
	}

	CmdStart(app)

	app.Run(os.Args)
}
