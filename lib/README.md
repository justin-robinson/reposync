```go
PACKAGE DOCUMENTATION

package lib
    import "."


CONSTANTS

const CONFIG_FILE_NAME = ".reposync.json"

FUNCTIONS

func GetConfig(dir string) (reposyncConfig, bool)
    reads and parses .reposync.json file inside of dir

func NewReposync(Directory string, RepoUrls []string, Prune bool) *reposync
```
