# Duh

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
         │    └── db.toml  #contains exports
         │                  and aliases     
         └──a_repo/                         
             │                              
             └── db.toml                                   
```

(Source link)[https://asciiflow.com/#/share/eJytkEFqwzAQRa8ipl20EGIItDQ%2BiyCZ2FOsokhGIwenIaHkBF140UVP55NUbhzaUBfbkNloRujrzf87MLgmiE2h9QQ0bslBDDsJG3KsrJEQzyYSynDOH%2Beh2zY3Tw%2Bh81T6MEg4RNNXzlwi6uqtrt5FX0lpfoa6%2BgwysaQNanF7lxaZUOaFEn%2B%2FvJQEygobzDflYxylv8Y%2FDwtpm6COOENHUdg86vu9ro7j1jml0wRUMLlF7uiZHJmEeOrtWl%2BNUrUUR7ll5a1TxP%2B5%2BZ3UINSl4GzolN1A0XEIrFNydpau2sRuEms8KsOCytw6z71%2B0KQCtUIm7uC0AFw02Q3zM8RMh%2BCPlQ4R7GH%2FBZNU7XY%3D)]

## Cmd

Local build:

```sh
GOOS=linux GOARCH=amd64 go build -o duh cmd/cli/main.go
```
