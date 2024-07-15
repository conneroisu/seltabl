# seltabls

This is a [language server](https://microsoft.github.io/language-server-protocol/) for [seltabl](https://github.com/conneroisu/seltabl).

## Installation

```sh
go install github.com/conneroisu/seltabl/tools/seltabls@latest
```

## Usage

### Neovim

```lua
---@diagnostic disable-next-line: missing-fields
local client = vim.lsp.start {
	name = "seltabls",
	cmd = { "seltabls", "lsp" },
	on_attach = require("lsp_attach").on_attach,
}

if not client then
	vim.notify("Failed to start seltabls")
	print("Failed to start seltabls")
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

Base Command:
```sh
seltabls
```
Generate a new seltabl struct:
```sh
seltabls generate
```
Get the configuration for the cwd:
```sh
seltabls config
```

## Development

### Run tests

Using makefile (generally used for CI) one can run the tests with:

```sh
make test
```
Using taskfile (generally used for local development), one can run the tests with:

```sh
task test
```
Manually, one can run the tests with:

```sh
go test ./...
```

