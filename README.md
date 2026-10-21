# amixr-api-go-client

[![CI](https://github.com/grafana/amixr-api-go-client/actions/workflows/ci.yml/badge.svg)](https://github.com/grafana/amixr-api-go-client/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/grafana/amixr-api-go-client.svg)](https://pkg.go.dev/github.com/grafana/amixr-api-go-client)

Go client for the [Grafana OnCall public API](https://grafana.com/docs/grafana-cloud/observe-and-act/respond-to-incidents/reference/oncall-api/).

It is primarily maintained for the [Grafana Terraform provider](https://github.com/grafana/terraform-provider-grafana). The module is still `v0.x`, so the API may change between minor versions.

## Install

```sh
go get github.com/grafana/amixr-api-go-client
```

The package name is `aapi`.

## Usage

### Grafana Cloud (recommended)

Pass your Grafana stack URL and a Grafana service account token. The client looks up the OnCall API URL from the IRM plugin settings of your stack on the first request, so you don't need to know which OnCall region the stack uses.

```go
package main

import (
	"fmt"
	"log"

	aapi "github.com/grafana/amixr-api-go-client"
)

func main() {
	client, err := aapi.NewWithGrafanaAutodiscovery(
		"https://your-stack.grafana.net", // Grafana stack URL
		"<grafana-service-account-token>",
		"", // OnCall token; empty reuses the service account token
		"", // optional OnCall API URL, used only if the lookup fails
	)
	if err != nil {
		log.Fatal(err)
	}

	users, _, err := client.Users.ListUsers(&aapi.ListUserOptions{})
	if err != nil {
		log.Fatal(err)
	}
	for _, u := range users.Users {
		fmt.Println(u.ID, u.Username)
	}
}
```

### Explicit OnCall API URL

If you already know the OnCall API URL, set it directly:

```go
// Grafana Cloud: OnCall API URL, token, and Grafana stack URL.
client, err := aapi.NewWithGrafanaURL(
	"https://oncall-prod-us-central-0.grafana.net/oncall",
	"<grafana-service-account-token>",
	"https://your-stack.grafana.net",
)

// OnCall API URL and an OnCall API token, with no Grafana stack URL.
client, err := aapi.New("https://oncall.example.com", "<oncall-api-token>")
```

Never hard-code tokens. Read them from environment variables or a secret store.

### Pagination

List methods return one page at a time. Set `Page` on the options and follow `Next` until it is `nil`:

```go
opts := &aapi.ListIntegrationOptions{}
for page := 1; ; page++ {
	opts.Page = page
	resp, _, err := client.Integrations.ListIntegrations(opts)
	if err != nil {
		log.Fatal(err)
	}
	// use resp.Integrations
	if resp.Next == nil {
		break
	}
}
```

## Development

```sh
go test -race ./...
go vet ./...
gofmt -l .
```

CI runs the same checks, plus `go mod tidy`, on every pull request.

## Contributing

Issues and pull requests are welcome. Please include tests for new endpoints or behaviour changes.

To report a security issue, follow the [Grafana security policy](https://github.com/grafana/.github/blob/main/SECURITY.md). Do not open a public issue.

## License

[Apache License 2.0](LICENSE)
