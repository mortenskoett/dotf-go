```    
    _       _     __         __ _      
 __| | ___ | |_  / _|  ___  / _` | ___ 
/ _` |/ _ \|  _||  _| |___| \__. |/ _ \
\__/_|\___/ \__||_|         |___/ \___/

dotf-go - dotfiles handler written in go
```
The intended functionality of this tool is to manage dotfiles locally and to the dotfiles up-to-date
across two or more systems using an intermediary git repository and a tray icon based service.

The dotf-go project consists of the following components:
- `dotf-cli` is a command-line tool to manage the dotfiles in the repository and their relationship
	to the userspace (system) through symlinks.
- `dotf-tray` is a trayicon-based process that facilitates automatic updates between the local repo and the remote.

### Dependencies
go, gotk, gcc, libappindicator-gtk3, git

### Prerequisites to use dotf-go
- It is required that a remote repository is setup with ssh keys so that both push/pull does not require user/pass.

