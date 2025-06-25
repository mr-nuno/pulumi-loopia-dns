# Pulumi Loopia DNS Provider

A Pulumi provider for managing DNS records in Loopia using the XML-RPC API. This project is a template for building custom Pulumi providers in Go, following best practices and the structure of the official Pulumi Go Provider configurable example.

## Features
- Full CRUD support for Loopia DNS records
- Secure configuration via Pulumi secrets
- Dependency-injected client for extensibility and testability
- Well-documented codebase for easy adaptation to other providers

## Project Structure

```
provider/
  provider.go      # Provider definition, config, and registration
  client.go        # Loopia API client and interface
  dns_record.go    # DNS record resource implementation
  subdomain.go     # (Example) Subdomain resource implementation
main.go            # Provider entry point
```

## Usage

1. **Build the provider:**
   ```sh
   go build -o bin/pulumi-resource-loopia-dns main.go
   ```
2. **Configure the provider in your Pulumi program:**
   ```yaml
   config:
     username: <your-username>
     password: <your-password>
     endpoint: https://api.loopia.se/RPCSERV
   ```
3. **Use the `DnsRecord` resource in your Pulumi code.**

## .NET SDK Generation and Provider Installation

- The .NET SDK is generated from the provider's schema and is located in `sdk/dotnet/`.
- To use the provider in your Pulumi .NET program, ensure the generated Go provider executable (e.g., `pulumi-resource-loopia-dns.exe`) is available in your `PATH` or copied to your Pulumi project's working directory.
- You can build the provider executable with:
  ```sh
  go build -o bin/pulumi-resource-loopia-dns main.go
  ```
- If you update the provider schema or resources, regenerate the .NET SDK by running (from the repo root):
  ```sh
  # Example: using Pulumi's codegen tools (adjust as needed for your setup)
  pulumi-gen-dotnet --out sdk/dotnet/ schema.json
  ```
- After building, you may need to copy the executable to your Pulumi project's directory or ensure it's discoverable by the Pulumi CLI.

## Example: Using DnsRecord in a Pulumi .NET Program

```csharp
using Pulumi;
using Pulumi.LoopiaDns;

class MyStack : Stack
{
    public MyStack()
    {
        // Provider configuration from Pulumi config
        var config = new Config();
        var username = config.Require("username");
        var password = config.RequireSecret("password");
        var endpoint = config.Get("endpoint") ?? "https://api.loopia.se/RPCSERV";

        var provider = new LoopiaDns.Provider("loopia", new ProviderArgs
        {
            Username = username,
            Password = password,
            Endpoint = endpoint,
        });

        var record = new DnsRecord("my-record", new DnsRecordArgs
        {
            Zone = "example.com",
            Name = "www",
            Type = "A",
            Value = "1.2.3.4",
            Ttl = 3600,
        }, new CustomResourceOptions { Provider = provider });
    }
}
```

## Example: Adding a New Resource (Subdomain)

Suppose you want to add a new resource for managing subdomains. You would:

1. **Create `provider/subdomain.go`** with a resource struct and CRUD methods, following the pattern in `dns_record.go`.
2. **Register the resource in `provider.go`:**
   ```go
   infer.Resource(&Subdomain{getClient: factory}),
   ```
3. **Use the new resource in your Pulumi .NET program:**
   ```csharp
   using Pulumi;
   using Pulumi.LoopiaDns;

   class MyStack : Stack
   {
       public MyStack()
       {
           var subdomain = new Subdomain("my-subdomain", new SubdomainArgs
           {
               Zone = "example.com",
               Name = "blog",
           });
       }
   }
   ```

## Extending This Provider
- To add new resources, create a new file in `provider/` and follow the pattern in `dns_record.go`.
- To support a different API, implement a new `Client` and inject it via the provider factory.
- All types and methods are documented for easy onboarding.

## Development
- Run `go mod tidy` to manage dependencies.
- Run `go build ./...` to check for errors.
- See `provider/` for code documentation and extension points.

## License
MIT
