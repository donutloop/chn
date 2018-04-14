package model

import "github.com/donutloop/chn/internal/handler"

func NewStoryFrom(story *handler.Story) *Story {
	return &Story{
		By: story.By,
		Descendants: story.Descendants,
		Kids: story.Kids,
		Score: story.Score,
		Type: story.Type,
		Title: story.Title,
		Url: story.Url,
		DomainName: story.DomainName,
		Langauges: story.Langauges,
	}
}

type Story struct {
	ID          string   `bson:"id,omitempty" json:"id,omitempty"`
	By          string   `bson:"by,omitempty" json:"by,omitempty"`
	Descendants int64    `bson:"descendants,omitempty" json:"descendants,omitempty"`
	Kids        []int64  `bson:"kids,omitempty" json:"kids,omitempty"`
	Score       int64    `bson:"score,omitempty" json:"score,omitempty"`
	Type        string   `bson:"type,omitempty" json:"type,omitempty"`
	Title       string   `bson:"title,omitempty" json:"title,omitempty"`
	Url         string   `bson:"url,omitempty" json:"url,omitempty"`
	DomainName  string   `bson:"domainName,omitempty" json:"domainName,omitempty"`
	Langauges   []string `bson:"Langauges,omitempty" json:"Langauges,omitempty"`
}

func (s *Story) GetId() string {
	return s.ID
}

func (s *Story) SetId(id string) {
	s.ID = id
}

func (s *Story) GetNamespace() string {
	return "stories"
}
