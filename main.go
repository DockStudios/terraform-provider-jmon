package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/matthewjohn/jmon-terraform-provider/jmon"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: jmon.Provider})
}
