package scraper

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
)

func NewGithubScraper() *Github {
	return &Github{}
}

type Github struct{}

func (g *Github) ExtractBaseURL(url string) (string, error) {
	// todo check if is an github url

	doc, err := goquery.NewDocument(url)
	if err != nil {
		return "", err
	}

	var baseUrl string
	doc.Find(".repohead-details-container a").Each(func(i int, s *goquery.Selection) {

		attr, ok := s.Attr("data-pjax")
		if !ok {
			return
		}

		if attr != "#js-repo-pjax-container" {
			return
		}

		subPath, ok := s.Attr("href")
		if !ok {
			return
		}

		baseUrl = "https://github.com" + subPath
	})

	if baseUrl == "" {
		return "", errors.New("url not found")
	}

	return baseUrl, nil
}
