package git

import (
	"io/ioutil"
	"log"
	"os/exec"
	"path/filepath"
	"strings"
)

// gitRepo provides interaction with git repos on the local machine
type Repo struct {
	WorkingDirectory string
	RepoName         string
	OriginUrl        string
}

// constructor for gitRepo
func NewRepo(workingDirectory, repoName, originUrl string) Repo {
	gr := Repo{
		workingDirectory,
		repoName,
		originUrl,
	}

	if gr.RepoName == "" && gr.OriginUrl != "" {
		_, repoName := filepath.Split(gr.OriginUrl)
		gr.RepoName = strings.TrimSuffix(repoName, filepath.Ext(repoName))
	}

	if gr.RepoName == "" {
		log.Fatal("Could not determine git repo name from origin url")
	}

	return gr
}

// determine if this gitRepo is actually a git repo
func (gr *Repo) IsGitRepo() bool {

	// is this a directory at all?
	_, err := ioutil.ReadDir(gr.GetFullRepoPath())
	if err != nil {
		return false
	}

	// is there a .git folder in there?
	_, err = ioutil.ReadDir(gr.GetFullRepoPath() + "/.git")
	if err != nil {
		return false
	}

	// is this actually a git repo?
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	cmd.Dir = gr.GetFullRepoPath()
	out, _ := cmd.Output()
	if strings.Trim(string(out), " \n\r") != "true" {
		return false
	}

	return true
}

// determines if there are any uncommited files
func (gr *Repo) IsClean() bool {
	diffOutput := gr.RunCmd(gr.GetFullRepoPath(), "diff", "--name-only")
	return len(diffOutput) == 0
}

// returns full path to the git repo
func (gr *Repo) GetFullRepoPath() string {
	return filepath.Join(gr.WorkingDirectory, gr.RepoName)
}

// runs git commands
func (gr *Repo) RunCmd(workingDirectory string, command ...string) string {

	// create the command
	cmd := exec.Command("git", command...)

	// run the command in our git repo directory
	cmd.Dir = workingDirectory

	// redirect output
	stdoutPipe, _ := cmd.StdoutPipe()
	stderrPipe, _ := cmd.StderrPipe()

	// star the command
	cmd.Start()

	// wait until the command finishes
	defer cmd.Wait()

	// get output
	stdout, _ := ioutil.ReadAll(stdoutPipe)
	stderr, _ := ioutil.ReadAll(stderrPipe)

	return string(append(stdout, stderr...))
}

// runs a git fetch
func (gr *Repo) Fetch(fetchArgs []string) string {
	fetchArgs = append([]string{"fetch"}, fetchArgs...)
	return gr.RunCmd(gr.GetFullRepoPath(), fetchArgs...)
}

// runs a git merge
func (gr *Repo) Merge(mergeArgs []string) string {
	mergeArgs = append([]string{"merge"}, mergeArgs...)
	return gr.RunCmd(gr.GetFullRepoPath(), mergeArgs...)
}

// runs a git clone
func (gr *Repo) Clone(cloneArgs []string) string {
	cloneArgs = append([]string{"clone", gr.OriginUrl, "--progress"}, cloneArgs...)
	return gr.RunCmd(gr.WorkingDirectory, cloneArgs...)
}

func (gr*Repo) Checkout(checkoutArgs []string) string {
	checkoutArgs = append([]string{"checkout"}, checkoutArgs...)
	return gr.RunCmd(gr.GetFullRepoPath(), checkoutArgs...)
}
