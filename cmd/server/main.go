package main

import (
	"github.com/aalpern/svc"
	"github.com/aalpern/svc/components"
	"github.com/aalpern/svc/httpsvc"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)

	global, _ := svc.NewCompositeComponent(
		svc.WithComponent(&components.LogConfigComponent{}),
		svc.WithComponent(&components.RuntimeMetricsComponent{}),
		svc.WithComponent(components.NewShutdownWatcher()))

	c, _ := svc.NewCompositeComponent(
		svc.WithNamedComponent("profile-server", &components.ProfileServer{}),
		svc.WithNamedComponent("http-server", httpsvc.New()))

	svc.ServiceMain("luminosity-server", "Web access for Lightroom catalogs",
		svc.WithGlobal(global),
		svc.WithCommandHandler("start", "Start the API service", c))
}
