package issues

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type httpClient interface {
	Get(url string) (*http.Response, error)
}

type Issues struct {
	Number int
	Title  string
}

type Label struct {
	Name string
}

type Service struct {
	repo   string
	client httpClient
}

func NewService(repo string, client httpClient) Service {
	return Service{
		repo:   repo,
		client: client,
	}
}

func (i Service) Get(labelFilter, creatorFilter string) ([]Issues, error) {
	rawQueries := []string{}
	if labelFilter != "" {
		rawQueries = append(rawQueries, fmt.Sprintf("labels=%s", url.QueryEscape(labelFilter)))
	}

	if creatorFilter != "" {
		rawQueries = append(rawQueries, fmt.Sprintf("creator=%s", url.QueryEscape(creatorFilter)))
	}

	u := &url.URL{
		Scheme:   "https",
		Path:     fmt.Sprintf("api.github.com/repos/%s/issues", i.repo),
		RawQuery: strings.Join(rawQueries, "&"),
	}

	res, err := i.client.Get(u.String())
	if err != nil {
		return []Issues{}, err
	}

	if res.StatusCode != 200 {
		return []Issues{}, fmt.Errorf("invalid path %s", u.String())
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []Issues{}, err
	}

	var issues []Issues
	err = json.Unmarshal(body, &issues)
	if err != nil {
		return []Issues{}, err
	}

	return issues, nil
}
