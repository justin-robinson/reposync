package git

func InitFlow (repo *Repo) string {
	return repo.RunCmd(repo.GetFullRepoPath(), "flow", "init", "-d")
}
