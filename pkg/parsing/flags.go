// The flags package encapsulates behaviour to share flag names across pkgs
package parsing

import "fmt"

type Flag struct {
	Name        string
	Description string
	ValueName   string
}

func NewFlag(name, description string) *Flag {
	return &Flag{
		Name:        name,
		Description: description,
		ValueName:   "",
	}
}

func NewValueFlag(name, description, valuename string) *Flag {
	return &Flag{
		Name:        name,
		Description: description,
		ValueName:   valuename,
	}
}

// Contains flags with/without affixed value as parsed from commandline
type FlagHolder struct {
	flags map[string]string
}

// Wrapper around a map to contain both value-carrying and non-value-carrying flags.
// A non-value-carrying flag (key) simply has no value (empty string). If a flag does not contain
// a value an error is returned if the holder is asked for the value.
func NewFlagHolder(flags map[string]string) *FlagHolder {
	return &FlagHolder{
		flags: flags,
	}
}

func NewEmptyFlagHolder() *FlagHolder {
	return &FlagHolder{
		flags: map[string]string{},
	}
}

// Check if one of multiple flags exists
func (cl *FlagHolder) OneOf(fs []*Flag) bool {
	for _, f := range fs {
		if _, ok := cl.flags[f.Name]; ok {
			return ok
		}
	}
	return false
}

// Check if flag exists
func (cl *FlagHolder) Exists(f *Flag) bool {
	if _, ok := cl.flags[f.Name]; ok {
		return ok
	}
	return false
}

// Get keys of contained flags
func (cl *FlagHolder) GetAllKeys() []string {
	var fs []string
	for k := range cl.flags {
		fs = append(fs, k)
	}
	return fs
}

// Get number of contained flags
func (cl *FlagHolder) Count() int {
	return len(cl.flags)
}

// Get value of a value carrying flag or empty string
func (cl *FlagHolder) GetOrEmpty(f *Flag) string {
	if val, ok := cl.flags[f.Name]; ok {
		return val
	}
	return ""
}

// Get value of a value carrying flag. Fails if given flag does not carry a value.
func (cl *FlagHolder) Get(f *Flag) (string, error) {
	if val, ok := cl.flags[f.Name]; ok && val != "" {
		return val, nil
	}
	return "", fmt.Errorf("flag '--%s' does not carry a value", f.Name)
}
