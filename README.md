# Duh
Dotfiles (duh) CLI to handle and share shell aliases, exports, functions, and gitaliases. 

> ⚠️
>
> **This project is still in ALPHA** (v0.x.x), meaning nothing is fully tested, and updates might be breaking.
>
> **No supports will be given until the Beta (v1.x.x), is launched**

## Duh CLI Documentation
## Install

### Quick install
```sh
curl -sSL https://raw.githubusercontent.com/Fabbbou/duh/main/install.sh | sh
```

### Custom install directory
Just copy-paste this, make sure to update the `INSTALL_DIR` with your own.
```sh
INSTALL_DIR=$HOME/.local/bin curl -sSL https://raw.githubusercontent.com/Fabbbou/duh/main/install.sh | sh
```

### Setup

Add to your shell config (`~/.bashrc`, `~/.zshrc`):
```bash
eval "$(duh inject)"
```

### Commands

```bash
# Aliases
duh alias set <name> <command>    # Set alias
duh alias unset <name>            # Remove alias  
duh alias list                    # List all

# Exports
duh exports set <var> <value>     # Set export
duh exports unset <var>           # Remove export
duh exports list                  # List all

# Generate injection
duh inject                        # Output all aliases/exports
```

### Example

```bash
duh alias set ll "ls -la"
duh exports set EDITOR "vim"
duh inject
# Output:
# alias ll="ls -la"
# export EDITOR="vim"
```
