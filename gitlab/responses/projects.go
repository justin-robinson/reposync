package responses

// Project is the expected response from gitlab's projects api
// http://doc.gitlab.com/ee/api/projects.html
type Projects struct {
	Id               int
	Description      string
	Default_branch   string
	Tag_list         []string
	Public           bool
	Archived         bool
	Visibility_level int
	Ssh_url_to_repo  string
	Http_url_to_repo string
	Web_url          string
	Owner            struct {
		Name       string
		Username   string
		Id         int
		State      string
		Avatar_url string
		Web_url    string
	}
	Name                   string
	Name_with_namespace    string
	Path                   string
	Path_with_namespace    string
	Issues_enabled         bool
	Merge_requests_enabled bool
	Wiki_enabled           bool
	Snippets_enabled       bool
	Created_at             string
	Last_activity_date     string
	Creator_id             int
	Namespace              struct {
		Id          int
		Name        string
		Path        string
		Owner_id    int
		Created_at  string
		Updated_at  string
		Description string
		Avatar      string
	}
	Avatar_url  string
	Star_count  int
	Forks_count int
}
