# Qng MCP Server Example

## use stdio mode
```bash
.\mcphost.exe  -m ollama:qwen2.5:3b --config .\conf\stdio.json --debug
```
## use sse mode

```bash
# start server
.\qng_server -t sse
.\mcphost.exe  -m ollama:qwen2.5:3b --config .\conf\stdio.json --debug
```