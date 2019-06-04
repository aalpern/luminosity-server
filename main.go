package main

import (
	. "github.com/aalpern/svc"
	"github.com/aalpern/svc/components"
)

func main() {
	ServiceMain(
		"luminosity-server",
		"Web access for Lightroom catalogs",
		WithGlobal(NewCompositeComponent(
			WithComponent(&components.LogConfigComponent{}),
			WithComponent(&components.RuntimeMetricsComponent{}),
			WithComponent(components.NewShutdownWatcher()))),
		WithCommandHandler("start", "Start the API service",
			NewCompositeComponent(
				WithComponent(&components.ProfileServer{
					Enable: true,
				}),
				WithComponent(&LuminosityApiServer{}))),
	)
}
