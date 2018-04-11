package mediator

import (

	"time"
	"github.com/donutloop/chn/internal/client"

	"github.com/donutloop/chn/internal/scraper"
	"sort"
	"strings"
	"github.com/pkg/errors"
	url2 "net/url"
)

func NewGithub(c *client.Github, s *scraper.Github, baseURL string, timeoutAfter time.Duration) *Github {
	return &Github{
		baseURL:      baseURL,
		timeoutAfter: timeoutAfter,
		s: s,
		c: c,
	}
}

type Github struct {
	baseURL      string
	timeoutAfter time.Duration
	c *client.Github
	s *scraper.Github
}

// todo extract more data
func (c Github) GetDataBy(url string) ([]string, error) {

	baseUrl, err := c.s.ExtractBaseURL(url)
	if err != nil {
		return nil, err
	}

	u, err := url2.Parse(baseUrl)
	if err != nil {
		return nil, errors.Wrapf(err, "error get stories")
	}

	pathParts := strings.Split(strings.TrimLeft(u.Path, "/"), "/")
	if len(pathParts) == 2 {
		raw, err := c.c.ListsLanguages(pathParts[0], pathParts[1])
		if err != nil {
			return nil, errors.Wrap(err, "github mediator")
		} else {
			languages := make([]string, 0)
			for l := range raw {
				languages = append(languages, l)
			}
			sort.Strings(languages)

			return languages, nil
		}
	}
	return nil, errors.Errorf("error get stories splitting github url has failed (len: %d)", len(pathParts))
}
