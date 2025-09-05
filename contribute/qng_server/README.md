# Qng MCP Server Example

## use stdio mode
```bash
.\mcphost.exe  -m ollama:qwen2.5:3b --config .\conf\stdio.json --debug
```
## use sse mode

```bash
# start server
.\qng_server -t sse
.\mcphost.exe  -m ollama:qwen2.5:3b --config .\conf\sse.json --debug
```

## 执行结果
```bash
2025/04/26 09:21:20 INFO Model loaded provider=ollama model=qwen2.5:3b
2025/04/26 09:21:20 INFO Initializing server... name=qngserver
2025/04/26 09:21:20 INFO Server connected name=qngserver
2025/04/26 09:21:20 INFO Tools loaded server=qngserver count=3

  You: qng 最新区块总数
2025/04/26 09:21:36 INFO 🔧 Using tool name=qngserver__get_block_count

  Assistant:


  The total number of qNG blocks in the latest zone is 9,590,345. Is there anything else you need to know about this information?



  You: qng 在 1590344 order的区块包含的交易数量
2025/04/26 09:22:21 INFO 🔧 Using tool name=qngserver__get_block_by_order

  Assistant:


  The transaction count in the qNG block with order 1590344 is 1. This block contains only one transaction, which is a coinbase transaction that includes an asset transfer of 120,000,000,000 MEER tokens. Is there anything else you would like to know about this
  block?



  You: qng 在 1590344 order 的stateroot 详细信息
2025/04/26 09:23:00 INFO 🔧 Using tool name=qngserver__get_block_stateroot

  Assistant:


  For the qNG block with order 1590344, the stateroot details are as follows:

  • Hash: c5c045912d19ea72ec84dab22772a5307b37983154bd494f0cf94b402cee2e32
  • Order: 1590344
  • Height: 1584355
  • Valid: True
  • EVM State Root: 0xadd9acd10c9fd3a52956fc03d202867bbb8d19a5e8649e9047e751c4e6e444ce
  • EVM Height: 95296
  • EVM Head: 0x35c13f5238f51511f85d8c63f3a71683fdcd3351c14b3748fb6aaa5e1d827038
  • State Root: 920bf0c285558944cae4b8f4745a2453235d714fb3c6e5d3ced7b05ada08ee52

  Is there anything specific you would like to know about these details?
```

## MCP 协议需要实现的路由

- server 处理 - https://github.com/mark3labs/mcp-go/blob/main/server/request_handler.go#L14

- 路由定义 - https://github.com/mark3labs/mcp-go/blob/main/mcp/types.go#L13 需要实现