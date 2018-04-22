package storiesservice

import (
	"github.com/donutloop/chn/storiesservice/internal/cache"
	"github.com/donutloop/chn/storiesservice/internal/handler"
	"github.com/donutloop/chn/storiesservice/internal/mediator"
	"github.com/donutloop/chn/storiesservice/internal/scraper"
	"github.com/donutloop/chn/storiesservice/internal/service"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	log "github.com/sirupsen/logrus"
	"time"
	c "github.com/coreos/etcd/client"
	"context"
	"github.com/BurntSushi/toml"
	"bytes"

	"github.com/donutloop/chn/storiesservice/internal/config"
	"github.com/donutloop/chn/storiesservice/internal/client"
	"net/http"
	"github.com/pkg/errors"
	"github.com/go-chi/cors"
	"github.com/donutloop/chn/storiesservice/internal/storage"
)

func NewAPIService(etcdAddr string, env string, dbAddr string) *APIService {
	return &APIService{
		etcdAddr: etcdAddr,
		env: env,
		dbAddr:dbAddr,
	}
}

type APIService struct {
	config *config.Config
	etcdAddr string
	Router *chi.Mux
	env    string
	dbAddr string
}

func (s *APIService) Init() error {

	cetcd, err := c.New(c.Config{
		Endpoints:   []string{s.etcdAddr},
		HeaderTimeoutPerRequest: 5 * time.Second,
	})
	if err != nil {
		return err
	}
	kapi := c.NewKeysAPI(cetcd)

	configResponse, err := kapi.Get(context.Background(), "/stories-config", nil)
	if err != nil {
		return err
	}

	config := &config.Config{}
	_, err = toml.DecodeReader(bytes.NewReader([]byte(configResponse.Node.Value)), config)
	if err != nil {
		return err
	}
	s.config = config
	config.Storage.Address = s.dbAddr

	setupLogger(s.env)

	// components ...
	hn := client.NewHackerNews(s.config.HackerNews.BaseURL, s.config.TimeoutAfter)
	githubClient := client.NewGithub(s.config.Github.BaseURL, s.config.TimeoutAfter)
	githubScraper := scraper.NewGithubScraper()

	githubMediator := mediator.NewGithub(githubClient, githubScraper, s.config.Github.BaseURL, s.config.TimeoutAfter)

	storiesCache := cache.NewStoriesCache(s.config.StoriesCache.DefaultExpirationInMinutes, s.config.StoriesCache.CleanupIntervalInMinutes)


	st, err := storage.New(s.config)
	if err != nil {
		return errors.Wrap(err, "mongo db")
	} else {
		log.Infof("storage is connected (%s)", s.config.Storage.Address)
	}

	storiesStorage := storage.NewStories(st, s.config.StoriesStorage.Tries, s.config.StoriesStorage.InitialInterval, s.config.StoriesStorage.MaxInterval)

	// router and middleware ...
	r := chi.NewRouter()

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})

	r.Use(
		middleware.DefaultLogger,
		middleware.Timeout(s.config.TimeoutAfter*time.Second),
		middleware.Recoverer,
		cors.Handler,
	)

	// services ...
	storiesService := service.NewStoriesService(hn, storiesCache, githubMediator, storiesStorage)

	// routes ...
	r.Method(http.MethodPost, "/xservice/service.chn.StoryService/Stories", handler.NewStoryServiceServer(storiesService, nil, log.Errorf))

	s.Router = r
	return nil
}

func (s *APIService) Start(addr string) error {
	// start the server up on our port
	return http.ListenAndServe(addr, s.Router)
}

func setupLogger(env string) {
	switch env {
	case "DEV":
		log.SetLevel(log.DebugLevel)
		log.SetFormatter(&log.TextFormatter{})
	default:
		log.SetFormatter(&log.JSONFormatter{})
	}
}
