package filters

import (
	"os"
)

func RepoExistsLocally(repoPath string) bool {
	_, err := os.Stat(repoPath)

	return !os.IsNotExist(err)
}
