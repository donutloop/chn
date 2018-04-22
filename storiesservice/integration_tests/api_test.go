package integration_tests

import (
	"context"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/donutloop/chn/stories"
	"github.com/donutloop/chn/storiesservice/internal/api"
	"github.com/donutloop/chn/storiesservice/internal/handler"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var url string

func TestMain(m *testing.M) {

	config := &api.Config{}
	if _, err := toml.DecodeFile("../cfg/config_local.toml", config); err != nil {
		log.Fatal(fmt.Sprintf("couldn't load config (%v)", err))
	}

	// services ...
	apiService := stroies.NewAPIService(config)
	apiService.Init()

	server := httptest.NewServer(apiService.Router)
	defer server.Close()
	url = server.URL

	// call flag.Parse() here if TestMain uses flags
	os.Exit(m.Run())
}

func TestStoriesJSONCall(t *testing.T) {

	client := handler.NewStoryServiceJSONClient(url, new(http.Client))

	storyResp, err := client.Stories(context.Background(), &handler.StoryReq{Category: "best", Limit: 10})
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
