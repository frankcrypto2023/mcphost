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

## MCP 协议需要实现的路由

- server 处理 - https://github.com/mark3labs/mcp-go/blob/main/server/request_handler.go#L14

- 路由定义 - https://github.com/mark3labs/mcp-go/blob/main/mcp/types.go#L13 需要实现