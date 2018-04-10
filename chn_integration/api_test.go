package chn_integration

import (
	"testing"
	"net/http/httptest"
	"os"
	"github.com/donutloop/chn/internal/service"
	"github.com/BurntSushi/toml"
	"github.com/donutloop/chn/internal/api"
	"log"
	"fmt"
	"github.com/donutloop/chn/internal/handler"
	"net/http"
	"context"
)

var url string

func TestMain(m *testing.M) {

	config := &api.Config{}
	if _, err := toml.DecodeFile("../cfg/config_local.toml", config); err != nil {
		log.Fatal(fmt.Sprintf("couldn't load config (%v)", err))
	}

	// services ...
	apiService := service.NewAPIService(config)
	apiService.Init()

	server := httptest.NewServer(apiService.Router)
	defer server.Close()
	url = server.URL

	// call flag.Parse() here if TestMain uses flags
	os.Exit(m.Run())
}

func TestStoriesJSONCall(t *testing.T) {

	client := handler.NewStoryServiceJSONClient(url, new(http.Client))

	storyResp, err := client.Stories(context.Background(), &handler.StoryReq{Category: "best", Limit:10})
	if err != nil {
		t.Fatal(err)
	}

	if len(storyResp.Stories) == 0 {
		t.Fatal("stories is empty")
	}

	var langaugesCounter int
	for _, story := range storyResp.Stories {
		if len(story.Langauges) != 0 {
			langaugesCounter++
		}
	}

	if langaugesCounter == 0 {
		t.Error("langauges counter is zero")
	}
}
