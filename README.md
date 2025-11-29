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
