> ⚠️
>
> **This project is still in ALPHA** (v0.x.x), meaning nothing is fully tested, and updates might be breaking.
>
> **No supports will be given until the Beta (v1.x.x), is launched**


# Duh
> 🔧 Simple dotfiles manager that actually works

Keep your shell aliases, shell functions, and git aliases, environment exports, synchronized across all your machines, without efforts.
You can finally have your alias everywhere easily.
Share your configs to your team, so you have the same alias and functions everywhere.



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

## Usage

```bash

# Aliases
duh alias set <name> <command>    # Set alias
duh alias unset <name>            # Remove alias  
duh alias list                    # List all

# Exports
duh exports set <var> <value>     # Set export
duh exports unset <var>           # Remove export
duh exports list                  # List all

# Repositories
duh repository list                           # List all repositories
duh repository enable <name>                  # Enable a repository
duh repository disable <name>                 # Disable a repository
duh repository delete <name>                  # Delete a repository
duh repository default                        # Show current default repository
duh repository default set <name>             # Set repository as default
duh repository rename <old> <new>             # Rename a repository
duh repository add <repo-url> [<custom name>] # Add a new repo from a remote git server
duh repository create <name>                  # Create new empty repository
duh repository update                         # Update repositories from remote sources
duh repository update --commit                # Update repositories, commit local changes first
duh repository update --force                 # Update repositories, discard local changes
duh repository edit <name>                    # Edit the export and aliases file for the given repo, using default editor

# Force duh to reload
duh_reload
```

> About `duh_reload`
>
> You can use this alias to force duh to reload from the config files
>
> Can be usefull when you just edited the configs, to avoid having to start a new shell

### Repository Updates

When working with repositories that have git remotes, you can update them:

- `duh repository update` - Safe update (fails if local changes exist)
- `duh repository update --commit` - Commits local changes before updating
- `duh repository update --force` - Discards local changes (destructive!)

### Example

```bash
duh alias set ll "ls -la"
duh exports set EDITOR "vim"
duh_reload
 
# Injects:
# alias ll="ls -la"
# export EDITOR="vim"
```
