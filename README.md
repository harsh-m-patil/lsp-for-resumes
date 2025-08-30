# ATS Optimus Prime LSP (Go)

A minimal Language Server Protocol (LSP) implementation in Go tailored for resume text analysis. It demonstrates JSON-RPC framing, core LSP methods, and simple diagnostics, hovers, definitions, code actions, and completions.

## Features
- Initialize handshake with advertised capabilities
- DidOpen/DidChange document tracking and diagnostics
- Diagnostics:
  - Flags low-impact bullet lines (e.g., lines like "- something" with no numbers)
  - Detects overused action verbs and suggests varying language

## Project layout
- `main.go` — server entry; stdin/stdout JSON-RPC loop
- `rpc/` — message framing (Content-Length) and tests
- `lsp/` — LSP request/response types
- `analysis/` — document state, diagnostics, and simple features

## Requirements
- Go

## Build
```sh
# from repo root
go build -o main .
```

## Run with an editor
The server communicates over stdio. Configure your editor to spawn the built binary.

- Neovim (nvim-lspconfig):
```lua
require('lspconfig').optimus_prime.setup({
  cmd = { '/absolute/path/to/main' },
  name = 'ats-optimus-prime-lsp',
})
```

- VS Code: use an extension or custom client that launches the binary and connects via stdio.

## Supported LSP methods
- initialize
- textDocument/didOpen
- textDocument/didChange
- textDocument/hover
- textDocument/definition
- textDocument/codeAction
- textDocument/completion

Note: shutdown/exit and didClose are not yet implemented.

## Logging
By default logs to a file path in `main.go`. Adjust `getLogger` to log to stderr or use an env var for portability.

## Testing
```sh
go test ./...
```
