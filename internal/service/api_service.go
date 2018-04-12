package service

import (
	"github.com/donutloop/chn/internal/api"
	"github.com/donutloop/chn/internal/cache"
	"github.com/donutloop/chn/internal/client"
	"github.com/donutloop/chn/internal/handler"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
	"github.com/donutloop/chn/internal/mediator"
	"github.com/donutloop/chn/internal/scraper"
	"github.com/donutloop/chn/internal/storage"
)

func NewAPIService(config *api.Config) *APIService {
	return &APIService{
		config: config,
	}
}

type APIService struct {
	config *api.Config
	Router *chi.Mux
}

func (s *APIService) Init() error {

	// components ...
	hn := client.NewHackerNews(s.config.HackerNews.BaseURL, s.config.TimeoutAfter)
	githubClient := client.NewGithub(s.config.Github.BaseURL, s.config.TimeoutAfter)
	githubScraper := scraper.NewGithubScraper()

	githubMediator := mediator.NewGithub(githubClient, githubScraper, s.config.Github.BaseURL, s.config.TimeoutAfter)

	storiesCache := cache.NewStoriesCache(s.config.StoriesCache.DefaultExpirationInMinutes, s.config.StoriesCache.CleanupIntervalInMinutes)

	st, err := storage.New(s.config)
	if err != nil {
		return err
	} else {
		log.Infof("storage is connected (%s)", s.config.Storage.Address)
	}

	// router and middleware ...
	r := chi.NewRouter()

	r.Use(
		middleware.DefaultLogger,
		middleware.Timeout(s.config.TimeoutAfter*time.Second),
		middleware.Recoverer,
	)

	// services ...
	storiesService := NewStoriesService(hn, storiesCache, githubMediator, st)

	// routes ...
	r.Method(http.MethodPost, "/xservice/service.chn.StoryService/Stories", handler.NewStoryServiceServer(storiesService, nil, log.Errorf))
	handler.FileServer(r, "/static", http.Dir("../../static"))
	r.Get("/", handler.File("index"))

	s.Router = r
	return nil
}

func (s *APIService) Start(port int) error {
	// start the server up on our port
	return http.ListenAndServe(":"+strconv.Itoa(port), s.Router)
}
