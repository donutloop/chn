package client

import (
	"fmt"
	"encoding/json"
	"net/http"
	"time"
	"github.com/pkg/errors"
)

func NewGithub(baseURL string, timeoutAfter time.Duration) *Github {
	return &Github{
		baseURL:      baseURL,
		timeoutAfter: timeoutAfter,
	}
}

type Github struct {
	baseURL string
	client http.Client
	timeoutAfter time.Duration
}


// ListsLanguages for the specified repository. The value shown for each language is the number of bytes of code written in that language.
func (c Github) ListsLanguages(owner string, repo string) (map[string]int, error) {

	req, cancel, err := newRequestWihtoutTimeout(http.MethodGet, fmt.Sprintf("%s/repos/%s/%s/languages", c.baseURL, owner, repo), nil, c.timeoutAfter)
	if err != nil {
		return nil, err
	}
	defer cancel()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, errors.Wrap( err, "error while preparing request is a error occurred")
	}
	defer resp.Body.Close()

	languages := make(map[string]int)
	if err := json.NewDecoder(resp.Body).Decode(&languages); err != nil {
		return nil, errors.Wrap(err,"error decoding json has failed")
	}
	return languages, nil
}
