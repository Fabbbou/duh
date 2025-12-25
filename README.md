# Duh


## Duh CLI Documentation

A simple dotfiles manager for shell aliases and environment exports.

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
