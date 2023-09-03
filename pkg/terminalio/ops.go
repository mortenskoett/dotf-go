package terminalio

import (
	"fmt"
	"os"
	"path/filepath"
)

// Copies file from userspace to the dotfiles directory and creates symlink from userspace file into
// the newly copied file in the dotfiles dirctory. The identical relative path is used for both
// 'homeDir' and 'dotfilesDir'.
func AddDotfile(userspaceFile, userspaceHomedir, dotfilesDir string) error {
	absUserspaceFile, err := GetAndValidateAbsolutePath(userspaceFile)
	if err != nil {

		return err
	}

	absHomedir, err := GetAndValidateAbsolutePath(userspaceHomedir)
	if err != nil {
		return err
	}

	absDotfilesDir, err := GetAndValidateAbsolutePath(dotfilesDir)
	if err != nil {
		return err
	}

	// Construct path inside dotfiles dir
	absNewDotFile, err := replacePrefixPath(absUserspaceFile, absHomedir, absDotfilesDir)
	if err != nil {
		return err
	}

	// Assert a file is not already in dotfiles dir at location
	exists, err := CheckIfFileExists(absNewDotFile)
	if err != nil {
		return err
	}
	if exists {
		return &ErrFileAlreadyExists{absNewDotFile}
	}

	// Backup file before copying it
	_, err = backupFile(absUserspaceFile)
	if err != nil {
		return err
	}

	// Copy file to dotfiles
	_, err = copyFileOrDir(absUserspaceFile, absNewDotFile)
	if err != nil {
		return err
	}

	// Remove file in userspace
	if err := deleteFileOrDir(absUserspaceFile); err != nil {
		return err
	}

	// Create symlink from userspace to the newly created file in dotfiles
	if err := createSymlink(absUserspaceFile, absNewDotFile); err != nil {
		return err
	}

	return nil
}

// Copies a file from an external location into current dotfiles directory. E.g. fromdir can be the
// root of another dotfiles directory and todir can be the path to current dotfiles root dir.
// Returns the path of the copied file.
// If confirm==true an empty string and an error containing the calucated path to the new dotfile
// are returned.
// If the file to be copied is a symlink, a symlink will be created in todir pointing back to the
// symlink in fromdir. This is to keep the semantics, that userspace links always point into their
// own dotfiles dir.
func CopyExternalDotfile(fpath, fromdir, todir string, confirm bool) (string, error) {
	absfilepath, err := GetAndValidateAbsolutePath(fpath)
	if err != nil {
		return "", err
	}

	absExtDotfilesDir, err := GetAndValidateAbsolutePath(fromdir)
	if err != nil {
		return "", err
	}

	absDotfilesDir, err := GetAndValidateAbsolutePath(todir)
	if err != nil {
		return "", err
	}

	// Construct path inside dotfiles dir.
	absNewDotfile, err := replacePrefixPath(absfilepath, absExtDotfilesDir, absDotfilesDir)
	if err != nil {
		return "", err
	}

	// If wanted by caller the process can abort here to show the calculated new path.
	if confirm {
		return "", &ErrConfirmProceed{Path: absNewDotfile}
	}

	// Assert a file is not already in dotfiles dir at location.
	exists, err := CheckIfFileExists(absNewDotfile)
	if err != nil {
		return "", err
	}
	if exists {
		return "", &ErrFileAlreadyExists{absNewDotfile}
	}

	// Determine whether given file is a symlink.
	ok, err := IsFileSymlink(absfilepath)
	if err != nil {
		return "", fmt.Errorf("failed to determine if given file was a symlink: %v", err)
	}

	// Whether given file from external flag is a symlink.
	if ok {
		// Get file pointed to by the symlink.
		relativeToSymlinkPath, err := os.Readlink(absfilepath)
		if err != nil {
			return "", fmt.Errorf("failed to get src of symlink: %v", err)
		}

		// Create absolute path to file pointed to by the symlink.
		absSrcPath := filepath.Join(absExtDotfilesDir, relativeToSymlinkPath)

		validatedAbsSrcPath, err := GetAndValidateAbsolutePath(absSrcPath)
		if err != nil {
			return "", fmt.Errorf("failed to validate file pointed to by symlink: %v", err)
		}

		// We can now create a symlink pointing to the file pointed to by the symlink.
		if err := createSymlink(absNewDotfile, validatedAbsSrcPath); err != nil {
			return "", err
		}
		return absNewDotfile, nil
	}

	// Backup file before copying it
	_, err = backupFile(absfilepath)
	if err != nil {
		return "", err
	}

	// Copy file to dotfiles
	return copyFileOrDir(absfilepath, absNewDotfile)
}

