package reposync

import (
	"github.com/justin-robinson/reposync/lib/pool"
	"github.com/justin-robinson/reposync/lib/git"
	"path/filepath"
)

// finds git repos directly in Directory then fetches and optionally
// Prunes.  After fetch if git repo is clean, a merge is performed
type reposync struct {
	Directories    []string
	RemoteRepoMeta []RemoteRepoMeta
	Prune          bool
	Output         chan ReposyncMessage
}

type RemoteRepoMeta struct {
	Url string
	WorkingDirectory string
}

type ReposyncMessage struct {
	Repo *git.Repo
	Output string
}

// create our resource pool of 10 git resources
var gitResourcePool = pool.NewGitResourcePool(10)

func NewReposync(directories []string, remoteRepoMeta []RemoteRepoMeta, prune bool) *reposync {

	r := reposync{
		Directories: directories,
		RemoteRepoMeta:    remoteRepoMeta,
		Prune:       prune,
	}

	// output from parallel git operations
	r.Output = make(chan ReposyncMessage)

	return &r
}

// Starts sync of all git repos
func (r *reposync) Sync() int {

	localReposCount := len(r.Directories)

	// pull all git folders already in the output directory
	for _, dir := range r.Directories {
		repoDir, repoName := filepath.Split(dir)
		repo := git.NewRepo(repoDir, repoName, "")

		// don't attempt a pull if the repo already exists
		// or is symlinked to the working directory
		if !repo.IsGitRepo() {
			localReposCount--
			continue
		}

		go func(repo git.Repo) {
			// get a resource out of the pool
			resource := gitResourcePool.Borrow()

			// use the resource to perform a git pull
			r.Output <- ReposyncMessage{
				&repo,
				resource.Pull(repo, r.Prune),
			}

			// give the resource back
			gitResourcePool.Return(resource)

			return
		}(repo)
	}

	// clone all repos from upstream
	for _, repoUrl := range r.RemoteRepoMeta {
		go func(repo git.Repo) {

			// get a resource out of the pool
			resource := gitResourcePool.Borrow()

			// use the resource to perform a git clone
			r.Output <- ReposyncMessage{
				&repo,
				resource.Clone(repo),
			}

			// give the resource back
			gitResourcePool.Return(resource)

			return
		}(git.NewRepo(repoUrl.WorkingDirectory, "", repoUrl.Url))
	}

	return localReposCount + len(r.RemoteRepoMeta)
}
