# Development notes
===================

# Sat Apr  2 01:34:41 PM CEST 2022
- Web app used to generate ASCII font: https://texteditor.com/ascii-art/

## Done from old todo
- [x] Implement reasonable CLI interface
	- OK Create an abstraction so that Commands can handle themselves (printing and accessing packages etc)
	- OK Conform CLI to some of the industry standard ones git, dotnet, jq
		- OK Add ability to show help specifically for each commands by suffixing --help | -h | help
		- OK Add ability to give flags specific to each command

- [x] Add ability for different distros to share some dotfile and each update it.
	- OK Add ability to symlinks inside separate distro dotfiles pointing to
		  a shared dotfiles directory. Apparently it is possible to add symlinks to Github which makes this feature possible.
- [x] Add ability to move the dotfiles dir and update all symlinks (implemented in dotf-move)
	- OK Move logic into pkg
	- OK Setup argument parsing in dotf-cli so eg. `dotf bla1 bla2 etc` is available.
	- OK dotf-go should make this command accesible
	- OK Fixup loose ends: Is it error resilient? Does it actually handle parameters?

- [x] When new files are pushed to the remote from sys A, they should immediately be downloaded down into the dotfiles dir of sys B
- [x] When a file is added locally to the dotfiles dir, it should be uploaded as soon as possible to the remote
- [x] If the added files cannot be uploaded when they are added, an attempt to upload them should be made every X time.
- [x] The application must then run in the background
- [x] A systray icon should be visible to give status info and signal that the service is running.


