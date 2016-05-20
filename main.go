package main

import (
	"flag"
	"fmt"
	"github.com/justin-robinson/reposync/gitlab"
	"github.com/justin-robinson/reposync/lib/config"
	"log"
	"os/user"
	"strconv"
	"github.com/justin-robinson/reposync/lib/reposync"
	"github.com/justin-robinson/reposync/lib"
	"github.com/justin-robinson/reposync/lib/filters"
	"path/filepath"
	"io/ioutil"
	"strings"
	"os"
)

func main() {

	workingDirectory := flag.String("dir", "", "path to git repos")
	prune := flag.Bool("prune", false, "prune on git fetch")
	version := flag.Bool("version", false, "prints version number")

	// parse cli args
	flag.Parse()

	if *version {
		fmt.Println(lib.GetAppVersion())
		os.Exit(0)
	}

	// sync ~/dev/code if no directoy was specified
	if *workingDirectory == "" {
		usr, err := user.Current()

		if err != nil {
			log.Fatal(err)
		}

		*workingDirectory = usr.HomeDir + "/dev/code"
	}

	// load the config file
	configFile := config.GetConfig(*workingDirectory)

	// urls for remote repos
	repoUrls := make([]reposync.RemoteRepoMeta, 0)

	// exhaust all pages of gitlab api for projects
	apiPageNumber := 0
	api := gitlab.Api{
		configFile.Source.Url,
		configFile.Source.Token,
	}
	for {
		// get the next page!
		apiPageNumber++

		// get list of projects
		projects, err := api.GetProjects(strconv.Itoa(apiPageNumber))

		// log api errors
		if err != nil {
			log.Print(err)
			break
		}

		// the last page will be an empty response
		if len(projects) == 0 {
			break
		}

		// add each url to our list of repos to clone
		for _, item := range projects {

			repoPath := filepath.Join(*workingDirectory, item.Path)

			// is the repo whitelisted by an existing config?
			whitelisted := filters.RepoIsWhitelisted(configFile.Source.Repos, item.Path_with_namespace)

			// does this repo exist already?
			existsLocally := filters.RepoExistsLocally(repoPath)

			// append repo to our list of repos to clone
			if whitelisted && !existsLocally {
				repoUrls = append(repoUrls, reposync.RemoteRepoMeta{item.Ssh_url_to_repo, *workingDirectory})
			}
		}
	}

	// get all files in the working directory
	fileNames, err := ioutil.ReadDir(*workingDirectory)
	if err != nil {
		log.Fatal(err)
	}

	// get all non symlinked folders in the working directory
	directories := make([]string, 0)
	for _, file := range fileNames {
		repoPath := filepath.Join(*workingDirectory, file.Name())
		isSymlink := filters.RepoIsSymlinkedToRepo(*workingDirectory, repoPath)
		if !isSymlink && file.IsDir() {
			directories = append(directories, repoPath)
		}
	}

	// create the repo syncer
	reposync := reposync.NewReposync(directories, repoUrls, *prune)

	// SYNC!
	messagesReceived, messagesExpected := 0, reposync.Sync()
	for messagesReceived < messagesExpected {
		// print output for each git operation
		messagesReceived++
		message := <-reposync.Output

		var background string
		if message.Repo.IsClean() {
			background = lib.OnGreen
		} else {
			background = lib.OnYellow
		}

		output := lib.Black + background + strings.ToTitle(message.Repo.RepoName) + lib.ColorOff +
			" (" + strconv.Itoa(messagesReceived) + "/" + strconv.Itoa(messagesExpected) + ")\n"

		fmt.Print( output + message.Output )
	}

	return
}
