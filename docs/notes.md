# Development notes
===================

# Sat Sep 24 12:46:57 PM CEST 2022
Rethinking shared dotfiles into multiple repos to fix git single repo issues as currently
encountered.

## Problem
Problem is the dir structure of my dotfiles setup. The setup folder is on same level as dotfiles
under each distro. This will confuse git because files are not added in the full repo.

The complexity lies in the shared setup because then you'd need a pointer both to the root dotfiles
dir but also the specific distro in that setup.

Config is currently only read when any command is given i.e. if there is a config error is not read
on ./dotf

## New idea
Basically to put repos like `shared`, `configs` etc in the config and then be able to pick and
choose between these when setting up dotf locally.

### Pros
	- Then I would not have to have dotfiles from all systems placed on each system
	- Then I would only have a single repo source of truth for each dotf instance
	- Then I would be able to more meaningfully handle shared dotfiles between dotf instances

### Cons
	- How to deal with filenames and system env vars
	- If repos are in the config, then we need to make assumptions about filepaths

## New config
```
	dotfiles_main		= "www.github.com/mortenskoett/dotfiles"
	dotfiles_sub		= "www.github.com/mortenskoett/shared_dfiles"
	dotfiles_sub		= /home/mskk/dotfiles/configuration/
	dotfiles_sub		= /home/mskk/dotfiles/dockerfiles/
	userspace_root		= "~/"
	updateintervalsec 	= 1200
```

## Minimum viable changes
	- configuration must be read differently
	- dotf must have another layer of flags to handle specific repos
	- it should be possible to designate whether a configured repo is dotfiles or something else

# Fri May 13 11:34:14 AM CEST 2022
If revert is called. How do we know whether it is called on the symlink in userspace or the actual
file in dotfiles?

/home/msk/
/home/msk/dotfiles/

## detach absfiles and abshomedir to get name of dfiles or use filepath.Base
= dotfiles

/home/msk/dotfiles/bla/dotfiles/fil
/home/msk/bla/dotfiles/fil

## detach filepath w. homedir to get relative path
= bla/dotfiles/fil

1. we get the name of the dotfiles folder
2. we detach the file using homedir:
	- if prefix == dotfiles folder name then we know it is inside dotfiles
	- if not then it is in userspace
3. now we can assert correctly on symlink in userspace and file/folder in dotfiles


# Tue May  3 09:30:50 PM CEST 2022
Further thougths of config file reading and default settings.
A need to refactor lib functions is arising:
	- separate some of the pkg libs into sub modules e.g. terminalio
	- make commands use the lib functions themselves in order to be able to reuse the lib code

# Sun Apr 17 09:06:31 AM CEST 2022
Reading config:
	1. Read config at default location: ~/.config/dotf/config
	2. Read config at location w. flag: --config (try to use flags pkg)
	3. Default behaviour if no config: fail and print basic config to terminal.

Arg-parser that handles:
	1. Positional args (command arg)
	2. Both bool and Value Flags specific to each command
	3. General flags should be parsed first or propogated to command

	Order of parsing:
	Command -> General flags -> Specific flags

	Types
	--flag
	command

# Tue Apr  5 09:08:00 PM CEST 2022
- Moved todo.md into depcrecated and moved all bullets into Todoist.

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

