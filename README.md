dotf-go - Dotfiles handler written in go
-------------------------------------
The intended functionality revolves around keeping local dotfiles and remote dotfiles up-to-date
across two or more systems using an intermediary git repository.
The dotf-go project consists of the following components:
- `dotf` is a command-line tool to add files to the local dotfiles repository.
- `dotf-tray` is a trayicon-based process that facilitates automatic updates between the local repo and the remote.

Dependencies
------------
go, gotk, gcc, libappindicator-gtk3, git

Prerequisites to use dotf-go
------------
- It is required that a remote repository is setup with ssh keys so that both push/pull does not require user/pass.

Todo
----

### Basic dotf functionality already implemented in bash script
- [ ] Move a file or directory into a dotfiles dir and replace the file with a symbolic link pointing to this location

### New features
- [ ] It should be possible to change the location of a dotfile in user space
			and by updating symlinks and file location in dotfiles dir, e.g. `dotf move <current_symlink_location> <new_symlink_location>`

- [ ] It should be possible to install downloaded dotfiles, i.e. create symlinks for specific files in the dotfiles repo to that 
			same location in user space. (logic from dotf-move can be used here)

- [ ] It should be possible to move the dotfiles dir and update all symlinks 
			(this needs only to integrated into code base. It is implemented already in the dotf-move CLI.)

- [ ] A CLI UI should be implemented to give an overview of the status of both dotfiles and user space w. functionality:
	- Install specific dotfile.
	- Move specific dotfile.
	- See dotfiles that are not installed.
	- Revert a dotfile back to its original location.

- [ ] It should be possible to configure the settings of the application.
- [x] When new files are pushed to the remote from sys A, they should immediately be downloaded down into the dotfiles dir of sys B
- [x] When a file is added locally to the dotfiles dir, it should be uploaded as soon as possible to the remote
- [x] If the added files cannot be uploaded when they are added, an attempt to upload them should be made every X time.
- [x] The application must then run in the background
- [x] A systray icon should be visible to give status info and signal that the service is running.

Notes
-----
- Web app used to generate ASCII font: https://texteditor.com/ascii-art/
