package main

import (
	"github.com/BurntSushi/toml"
	"github.com/donutloop/chn/internal/api"
	"github.com/donutloop/chn/internal/service"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	stdlog "log"
	"os"
)

func main() {

	stdlog.SetFlags(stdlog.Lshortfile | stdlog.Ldate | stdlog.Ltime)

	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "port",
			Value: 8080,
			Usage: "server is listing on port",
		},
		cli.StringFlag{
			Name:  "config",
			Value: "../../cfg/config_local.toml",
			Usage: "server is listing on port",
		},
		cli.StringFlag{
			Name:  "env",
			Value: "DEV",
		},
	}

	app.Action = func(c *cli.Context) error {
		port := c.GlobalInt("port")

		config := &api.Config{
			ENV: c.GlobalString("env"),
		}
		_, err := toml.DecodeFile(c.GlobalString("config"), config)
		if err != nil {
			return err
		}

		apiService := service.NewAPIService(config)

		if err := apiService.Init(); err != nil {
			return err
		}

		if err := apiService.Start(port); err != nil {
			return err
		}

		log.Infof("Running on %d", port)
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.WithError(err).Error("could start api service")
	}
}
