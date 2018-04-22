package frontendservice

import (
	"github.com/donutloop/chn/frontendservice/internal/handler"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	log "github.com/sirupsen/logrus"
	"net/http"
	"github.com/go-chi/cors"
)

func NewAPIService(etcdAddr string, env string, staticPath string) *APIService {
	return &APIService{
		etcdAddr: etcdAddr,
		ENV: env,
		staticPath: staticPath,
	}
}

type APIService struct {
	staticPath string
	ENV 	string
	etcdAddr string
	Router *chi.Mux
}

func (s *APIService) Init() error {

	setupLogger(s.ENV)

	// router and middleware ...
	r := chi.NewRouter()

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})

	r.Use(
		middleware.DefaultLogger,
		//	middleware.Timeout(s.config.TimeoutAfter*time.Second),
		middleware.Recoverer,
		cors.Handler,
	)

	handler.FileServer(r, "/static", http.Dir(s.staticPath))
	r.Get("/", handler.File(s.staticPath,"index"))

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
