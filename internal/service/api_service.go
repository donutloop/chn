package service

import (
	"github.com/donutloop/chn/internal/handler"
	"log"
	"net/http"
	"strconv"
)

func NewAPIService(port int) *APIService {
	return &APIService{port: port}
}

type APIService struct {
	port int
}

func (s *APIService) Init() error {

	http.Handle("/xservice/service.chn.StoryService/Stories", handler.NewStoryServiceServer(NewStoriesService(), nil))

	// serve the favicon and logo files
	http.HandleFunc("/favicon.ico", handler.File("favicon"))
	http.HandleFunc("/logo.gif", handler.File("logo"))
	http.HandleFunc("/bundle.js", handler.File("bundle.js"))
	http.HandleFunc("/main.css", handler.File("main.css"))
	http.HandleFunc("/", handler.File("index"))

	// start the server up on our port
	log.Printf("Running on %d\n", s.port)
	err := http.ListenAndServe(":"+strconv.Itoa(s.port), nil)
	if err != nil {
		return err
	}

	return nil
}
