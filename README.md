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
- Ubuntu: `apt-get install go git gcc libgtk-3-dev libayatana-appindicator3-dev`
- Arch: `pacman -Ss go git gcc libayatana-appindicator`

### Installation
1. Install dependencies
2. Clone repo and run `make install` to install both cli and tray application
3. Debug your Go env vars to be correctly setup (in case 2 errors): https://go.dev/doc/tutorial/compile-install
4. Finally add a configuration file by following description below

### Setup dotf configuration
- Create a local git repository and make sure it connects to remote using ssh.
- Setup a config in `~${HOME}/.config/dotf/config` with the exact key/value pairs seen below. *Remember to substitute values in the right column for actual ones.*
```
userspacedir        = "~/"
dotfilesdir         = "path/to/dotfiles/dir/replicating/userspace"
syncdir             = "path/to/repo/root/used/to/sync"
autosync            = false
syncinterval        = "1200"
```

