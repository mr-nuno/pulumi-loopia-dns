# Loopia DNS Pulumi Provider

This project is based on the Pulumi provider boilerplate and the structure of the pulumi-command provider.

- `cmd/pulumi-resource-dns/`: Provider binary entrypoint
- `pkg/`: Provider logic and resources
- `tests/`: Provider-level tests

## Getting Started
1. Install Go 1.21 or later.
2. Build the provider:
   ```sh
   go build -o bin/pulumi-resource-dns ./cmd/pulumi-resource-dns
   ```
3. Implement your provider logic in the `pkg/` directory.

## Development
- Follow Pulumi provider best practices: https://www.pulumi.com/docs/guides/implement-providers/go/
- Implement resource CRUD logic in `pkg/`.
- Define schema in `internal/schema.go`.

## License
MIT
