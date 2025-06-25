package provider

import (
	"context"

	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
)

type Config struct {
	Username string `pulumi:"username" required:"true"`
	Password string `pulumi:"password" required:"true" provider:"secret"`
	Endpoint string `pulumi:"endpoint"`
}

func (c *Config) Annotate(a infer.Annotator) {
	a.Describe(&c.Username, "Loopia API username.")
	a.Describe(&c.Password, "Loopia API password.")
	a.Describe(&c.Endpoint, "Loopia API endpoint (optional, for testing)")
}

type ClientFactory func(ctx context.Context, config Config) (Client, error)

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
