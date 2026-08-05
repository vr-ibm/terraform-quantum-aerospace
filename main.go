package main

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/vr-ibm/terraform-quantum-aerospace/internal/provider"
)

var version = "dev"

func main() {
	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/vr-ibm/quantum",
	}
	err := providerserver.Serve(context.Background(), provider.New(version), opts)
	if err != nil {
		log.Fatal(err)
	}
}
