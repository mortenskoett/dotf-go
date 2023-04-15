package cli

import (
	"fmt"

	"github.com/mortenskoett/dotf-go/pkg/logging"
	"github.com/mortenskoett/dotf-go/pkg/parsing"
	"github.com/mortenskoett/dotf-go/pkg/terminalio"
)

type installCommand struct {
	*commandBase
}

func NewInstallCommand() *installCommand {
	name := "install"
	desc := `
	Will install a file or directory from dotfiles into the same location in userspace. If an
	identically named file is found in userspace a prompt will ask whether to delete the file or
	abort the operation. The only sane thing is to remove the file in userspace as the idea is for
	the previously created dotfile to take its place.

	The command can be used both on files inside the dotfiles directory as well as files in
	userspace and will do the same thing. A file inside dotfiles can be considered the source and
	a file in userspace will be considered the target. The source must exist and the target will be
	overwritten (after a backup is made).

	- If a file from inside dotfiles is given this file will be used as installation source.
	- If a file from userspace is given this file will be used as target.
	- If the '--external <directoy-path>' flag is given, it is possible to install a dotfile from an
	external dotfiles directory by giving the path of that directory. The file is copied into the
	dotfiles directory of the current distribution using the relative path from the given directory
	path and installed into userspace .`

	return &installCommand{
		&commandBase{
			Name:     name,
			Overview: "Install file/dir from dotfiles into userspace.",
			Usage:    name + " <filepath> [--help]",
			Args:     []arg{{Name: "file/dir", Description: "Path to file/dir inside dotfiles or path to file/dir in userspace."}},
			Flags: []*parsing.Flag{
				parsing.NewValueFlag(flagExternal, "Install a dotfile from an external location.", "directory-path"),
			},
			Description: desc,
		},
	}
}

func (c *installCommand) Run(args *parsing.CommandlineInput, conf *parsing.DotfConfiguration) error {
	filepath := args.PositionalArgs[0]

	// Handle flags
	for _, f := range c.Flags {
		switch f.Name {
		case flagExternal:
			if args.Flags.Exists(f) {
				// TODO: Implement the external flag
				// copy the given path into the right place in dotfiles
				// then proceed installation as usual
				extDotfilesDir, err := args.Flags.Get(f)
				if err != nil {
					return err
				}
				logging.Info(flagExternal, "flag given with value:", extDotfilesDir)
				dsuffix, err := terminalio.FindCommonPathPrefix(filepath, extDotfilesDir)
				logging.Info("suffix:", dsuffix)
				if err != nil {
					return err
				}
			}
		}
	}

	logging.Info("Exiting early")
	return nil
	logging.Info("Shouldn't happen")

	err := terminalio.InstallDotfile(filepath, conf.UserspaceDir, conf.DotfilesDir, false)
	if err != nil {
		switch e := err.(type) {
		case *terminalio.AbortOnOverwriteError:
			logging.Warn(fmt.Sprintf("A file already exists in userspace: %s", logging.Color(e.Path, logging.Green)))
			logging.Warn(fmt.Sprintf("It is required to backup and delete this file to install the dotfile."))

			ok := confirmByUser("Do you want to continue?")
			if ok {
				return terminalio.InstallDotfile(filepath, conf.UserspaceDir, conf.DotfilesDir, ok) // Overwrite file
			} else {
				logging.Info("Aborted by user")
				return nil
			}
		default:
			return err
		}
	}

	return nil
}
