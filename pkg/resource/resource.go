/*
Package resource maintains filepaths to used resource e.g. icons.
*/
package resource

import (
	"embed"
)

type IconPath string

const (
	PinkLowerCase          IconPath = "d_pink_lower_case.png"
	PinkLowerCaseDragon    IconPath = "d_pink_lower_case_dragon.png"
	PinkLowerCaseTimeGlass IconPath = "d_pink_lower_case_timeglass.png"
)

//go:embed assets/icons
var icons embed.FS

func GetIcon(resource IconPath) ([]byte, error) {
	return icons.ReadFile(string(resource))
}
