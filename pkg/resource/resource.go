/*
Package resource maintains filepaths to used resource e.g. icons.
*/
package resource

import (
	"errors"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var (
	resourcesDirName         = path.Join(ProjectRoot, "assets")
	resourceManagerSingleton = newResourceContainer()
)

type resourceContainer struct {
	storage map[string][]byte
}

// Publicly accessible through defined functions.

/* Get resource found at 'path' relative to resource dir. */
func Get(path string) ([]byte, error) {
	return resourceManagerSingleton.get(path)
}

/* Add resource to be available for the duration of the runtime. */
func Add(path string, file []byte) {
	resourceManagerSingleton.add(path, file)
}

/* Contains checks whether a resource exists at 'path' relative to resource dir. */
func Contains(path string) bool {
	return resourceManagerSingleton.contains(path)
}

// Loads and returns a container with the contents of the resource dir.
func newResourceContainer() *resourceContainer {
	container := &resourceContainer{storage: make(map[string][]byte)}
	container.loadResources(resourcesDirName)
	return container
}

/* Traverses resource dir and maps path to byte[] with the resource file. */
func (r *resourceContainer) loadResources(resourcePath string) {
	err := filepath.WalkDir(resourcePath, func(path string, fileInfo fs.DirEntry, err error) error {

		// Remove initial slash for brevity.
		relativePath := strings.TrimPrefix(strings.TrimPrefix(path, resourcePath), "/")

		if fileInfo.IsDir() {
			// Ignore.
			return nil

		} else {
			file, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			r.storage[relativePath] = file
		}

		return err
	})

	if err != nil {
		log.Fatal("Fatal error encountered while parsing resource:", err)
	}
}

func (r *resourceContainer) add(path string, file []byte) {
	r.storage[path] = file
}

func (r *resourceContainer) get(path string) ([]byte, error) {
	if f, ok := r.storage[path]; ok {
		return f, nil
	}
	return nil, errors.New("resource could not be found in resource container: " + path)
}

func (r *resourceContainer) contains(path string) bool {
	if _, ok := r.storage[path]; ok {
		return true
	}
	return false
}
