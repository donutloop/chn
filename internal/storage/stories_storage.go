package storage

import (
	"github.com/donutloop/chn/internal/model"
	"github.com/globalsign/mgo"
	"github.com/donutloop/chn/internal/handler"
	"github.com/donutloop/toolkit/retry"
	"context"
)

type Stories struct {
	storage Interface
	retrier retry.Retrier
}

func NewStories(storage Interface, maxTries uint, initialInterval, maxInterval float64) (*Stories) {

	r := retry.NewRetrier()
	r.InitialInterval = initialInterval
	r.MaxInterval = maxInterval
	r.Tries = maxTries

	return &Stories{
		storage: storage,
	}
}

func (s *Stories) SaveNewStoriesWithRetry(stories []*handler.Story) ([]error) {
	storiesFailed := make([]*model.Story, 0)

	for _, story := range stories {
		if err := s.storage.FindBy("url", story.Url,  new(model.Story)); err == mgo.ErrNotFound {
			ms := model.NewStoryFrom(story)
			if err := s.storage.Insert(ms); err != nil {
				storiesFailed = append(storiesFailed, ms)
			}
		} else if err != nil {
			return []error{err}
		}
	}

	errs := make([]error, 0)
	for _, story := range storiesFailed {
		err := s.retrier.Retry(context.Background(), func() (bool, error) {
			if err := s.storage.Insert(story); err != nil {
				return false, err
			}
			return true, nil
		})

		if err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

