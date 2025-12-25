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
Just copy-paste this in your favorite shell.
```sh
curl -sSL https://raw.githubusercontent.com/Fabbbou/duh/main/install.sh | sh
```

### Custom install directory
Just copy-paste this, make sure to update the `INSTALL_DIR` with the duh executable path you'd like
```sh
INSTALL_DIR=$HOME/.local/bin curl -sSL https://raw.githubusercontent.com/Fabbbou/duh/main/install.sh | sh
```

### Setup

Add to your shell config (`~/.bashrc`, `~/.zshrc`):
```bash
eval "$(duh inject)"
```

>If duh has already been used before, you might have an alias command for this: `duh_reload`

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

# Force duh to reload
duh_reload
```

> About `duh_reload`
>
> You can use this alias to force duh to reload from the config files
>
> Can be usefull when you just edited the configs, to avoid having to start a new shell

### Example

```bash
duh alias set ll "ls -la"
duh exports set EDITOR "vim"
duh inject
# Output:
# alias ll="ls -la"
# export EDITOR="vim"
```
