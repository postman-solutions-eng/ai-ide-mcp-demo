# AI IDE MCP Demo

Demo project for generating a Go API server from a Postman Collection using AI + MCP.

## The Goal

Use the Postman MCP server to read the **Platform Service Catalog** collection and generate a working Go + Chi API server that implements it.

## Postman Assets

- **Workspace:** Service Catalog - Devspace
- **Collection:** Platform Service Catalog (MCP Demo)
- **Environment:** Local (baseUrl: http://localhost:8080)

## Demo Flow

1. Start with this empty project (no Go code)
2. Point AI at the Postman collection via MCP
3. Generate Go server that implements the collection
4. Run and test locally
5. Verify against Postman tests

## After Generation

```bash
# Run the server
go run ./cmd/server

# Test health
curl http://localhost:8080/health

# Test with Postman CLI
postman collection run <collection-id> -e <environment-id>
```
