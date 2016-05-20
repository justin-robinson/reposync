package pool

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"github.com/justin-robinson/reposync/lib/git"
)

// gitResource wraps gitRepo operations
type gitResource struct{}

// runs a pull on a gitRepo with formatted output
func (resource *gitResource) Pull(repo git.Repo, prune bool) string {

	// to be printed to console
	out := bytes.Buffer{}

	// git fetch
	fetchArgs := []string{"--all"}
	if prune {
		fetchArgs = append(fetchArgs, "--prune")
	}
	out.WriteString(repo.Fetch(fetchArgs))

	// if no files changed, do a merge
	if repo.IsClean() {
		out.WriteString(repo.Merge([]string{}))
	}

	// send output back home to be printed to console
	return out.String()
}

// clones a gitRepo with formatted output
func (s *gitResource) Clone(repo git.Repo) string {

	// to be printed to console
	out := bytes.Buffer{}

	out.WriteString(repo.Clone([]string{}))
	out.WriteString(git.InitFlow(&repo))
	out.WriteString(repo.Checkout([]string{"develop"}))
	out.WriteString(s.symlinkDeps(&repo))

	// send output back home to be printed to console
	return out.String()
}

// Runs a symlink script inside the repo if the script exists
func (s *gitResource) symlinkDeps(repo *git.Repo) string {

	linkScriptPath := filepath.Join(repo.GetFullRepoPath(), "scripts", "link_dependencies.sh")

	var output string

	// does the link script exist?
	if _, err := os.Stat(linkScriptPath); err == nil {

		// split file path into dir and filename
		dir, _ := filepath.Split(linkScriptPath)

		// create command to run link script in dir
		cmd := exec.Command(linkScriptPath)
		cmd.Dir = dir

		// run command
		stdout, err := cmd.Output()

		// process output
		if err != nil {
			output = err.Error()
		} else {
			output = "Linking dependencies " + string(stdout) + "\n"
		}

	} else {
		output = "Dependency linking script not found at " + linkScriptPath + "\n"
	}

	return output
}
