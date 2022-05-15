```    
    _       _     __         __ _      
 __| | ___ | |_  / _|  ___  / _` | ___ 
/ _` |/ _ \|  _||  _| |___| \__. |/ _ \
\__/_|\___/ \__||_|         |___/ \___/

dotf-go - dotfiles handler written in go
```
The primary intended functionality of this tool is to manage dotfiles locally, but the tool also
makes it possible to easier keep dotfiles up-to-date across two or more systems using an
intermediary git repository and a tray application.

The dotf-go project consists of the following components:
- `dotf-cli` is a command-line tool to manage the dotfiles in the repository and their relationship
	to the userspace (system) through symlinks.
- `dotf-tray` is a trayicon-based process that facilitates automatic updates between the local repo and the remote.

### Usage
WIP

### Dependencies
go, gotk, gcc, libappindicator-gtk3, git

### Setup
- Create a local git repository and make sure it connects to remote using ssh.
- Setup a config in ~${HOME}/.config/dotf/config with the following contents:
```
remoteurl		= https://www.yourrepo.com/yourname/doesntexist
dotfilesdir		= "~/yourdotfilesdir"
homedir			= "~/"
updateintervalsec = 120
```

