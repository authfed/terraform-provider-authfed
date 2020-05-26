package main

import (
	"github.com/authfed/terraform-provider-authfed/authfed"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: authfed.Provider})
}
