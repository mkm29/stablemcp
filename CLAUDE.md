# StableMCP Development Guidelines

## Commands
- Build: `go build -o bin/stablemcp ./main.go`
- Run server: `./bin/stablemcp server [--config path/to/config.yaml]`
- Tests: `go test ./...`
- Run single test: `go test ./path/to/package -run TestName`
- Format: `go fmt ./...`
- Lint: `golangci-lint run`

## Code Style Guidelines
- Follow standard Go conventions (https://go.dev/doc/effective_go)
- Use snake_case for file names and package names
- Use camelCase for variable names and PascalCase for exported functions/types
- Group imports: stdlib first, then third-party, then internal packages
- Use explicit error handling with informative messages
- Prefer composition over inheritance
- Use context for propagating cancellation and timeouts
- Prefer dependency injection for testability
- Document all exported functions and types
- Keep functions short and focused on a single responsibility