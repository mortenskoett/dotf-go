package cli

import (
	"fmt"
	"os"

	"github.com/mortenskoett/dotf-go/pkg/logging"
	"github.com/mortenskoett/dotf-go/pkg/parsing"
	"github.com/mortenskoett/dotf-go/pkg/terminalio"
)

type installCommand struct {
	*commandBase
	UserInteractor UserInteractor
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
		commandBase: &commandBase{
			Name:     name,
			Overview: "Install file/dir from dotfiles into userspace.",
			Usage:    name + " <filepath> [--help]",
			Args:     []arg{{Name: "file/dir", Description: "Path to file/dir inside dotfiles or path to file/dir in userspace."}},
			Flags: []*parsing.Flag{
				parsing.NewValueFlag(FlagExternal, "Install a dotfile from an external location.", "directory-path"),
			},
			Description: desc,
		},
		UserInteractor: StdInUserInteractor{},
	}
}

func (c *installCommand) Run(args *parsing.CommandlineInput, conf *parsing.DotfConfiguration) error {
	fpath := args.PositionalArgs[0]

	// Handle flags
	for _, f := range c.Flags {
		switch f.Name {
		case FlagExternal:
			if args.Flags.Exists(f) {
				externaldir, err := args.Flags.Get(f)
				if err != nil {
					return err
				}
				return c.externalInstall(fpath, externaldir, conf)
			}
		}
	}
	return c.internalInstall(fpath, conf.UserspaceDir, conf.DotfilesDir)
}

// Install file outside current dotfiles directory.
func (c *installCommand) externalInstall(file, externaldir string, conf *parsing.DotfConfiguration) error {
	var dst string

	_, err := terminalio.CopyExternalDotfile(file, externaldir, conf.DotfilesDir, true)
	if err != nil {
		switch e := err.(type) {
		case *terminalio.ErrConfirmProceed:
			logging.Warn(fmt.Sprintf("The following path will be created: %s", e.Path))
			if !c.UserInteractor.ConfirmByUser("Do you want to continue?") {
				logging.Info("Aborted by user")
				return nil
			}
			dst, err = terminalio.CopyExternalDotfile(file, externaldir, conf.DotfilesDir, false)
			if err != nil {
				return err
			}
		default:
			return err
		}
	}
	return c.internalInstall(dst, conf.UserspaceDir, conf.DotfilesDir)
}

// Install file already inside current dotfiles directory.
func (c *installCommand) internalInstall(file, userspacedir, dotfilesdir string) error {
	err := terminalio.InstallDotfile(file, userspacedir, dotfilesdir, false)
	if err != nil {
		switch e := err.(type) {
		case *terminalio.ErrAbortOnOverwrite:
			logging.Warn(fmt.Sprintf("A file already exists in userspace: %s", logging.Color(e.Path, logging.Green)))
			logging.Warn(fmt.Sprintf("It is required to backup and delete this file to install the dotfile."))

			ok := ConfirmByUser("Do you want to continue?", os.Stdin)
			if ok {
				return terminalio.InstallDotfile(file, userspacedir, dotfilesdir, ok) // Overwrite file
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
