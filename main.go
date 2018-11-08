package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/ullbergm/terraform-provider-artifactory/artifactory"
)

func main() {
	opts := plugin.ServeOpts{
		ProviderFunc: artifactory.Provider,
	}
	plugin.Serve(&opts)
}
