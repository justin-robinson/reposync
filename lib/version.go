package lib

import (
	"fmt"
)

const Major = 0
const Minor = 2
const Patch = 0

func GetAppVersion () string {
	return fmt.Sprintf("%v.%v.%v", Major, Minor, Patch)
}
