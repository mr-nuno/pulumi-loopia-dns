package pkg

import (
	"context"

	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi-go-provider/middleware/schema"
)

const (
	Name    = "loopia-dns"
	Version = "0.1.0"
)

type ProviderConfig struct {
	Username string `pulumi:"username" required:"true"`
	Password string `pulumi:"password" required:"true"`
	Endpoint string `pulumi:"endpoint"`
}

func Main() {
	p.RunProvider(context.Background(), Name, Version, NewProvider())
}

// This provider uses the `pulumi-go-provider` library to produce a code-first provider definition.
func NewProvider() p.Provider {
	return infer.Provider(infer.Options{
		// This is the metadata for the provider
		Metadata: schema.Metadata{
			DisplayName: "Loopia DNS",
			Description: "The Pulumi Loopia DNS Provider enables you to manage DNS records in Loopia using the XML-RPC API. It supports full CRUD operations for DNS records and is designed for secure, automated DNS management in your Pulumi infrastructure projects.",
			LogoURL:     "https://raw.githubusercontent.com/pulumi/pulumi-command/master/assets/logo.svg",
		},
		Config: infer.Config(ProviderConfig{}),
		Resources: []infer.InferredResource{
			infer.Resource(&DnsRecord{}),
		},
	})
}
