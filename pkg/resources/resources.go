/*
Resources maintains pointers to all used resources.
*/
package resources

import (
	"errors"
	"io/fs"
	"log"
	"mskk/dotf-go/pkg/projectpath"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var (
	resourcesDirName = path.Join(projectpath.Root, "assets")
)

// Private
type resourceContainer struct {
	storage map[string][]byte
}

// Loads and returns a container with the contents of the resource dir.
func newResourceContainer() *resourceContainer {
	container := &resourceContainer{storage: make(map[string][]byte)}
	container.loadResources(resourcesDirName)
	return container
}

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
		log.Fatal("Fatal error encountered while parsing resources:", err)
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

// Public resource manager.
var resourceManagerSingleton = newResourceContainer()

// Get resource found at 'path' relative to resource dir.
func Get(path string) ([]byte, error) {
	return resourceManagerSingleton.get(path)
}

// Add resource to be available for the duration of the runtime.
func Add(path string, file []byte) {
	resourceManagerSingleton.add(path, file)
}

// Contains checks whether a resource exists at 'path' relative to resource dir.
func Contains(path string) bool {
	return resourceManagerSingleton.contains(path)
}
