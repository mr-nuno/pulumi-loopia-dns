package main

import (
	"context"
	"fmt"
	"os"

	"github.com/mr-nuno/pulumi-loopia-dns/provider"
)

// main is the entry point for the Pulumi Loopia DNS provider binary.
func main() {
	prov, err := provider.NewProvider(provider.RealClientFactory)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
		os.Exit(1)
	}
	err = prov.Run(context.Background(), "loopia-dns", "0.1.0")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
		os.Exit(1)
	}
}
