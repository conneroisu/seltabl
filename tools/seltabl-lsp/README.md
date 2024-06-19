# seltabl-lsp

This is a [language server](https://microsoft.github.io/language-server-protocol/) for [seltabl](https://github.com/conneroisu/seltabl).

## Installation

```sh
go intall github.com/conneroisu/seltabl/seltabl-lsp
```

## Usage

### Neovim

```lua
---@diagnostic disable-next-line: missing-fields
local client = vim.lsp.start {
	name = "tools",
	cmd = { "path to seltabl-lsp binary" },
	on_attach = require("lsp_attach").on_attach,
}

if not client then
	vim.notify("Failed to start seltabl-lsp")
	print("Failed to start seltabl-lsp")
	return
end


vim.api.nvim_create_autocmd("FileType", {
	pattern = "go",
	callback = function()
		local bufnr = vim.api.nvim_get_current_buf()
		vim.lsp.buf_attach_client(bufnr, client)
	end
})
```

## Manual Usage

```sh
seltabl-lsp
```

## Development

### Run tests

```sh
go test ./...
```
