package main

import (
	"github.com/bartoszj/terraform-provider-mongodb/mongodb"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{ProviderFunc: mongodb.Provider})
}
