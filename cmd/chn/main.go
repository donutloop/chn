package main

import (
	"github.com/donutloop/chn/internal/service"
	"github.com/urfave/cli"
	"os"
	log "github.com/sirupsen/logrus"
	stdlog "log"
)

func main() {

	stdlog.SetFlags(stdlog.Lshortfile| stdlog.Ldate | stdlog.Ltime)

	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "port",
			Value: "8080",
			Usage: "server is listing on port",
		},
	}

	app.Action = func(c *cli.Context) error {
		port := c.GlobalInt("port")
		api := service.NewAPIService(port)

		if err := api.Init(); err != nil {
			return err
		}

		log.Infof("Running on %d", port)
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.WithError(err).Error("could start api service")
	}
}
