package main

import (
	"github.com/aalpern/svc"
	"github.com/aalpern/svc/httpsvc"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)

	global, _ := svc.NewCompositeComponent(
		svc.WithNamedComponent("log", &svc.LogConfigComponent{}),
		svc.WithShutdownWatcher())

	c, _ := svc.NewCompositeComponent(
		svc.WithNamedComponent("profile-server", &svc.ProfileServer{}),
		svc.WithNamedComponent("http-server", httpsvc.New()))

	svc.ServiceMain("luminosity-server", "Web access for Lightroom catalogs",
		svc.WithGlobal(global),
		svc.WithCommandHandler("start", "Start the API service", c))
}