// Writes a file to disk
func WriteFile(fpath string, contents []byte, overwrite bool) error {
	exists, err := CheckIfFileExists(fpath)
	if err != nil {
		return err
	}

	if exists {
		if !overwrite {
			return &ErrAbortOnOverwrite{fpath}
		}

		// Backup file before deleting it
		if _, err = backupFile(fpath); err != nil {
			return err
		}

		// Delete file
		if err := deleteFile(fpath); err != nil {
			return err
		}
	}

	// Create new file
	if err := writeFile(fpath, contents); err != nil {
		return err
	}

	return nil
}

// Installs a dotfile into its relative equal location in userspace by way of a symlink in userspace
// pointing back to the file in dotfiles. The userspace file will be removed if 'overwrite' is true.
// Both the filepath inside dotfile as well as in userspace can be given.
func InstallDotfile(file, userspaceDir, dotfilesDir string, overwrite bool) error {
	info, err := getFileLocationInfo(file, userspaceDir, dotfilesDir)
	if err != nil {
		return err
	}

	// Check whether dotfile exists
	exists, err := CheckIfFileExists(info.dotfilesFile)
	if err != nil {
		return err
	}
	if !exists {
		return &ErrFileNotFound{info.dotfilesFile}
	}

	// Check whtether userspace file already exists
	exists, err = CheckIfFileExists(info.userspaceFile)
	if err != nil {
		return err
	}
	if exists {
		if !overwrite {
			return &ErrAbortOnOverwrite{info.userspaceFile}
		}

		// Backup file before copying it
		if _, err = backupFile(info.userspaceFile); err != nil {
			return err
		}

		// Remove file in userspace
		if err := deleteFile(info.userspaceFile); err != nil {
			return err
		}
	}
	// Create symlink in userspace pointing to dotfile
	if err := createSymlink(info.userspaceFile, info.dotfilesFile); err != nil {
		return err
	}

	return nil
}

// Reverts the insertion of a file into the dotfiles directory and return it to its original
// location in userspace. The symlink is removed first. The operation can be applied both to the
// symlink in userspace and the actual file in the dotfiles directory.
func RevertDotfile(file, userspaceDir, dotfilesDir string) error {
	info, err := getFileLocationInfo(file, userspaceDir, dotfilesDir)
	if err != nil {
		return err
	}

	dotfile := info.dotfilesFile
	usersymlink := info.userspaceFile

	// Check whtether file and symlink exists
	ok, err := CheckIfFileExists(dotfile)
	if err != nil {
		return err
	}
	if !ok {
		return &ErrFileNotFound{dotfile}
	}

	ok, err = IsFileSymlink(usersymlink)
	if err != nil {
		return err
	}
	if !ok {
		return &ErrSymlinkNotFound{usersymlink}
	}

	// Backup file before copying it
	if _, err = backupFile(dotfile); err != nil {
		return err
	}

	// Remove symlink in userspace
	if err := deleteFile(usersymlink); err != nil {
		return err
	}

	// Copy dotfile back to userspace
	if _, err = copyFileOrDir(dotfile, usersymlink); err != nil {
		return err
	}

	// Remove file in dotfiles
	if err := deleteFileOrDir(dotfile); err != nil {
		return err
	}

	return nil
}
