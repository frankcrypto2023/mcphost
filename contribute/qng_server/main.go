package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

const RPC_URL = "http://216.230.226.189:1234/"

// JSONRPCRequest 结构体用于构建 JSON-RPC 请求
type JSONRPCRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

// JSONRPCResponse 结构体用于解析 JSON-RPC 响应
type JSONRPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Result  interface{} `json:"result"`
	Error   *RPCError   `json:"error"`
}

// RPCError 结构体用于解析 JSON-RPC 错误
type RPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func JsonRpcResponse(method string, params []interface{}) ([]byte, error) {

	// 构建 JSON-RPC 请求体
	request := JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
		ID:      1,
	}

	// 将请求体编码为 JSON
	requestBody, err := json.Marshal(request)
	if err != nil {
		fmt.Println("Error marshaling JSON request:", err)
		return nil, err
	}

	// 发送 HTTP POST 请求
	resp, err := http.Post(RPC_URL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("Error sending HTTP request:", err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, err
	}
	return body, nil

}

// authKey is a custom context key for storing the auth token.
type authKey struct{}

// withAuthKey adds an auth key to the context.
func withAuthKey(ctx context.Context, auth string) context.Context {
	return context.WithValue(ctx, authKey{}, auth)
}

// authFromRequest extracts the auth token from the request headers.
func authFromRequest(ctx context.Context, r *http.Request) context.Context {
	return withAuthKey(ctx, r.Header.Get("Authorization"))
}

// authFromEnv extracts the auth token from the environment
func authFromEnv(ctx context.Context) context.Context {
	return withAuthKey(ctx, os.Getenv("API_KEY"))
}

type response struct {
	Args    map[string]interface{} `json:"args"`
	Headers map[string]string      `json:"headers"`
}

// Block represents a block in the blockchain.
type Block struct {
	Order int
	Data  string // or any other fields relevant to your blockchain
}

// Blockchain represents the blockchain.
type Blockchain struct {
	Blocks []Block
}

type MCPServer struct {
	server *server.MCPServer
}

func NewMCPServer() *MCPServer {
	mcpServer := server.NewMCPServer(
		"example-server",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
		server.WithToolCapabilities(true),
	)

	mcpServer.AddTool(mcp.NewTool("get_block_by_order",
		mcp.WithDescription("Retrieves a qng block by its order"),
		mcp.WithNumber("order",
			mcp.Description("Order of the qng block to retrieve"),
			mcp.Required(),
		),
	), handleGetBlockByOrderTool)

	mcpServer.AddTool(mcp.NewTool("get_block_count",
		mcp.WithDescription("Retrieves a qng block total count"),
	), handleGetBlockCount)

	mcpServer.AddTool(mcp.NewTool("get_block_stateroot",
		mcp.WithDescription("Retrieves a qng block stateroot by its order"),
		mcp.WithNumber("order",
			mcp.Description("Order of the qng block stateroot to retrieve"),
			mcp.Required(),
		),
	), handleGetStateRoot)

	return &MCPServer{
		server: mcpServer,
	}
}

func (s *MCPServer) ServeSSE(addr string) *server.SSEServer {
	return server.NewSSEServer(s.server,
		server.WithBaseURL(fmt.Sprintf("http://%s", addr)),
		server.WithSSEContextFunc(authFromRequest),
	)
}

func (s *MCPServer) ServeStdio() error {
	return server.ServeStdio(s.server, server.WithStdioContextFunc(authFromEnv))
}

// handleGetBlockByOrderTool handles the qng_getBlockByOrder tool request.
func handleGetBlockByOrderTool(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	order, ok := request.Params.Arguments["order"]
	if !ok {
		return nil, fmt.Errorf("missing or invalid order")
	}
	body, err := JsonRpcResponse("qng_getBlockByOrder", []interface{}{order, true})
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(string(body)), nil
}

// handleGetBlockCount handles the qng_getBlockCount tool request.
func handleGetBlockCount(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	body, err := JsonRpcResponse("qng_getBlockCount", []interface{}{})
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(string(body)), nil
}

// handleGetStateRoot handles the qng_getStateRoot tool request.
func handleGetStateRoot(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	order, ok := request.Params.Arguments["order"]
	if !ok {
		return nil, fmt.Errorf("missing or invalid order")
	}
	body, err := JsonRpcResponse("qng_getStateRoot", []interface{}{order, true})
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(string(body)), nil
}

func main() {
	var transport string
	flag.StringVar(&transport, "t", "stdio", "Transport type (stdio or sse)")
	flag.StringVar(
		&transport,
		"transport",
		"stdio",
		"Transport type (stdio or sse)",
	)
	flag.Parse()

	s := NewMCPServer()

	switch transport {
	case "stdio":
		if err := s.ServeStdio(); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	case "sse":
		sseServer := s.ServeSSE("localhost:8080")
		log.Printf("SSE server listening on :8080")
		if err := sseServer.Start(":8080"); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	default:
		log.Fatalf(
			"Invalid transport type: %s. Must be 'stdio' or 'sse'",
			transport,
		)
	}
}
