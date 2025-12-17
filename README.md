# AI IDE MCP Demo

Go + Chi API with Postman Collection testing via GitHub Actions.

## Project Structure

```
.
├── .github/workflows/     # CI/CD workflows
│   └── postman-tests.yml  # Newman collection runner
├── cmd/
│   └── api/
│       └── main.go        # Application entrypoint
├── internal/              # Private application code
├── postman/               # Postman collections & environments
│   ├── collection.json
│   └── environment.json
└── go.mod
```

## Development

```bash
# Run the server
go run ./cmd/api

# Run with hot reload (requires air)
air
```

## Testing with Postman/Newman

### Local Testing

```bash
# Install Newman
npm install -g newman

# Run collection
newman run ./postman/collection.json -e ./postman/environment.json
```

### CI/CD

The GitHub Actions workflow automatically runs the Postman collection on:
- Push to `main`/`master`
- Pull requests
- Manual trigger (workflow_dispatch)

## Environment Variables

| Variable | Description |
|----------|-------------|
| `API_BASE_URL` | Base URL for API requests |
| `API_KEY` | API authentication key (secret) |

Set these in GitHub repository settings:
- **Settings → Secrets and variables → Actions**
- Secrets: `API_KEY`
- Variables: `API_BASE_URL`


