package client

import (
	"fmt"
	"encoding/json"
	"net/http"
	"context"
	"io"
	"time"
	"strconv"
)

// Post -- represents the json object returned by the API
type Post struct {
	Id          int    `json:"id"`
	Deleted     bool   `json:"deleted"`
	Type        string `json:"type"`
	By          string `json:"by"`
	Time        int    `json:"time"`
	Text        string `json:"text"`
	Dead        bool   `json:"dead"`
	Parent      int    `json:"parent"`
	Poll        int    `json:"poll"`
	Kids        []int  `json:"kids"`
	Url         string `json:"url"`
	Score       int64    `json:"score"`
	Title       string `json:"title"`
	Parts       []int  `json:"parts"`
	Descendants int    `json:"descendants"`
}

func NewHackerNews(baseURL string, timeoutAfter int) *HackerNews {
	return &HackerNews{
		baseURL:baseURL,
		timeoutAfter: timeoutAfter,
	}
}

type HackerNews struct {
	story   string
	baseURL string
	http.Client
	timeoutAfter int
}

// GetCodesStory -- Return the ids of a story
func (c HackerNews) GetCodesForStory(story string) ([]int, error) {

	req, cancel, err := newRequestWihtoutTimeout(http.MethodGet, fmt.Sprintf("%s%s.json", c.baseURL, story), nil, time.Duration(c.timeoutAfter))
	if err != nil {
		return nil, err
	}
	defer cancel()

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error while fetching codes story : %v", err)
	}
	defer resp.Body.Close()

	keys := make([]int, 0)
	if err := json.NewDecoder(resp.Body).Decode(&keys); err != nil {
		return nil, fmt.Errorf("error while decoding json from codes story : %v", err)
	}
	return keys, nil
}

// GetPostStory -- Return the posts of story thanks their ids
func (c HackerNews) GetPost(code int) (*Post, error) {
	req, cancel, err := newRequestWihtoutTimeout(http.MethodGet, fmt.Sprintf("%s/item/%s.json", c.baseURL,  strconv.Itoa(code)), nil, time.Duration(c.timeoutAfter))
	if err != nil {
		return nil, err
	}
	defer cancel()

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error while fetching codes story : %v", err)
	}
	defer resp.Body.Close()

	post := new(Post)
	if err := json.NewDecoder(resp.Body).Decode(&post); err != nil {
		return nil, fmt.Errorf("error while decoding json from codes story : %v", err)
	}
	return post, nil
}

func newRequestWihtoutTimeout(method string, url string, body io.Reader, timeoutAfter time.Duration) (*http.Request, context.CancelFunc, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeoutAfter*time.Second)
	req = req.WithContext(ctx)
	return req, cancel, err
}



