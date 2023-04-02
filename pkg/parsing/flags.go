// The flags package encapsulates behaviour to share flag names across pkgs
package parsing

import "fmt"

// Contains flags with/without affixed value as parsed from commandline
type FlagHolder struct {
	flags map[string]string
}

func NewEmptyFlagHolder() *FlagHolder {
	return &FlagHolder{
		flags: map[string]string{},
	}
}

// Wrapper around a map to contain both value-carrying and non-value-carrying flags. If a flag does
// not contain a value an error is returned if the holder is asked for the value.
func NewFlagHolder(flags map[string]string) *FlagHolder {
	return &FlagHolder{
		flags: flags,
	}
}

type Flag struct {
	Name        string
	Description string
}

func NewFlag(id, desc string) *Flag {
	return &Flag{
		Name:        id,
		Description: desc,
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

// Get value of a value carrying flag. Fails if given flag does not carry a value.
func (cl *FlagHolder) GetValue(f *Flag) (string, error) {
	if val, ok := cl.flags[f.Name]; ok {
		return val, nil
	}
	return "", fmt.Errorf("flag %s does not carry a value", f.Name)
}
