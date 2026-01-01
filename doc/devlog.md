
## Config file structure

```txt
~/.zshrc ─┐                                 
          ├─ `eval $(duh inject)`           
~/.bashrc─┘                                 
                                            
                                            
~/.local/share/duh/                         
    │                                       
    ├── user_preferences.toml               
    │                                       
    └── repositories/                       
         │                                  
         ├── local/                         
         │    │                             
         │    ├── db.toml  #contains exports
         │    │             and aliases     
         │    │                             
         │    └── functions                 
         │         │                        
         │         ├── script_1.sh          
         │         │                        
         │         └── script_2.sh          
         │                                  
         │                                  
         │                                  
         └──a_repo/                         
             │                              
             └── db.toml                          
```

(Source link)[https://asciiflow.com/#/share/eJy9k9FqwjAUhl8lZLvYQCwTNqbPUtBjekozYlJyUtGJY%2FgEu8jFLvZ0fZKl07J1KkZwOzc5Cfnz%2FflJVlzDDPmIZ1XBMnDAcqmQkbOVcJVF3uMKlmjDjlXK52hJGp3y0aCX8kUYhw%2FD0C2blcf70DlcuDBJ%2BUvSf6bCClb719q%2FsVOVpvp7UvuPIGMTnINi1zeNN6mfULjbSVcSKFNoMF%2BU9%2FMop%2Bv87cGQMgJUQgVYTILz5NTptd%2BcZ2ebThNQRWjHpcUcLWqB1Hdmpi5G8TuKxdKQdMZKpGO3%2BZlUFKoraC%2B0zS5StImBHZS0uGy6S%2BxKGO1AamK4KI11FEECnTFQEgjpouba2PNKCxd%2BG8XI9vpoSZsFCStLN77rU%2FFHJN8lDaJIR%2Bu%2FBTv3MG6%2BQ9wTjSLtC%2Fzv13lAxNd8%2FQkytWHQ)]

## Cmd

Local binary build for linux :
```sh
GOOS=linux GOARCH=amd64 go build -o duh cmd/cli/main.go
```

Roadmap:
- code cobra CLI :
    - repository/repo/repos
        - using git lib in go to SYNC: pull, commit and push (or even create pr ?) 
    - edit files from default editor (using editor available)

- GIT aliases handling
    - injecting git aliases using a simple file to do so, using a proper gitconfig parser that handle multi-keys for a group (so includes works with multiple files)
- functions injection and edit:
    - specify how functions, files, folders are defined in a repo
    - specify the injection (parse sh/bash/zsh scripts? or let it free for users?)
    - search how to open a file or folder with an Api/Lib to use vscode or something to open a file (or vim otherwise)

- installation from brew, choco
- VSCode extension direnv like if it makes sense to have injection in it (exports mainly)