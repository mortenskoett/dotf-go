package terminalio

// Installs a dotfile into its relative equal location in userspace by way of a symlink in userspace
// pointing back to the file in dotfiles. The userspace file will be removed if 'overwriteFiles' is
// true.
func InstallDotfile(file, homeDir, dotfilesDir string, overwriteAllowed bool) error {
	info, err := getFileLocationInfo(file, homeDir, dotfilesDir)
	if err != nil {
		return err
	}

	// Check whether dotfile exists
	exists, err := checkIfFileExists(info.dotfilesFile)
	if err != nil {
		return err
	}
	if !exists {
		return &FileNotFoundError{info.dotfilesFile}
	}

	// Check whtether userspace file already exists
	exists, err = checkIfFileExists(info.userspaceFile)
	if err != nil {
		return err
	}
	if exists {
		if !overwriteAllowed {
			return &AbortOnOverwriteError{info.userspaceFile}
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
	ok, err := checkIfFileExists(dotfile)
	if err != nil {
		return err
	}
	if !ok {
		return &FileNotFoundError{dotfile}
	}

	ok, err = isFileSymlink(usersymlink)
	if err != nil {
		return err
	}
	if !ok {
		return &SymlinkNotFoundError{usersymlink}
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

// Copies file in userspace to dotfiles dir using same relative path between 'homeDir' and
// 'dotfilesDir'. The file is backed up first.
func AddFileToDotfiles(userspaceFile, homeDir, dotfilesDir string) error {
	absUserspaceFile, err := GetAndValidateAbsolutePath(userspaceFile)
	if err != nil {

		return err
	}

	absHomedir, err := GetAndValidateAbsolutePath(homeDir)
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
	exists, err := checkIfFileExists(absNewDotFile)
	if err != nil {
		return err
	}
	if exists {
		return &FileAlreadyExistsError{absNewDotFile}
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

