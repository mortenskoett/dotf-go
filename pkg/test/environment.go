// Environment to replicate dotfiles environment in memory used for testing.
package test

import (
	"log"
	"os"
)

// Path file handle used by test Environment
type DirectoryHandle struct {
	Name string // Will be contained in the path
	Path string // To the directory
}

// File hierachy test environment. Remember to call cleanup after use
type Environment struct {
	DotfilesDir, UserspaceDir, BackupDir *DirectoryHandle
}

// Returns a filepath handle used in the test Environment
func NewTestFilePathHandle(name string) *DirectoryHandle {
	dir, err := os.MkdirTemp("", name)
	if err != nil {
		log.Fatal("could not create testFileHandle:", err)
	}
	return &DirectoryHandle{Name: name, Path: dir}
}

// Returns a file hierachy based testing environment
func NewTestEnvironment() Environment {
	return Environment{
		DotfilesDir:  NewTestFilePathHandle("dotfiles"),
		UserspaceDir: NewTestFilePathHandle("userspace"),
		BackupDir:    NewTestFilePathHandle("backup"),
	}
}

// Adds a file to the filepath and returns its handle
func (e *DirectoryHandle) AddTempFile() *os.File {
	f, err := os.CreateTemp(e.Path, e.Name+"-*") // Suffixes a random string
	if err != nil {
		log.Fatal(err)
	}
	return f
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
