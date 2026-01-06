> ⚠️ **This project is still in ALPHA** (v0.x.x) - features to add
>
> Some features might be broken, and updates might not support previous configs.
>
> **No user support will be given until the Beta (v1.x.x), is launched**
>
> Contributions are welcome !

# Duh
> 🔧 Simple Dotfiles (duh) CLI to handle and share shell aliases, exports, functions, and gitaliases that actually works

Keep your shell aliases, shell functions, and git aliases, environment exports, synchronized across all your machines, without efforts.

You can finally have your aliases and functions everywhere easily.

Share your configs to your team with the **repositories**, so you have the same alias and functions everywhere.



## Duh CLI Documentation
## Install
See [Releases page](https://github.com/Fabbbou/duh/releases)

### Setup

Add to your shell config (`~/.bashrc`, `~/.zshrc`):
```sh
eval "$(duh inject --quiet)"
```
>If duh has already been used before, you might have an alias command for this: `duh_reload`

### (optional) Final step: add your first repo

Add the `my-duh` repo as an example of my own

```sh
duh repo add https://github.com/Fabbbou/my-duh
```
It contains some of the aliases and git aliases that I use  

## Usage



```bash
# Aliases
duh alias set <name> <command>                # Set alias
duh alias unset <name>                        # Remove alias  
duh alias list                                # List all

# Exports
duh exports set <var> <value>                 # Set export
duh exports unset <var>                       # Remove export
duh exports list                              # List all

# Functions
duh functions list                            # List active functions
duh functions list --all                      # List all functions  
duh functions list --core                     # List internal core functions
duh functions info <function-name>            # Show details of a specific function
duh functions add <function-name>             # Create new function script to the default repository (opens editor)

# Repositories
duh repository list                           # List all repositories
duh repository enable <name>                  # Enable a repository to be injected by duh
duh repository disable <name>                 # Disable a repository, so it wont be injected
duh repository delete <name>                  # Delete a repository
duh repository default                        # Show current default repository
duh repository default set <name>             # Set repository as default
duh repository rename <old> <new>             # Rename a repository
duh repository add <repo-url> [<custom name>] # Add a new repo from a remote git server
duh repository create <name>                  # Create new empty repository
duh repository update                         # Update repositories from remote sources
duh repository update --commit                # Update repositories, commit local changes first
duh repository update --force                 # Update repositories, discard local changes
duh repository push <name>                    # Push local changes to remote repository
duh repository edit <name>                    # Edit the export and aliases file for the given repo, using default editor
duh repository edit-gitconfig <name>          # Create and/or Edit the <repo>/gitconfig file

# Self
duh self version                              # The detailed Duh build version
duh self config-path                          # Print configuration directory path  
duh self repositories-path                    # Print repositories directory path
duh self update                               # Update duh to the latest version

# Force duh to reload
duh_reload
```

### Example

```bash
duh alias set ll "ls -la"
duh exports set EDITOR "vim"
duh function add myfunction  # Creates and opens script for editing
duh_reload
 
# Injects:
# alias ll="ls -la"
# export EDITOR="vim"
# myfunction() { ... }  # (if function was added to the script)
```


## Roadmap for v1.0.0
- ~~Handling shell functions injection~~ ✅ **DONE**
- ~~Reworking architecture~~ ✅ **DONE**
- Waiting for feedbacks

