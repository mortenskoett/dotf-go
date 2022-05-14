/*
Package resource maintains filepaths to used resource e.g. icons.
*/
package resource

import (
	"embed"
)

type IconPath string

const (
	icons                  IconPath = "assets/icons/"
	PinkLowerCase          IconPath = icons + "d_pink_lower_case.png"
	PinkLowerCaseDragon    IconPath = icons + "d_pink_lower_case_dragon.png"
	PinkLowerCaseTimeGlass IconPath = icons + "d_pink_lower_case_timeglass.png"
)

//go:embed assets/icons
var iconsEmbedded embed.FS

func GetIcon(resource IconPath) ([]byte, error) {
	return iconsEmbedded.ReadFile(string(resource))
}
