package provider

import (
	"context"

	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
)

// Config defines provider-level configuration for authenticating with Loopia.
type Config struct {
	Username string `pulumi:"username" required:"true"`                   // Loopia API username
	Password string `pulumi:"password" required:"true" provider:"secret"` // Loopia API password (secret)
	Endpoint string `pulumi:"endpoint"`                                   // Loopia API endpoint (optional, for testing)
}

// Annotate adds descriptions to provider config fields for documentation and codegen.
func (c *Config) Annotate(a infer.Annotator) {
	a.Describe(&c.Username, "Loopia API username.")
	a.Describe(&c.Password, "Loopia API password.")
	a.Describe(&c.Endpoint, "Loopia API endpoint (optional, for testing)")
}

// ClientFactory is a function that creates a Client from provider config.
type ClientFactory func(ctx context.Context, config Config) (Client, error)

// NewProvider constructs a Pulumi provider using the given client factory.
func NewProvider(factory ClientFactory) (p.Provider, error) {
	return infer.NewProviderBuilder().
		WithNamespace("loopia-dns").
		WithConfig(infer.Config(&Config{})).
		WithResources(
			infer.Resource(&DnsRecord{getClient: factory}),
		).
		WithModuleMap(map[tokens.ModuleName]tokens.ModuleName{
			"loopia-dns": "index",
		}).
		Build()
}
