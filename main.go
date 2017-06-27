package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/namikoda/terraform-provider-namikoda/ipsfor"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: ipsfor.Provider})
}
