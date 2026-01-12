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

Share your configs to your team with the **packages**, so you have the same alias and functions everywhere.



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
duh package add https://github.com/Fabbbou/my-duh
```
It contains some of the aliases and git aliases that I use  

## Usage

Common usage
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
duh functions add <function-name>             # Create new function script to the default package (opens editor)
```

> Note: all the common commands above are editing the `default` package duh is pointing to
>
> You can change this default package with the commands bellow 


Configuring packages
```bash
# Packages
duh package default                        # Show current default package
duh package default set <name>             # Set package as default
duh package edit-gitconfig <name>          # Create and/or Edit the <package>/gitconfig file
duh package enable <name>                  # Enable a package to be injected by duh
duh package disable <name>                 # Disable a package, so it wont be injected
duh package delete <name>                  # Delete a package
duh package rename <old> <new>             # Rename a package
duh package add <package-url> [<custom name>] # Add a new package from a remote git server
duh package list                           # List all packages
duh package create <name>                  # Create new empty package
duh package update                         # Update packages from remote sources
duh package update --commit                # Update packages, commit local changes first
duh package update --force                 # Update packages, discard local changes
duh package push <name>                    # Push local changes to remote package
duh package edit <name>                    # Edit the export and aliases file for the given package, using default editor
```

Misc
```bash
# Self
duh self version                              # The detailed Duh build version
duh self config-path                          # Print configuration directory path  
duh self packages-path                       # Print packages directory path
duh self update                               # Update duh to the latest version

duh_reload                                    # Force duh to reload

duh inject                                    # every items injected by duh is printed in this command
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

