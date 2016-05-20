package filters

import "github.com/justin-robinson/reposync/lib/config"

func RepoIsWhitelisted (repos map[string]bool, repoNamespace string) bool {

	// all repos are whitelisted if the config file doesn't exist
	if !config.ConfigExists() {
		return true
	}

	// repo must be defined and have a value of true
	enabled, defined := repos[repoNamespace]
	return (defined && enabled)
}
