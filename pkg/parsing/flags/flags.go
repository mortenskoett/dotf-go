// The flags package encapsulates behaviour to share flag names across pkgs
package flags

import "fmt"

const (
	Help   boolFlag  = "help"
	Select boolFlag  = "select"
	Config valueFlag = "config"
)

// The overall flag type
type Flag interface {
	str() string
}

// More specific flag type
type ValueFlag interface {
	Flag
	get() // empty
}

// Bool flags are true/false flags
type boolFlag string

func (f boolFlag) str() string {
	return string(f)
}

// Value flags contains a value
type valueFlag string

func (f valueFlag) str() string {
	return string(f)
}

func (f valueFlag) get() {
	// noop to restrict acces to GetValue method
}

// Contains flags with/without affixed value as parsed from commandline
type FlagHolder struct {
	valueFlags map[string]string // Example --name john or --name=john
	boolFlags  map[string]bool   // Example --verbose
}

func NewEmptyFlagHolder() *FlagHolder {
	return &FlagHolder{
		valueFlags: map[string]string{},
		boolFlags:  map[string]bool{},
	}
}

func NewFlagHolder(bfs map[string]bool, vfs map[string]string) *FlagHolder {
	return &FlagHolder{
		valueFlags: vfs,
		boolFlags:  bfs,
	}
}

// Check if flag exists
func (cl *FlagHolder) Exists(f Flag) bool {
	if _, ok := cl.boolFlags[f.str()]; ok {
		return ok
	}

	if _, ok := cl.valueFlags[f.str()]; ok {
		return ok
	}

	return false
}

// Get value of a value carrying flag
func (cl *FlagHolder) GetValue(f ValueFlag) (string, error) {
	if val, ok := cl.valueFlags[f.str()]; ok {
		return val, nil
	}
	return "", fmt.Errorf("flag %s does not carry a value", f.str())
}
