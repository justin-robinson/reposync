package gitlab

import (
	"encoding/json"
	"net/http"
	"net/url"
	"github.com/justin-robinson/reposync/gitlab/responses"
)

type Api struct {
	Url string
	Token string
}

func (a *Api) GetProjects(pageNumber string) ([]responses.Projects, error) {

	// build query params
	queryParams := a.getBaseQueryParams()
	queryParams.Set("per_page", "100")
	queryParams.Set("page", pageNumber)

	// make api call
	response, err := http.Get(a.Url + "projects?" + queryParams.Encode())

	defer response.Body.Close()

	if err != nil {
		return nil, err
	}

	// decode the response body to json
	projects := make([]responses.Projects, 0)
	jsonDecoder := json.NewDecoder(response.Body)
	jsonDecoder.Decode(&projects)

	return projects, nil
}

func (a *Api) getBaseQueryParams() url.Values {
	queryParams := url.Values{}
	queryParams.Set("private_token", a.Token)
	return queryParams
}
