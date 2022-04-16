# Prioritized todo
==================

# ------------------ MVP ------------------ 
## Come up with design to handle shared dotfiles properly
This should be implemented across the application and will affect the CLI design.

## Move functionality from the dotf.sh script into the go code base.
- Implement all features of dotf.sh legacy:
	- Day-to-day functionality
		- add		(from userspace to dotfiles)
		- remove	(from dotfiles and reinsert file into userspace)
		- push		(to git)
	- Overview functionality
		- info (show 
		- stat (git status)
		- list (list all dotfiles)

## Add ability to install dotfiles in user space
I.e. create symlinks for specific files in the dotfiles repo to that same location in user
space. (logic from dotf-move can be used here)

## Make sure icons are baked into the binary
I.e. by using the Go code gen features to actually write out files containing the icons in the
compiled binary.


# ------------------ FEATURES ------------------ 
## Add functionality to change the location of a dotfile in user space 
I.e. have symlinks and actual file location in the dotfiles dir updated, e.g. `dotf move
<current_symlink_location> <new_symlink_location>`

## Implement a CLI GUI 
I.e. to give an overview of the status of both dotfiles and user space w. functionality:
- Get overview of dotfiles / user space
- Install specific dotfile.
- Move specific dotfile.
- See dotfiles that are not installed.
- Revert a dotfile back to its original location.

## Add ability to configure the settings of the application.
A shortcut to view and edit the config file should be accessible from the dotf-tray
