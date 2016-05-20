package filters

import (
	"path/filepath"
	"strings"
)

func RepoIsSymlinkedToRepo(workingDirectory, repoPath string) bool {

	// attempt to evaluate the symlink
	symlinkedPath, err := filepath.EvalSymlinks(repoPath)

	// if it's not a symlink return false
	if err != nil || symlinkedPath == repoPath {
		return false
	}

	// shorten paths by removing . and ..
	symlinkedPath = filepath.Clean(symlinkedPath)
	workingDirectory = filepath.Clean(workingDirectory)

	// get the working directory of the symlinked location
	symlinkWorkingDirectory, _ := filepath.Split(symlinkedPath)

	// remove trailing path seperators
	pathSeparator := string(filepath.Separator)
	symlinkWorkingDirectory = strings.TrimRight(symlinkWorkingDirectory, pathSeparator)
	workingDirectory = strings.TrimRight(workingDirectory, pathSeparator)

	return symlinkWorkingDirectory == workingDirectory
}
