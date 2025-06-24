# Loopia DNS Pulumi Provider

This project is the Pulumi provider for Loopia DNS, named `loopia-dns`.

- `main.go`: Provider binary entrypoint
- `pkg/`: Provider logic and resources
- `tests/`: Provider-level tests
- `sdk/`: Generated SDKs for supported languages

## Getting Started
1. Install Go 1.21 or later.
2. Build the provider:
   ```sh
   go build -o bin/pulumi-resource-dns ./cmd/pulumi-resource-dns
   ```
3. Implement your provider logic in the `pkg/` directory.

## SDK Generation
To generate SDKs for all supported languages from your schema, run:

```sh
pulumi package gen-sdk
```

This will use the `schema.json` in the project root and output SDKs to the `sdk/` directory. Make sure to delete any old generated SDKs before running this command to avoid stale files.

## Development
- Follow Pulumi provider best practices: https://www.pulumi.com/docs/guides/implement-providers/go/
- Implement resource CRUD logic in `pkg/`.
- Define your provider schema in `schema.json`.

## License
MIT
