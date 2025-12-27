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

You can finally have your alias everywhere easily.

Share your configs to your team with the **repositories**, so you have the same alias and functions everywhere.



## Duh CLI Documentation
## Install

### Quick install
Just copy-paste this in your favorite shell.
```sh
curl -sSL https://raw.githubusercontent.com/Fabbbou/duh/main/install.sh | sh
```

#### Custom install directory
Just copy-paste this, make sure to update the `INSTALL_DIR` with the duh executable path you'd like
```sh
INSTALL_DIR=$HOME/.local/bin curl -sSL https://raw.githubusercontent.com/Fabbbou/duh/main/install.sh | sh
```

### Setup

Add to your shell config (`~/.bashrc`, `~/.zshrc`):
```bash
eval "$(duh inject --quiet)"
```

>If duh has already been used before, you might have an alias command for this: `duh_reload`

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
duh repository edit-gitconfig <name>          # Edit the <repo>/gitconfig file

# Path
duh path                                      # Show base configuration path  
duh path list                                 # Show base path and all repository paths

#version
duh --version
duh version [-d]

# Force duh to reload
duh_reload
```

### Example

```bash
duh alias set ll "ls -la"
duh exports set EDITOR "vim"
duh_reload
 
# Injects:
# alias ll="ls -la"
# export EDITOR="vim"
```


## Roadmap for v1.0.0
- Complete repo example at https://github.com/Fabbbou/my-duh
- Diagrams and GIF to explain the tool
- More detailled documentation about the main features
- Handling shell functions injection
- Waiting for feedbacks
- Tech
   - auto setup with installation script (the `eval $(duh inject)` thingy)
   - reworking architecture
