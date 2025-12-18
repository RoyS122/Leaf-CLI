package models

import "fmt"

type Version struct {
	Major, Minor, Patch uint8
}

func (v Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}
