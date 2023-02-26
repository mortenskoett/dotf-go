// Environment to replicate dotfiles environment in memory used for testing.
package test

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
)

// File hierachy test environment. **Remember to call cleanup after use**.
type Environment struct {
	// A contained directory in the environment
	DotfilesDir, UserspaceDir, BackupDir *DirectoryHandle
}

// Path file handle used by test Environment
type DirectoryHandle struct {
	Name string // The base of the path
	Path string // Path to directory
}

// Returns a file hierachy based testing environment
func NewTestEnvironment() *Environment {
	return &Environment{
		DotfilesDir:  newTestFilePathHandle("dotfiles-dir*"),
		UserspaceDir: newTestFilePathHandle("userspace-dir*"),
		BackupDir:    newTestFilePathHandle("backup-dir*"),
	}
}

// Returns a filepath handle used in the test Environment
func newTestFilePathHandle(name string) *DirectoryHandle {
	dir, err := os.MkdirTemp("", name)
	if err != nil {
		log.Fatal("could not create testFileHandle:", err)
	}
	return &DirectoryHandle{Name: name, Path: dir}
}

// Adds a directory to the filepath and returns its handle
func (e *DirectoryHandle) AddTempDir(relativePath string) *DirectoryHandle {
	path := filepath.Join(e.Path, relativePath)
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		log.Fatal("could not create nested directories:", err)
	}
	return &DirectoryHandle{Name: filepath.Base(path), Path: path}
}

// Adds a file to the filepath and returns its handle
func (e *DirectoryHandle) AddTempFile() *os.File {
	f, err := os.CreateTemp(e.Path, "file-*") // Suffixes a random string
	if err != nil {
		log.Fatal(err)
	}
	return f
}

// Creates a randomly named symlink pointing at 'tofile'. The path of the symlink is returned.
func (e *DirectoryHandle) CreateSymlink(tofile string) string {
	rnd := rand.Int31()
	symname := fmt.Sprintf("symlink-%d", rnd)
	symlinkpath := filepath.Join(e.Path, symname)

	err := os.Symlink(tofile, symlinkpath)
	if err != nil {
		log.Fatal(err)
	}

	return symlinkpath
}

// Cleans up the test environment. Should be called when done e.g. using defer
func (e *Environment) Cleanup() {
	if err := os.RemoveAll(e.DotfilesDir.Path); err != nil {
		log.Fatal(err)
	}
	if err := os.RemoveAll(e.UserspaceDir.Path); err != nil {
		log.Fatal(err)
	}
	if err := os.RemoveAll(e.BackupDir.Path); err != nil {
		log.Fatal(err)
	}
}
