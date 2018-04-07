package main

import (
	"github.com/donutloop/chn/internal/service"
	"github.com/urfave/cli"
	"log"
	"os"
)

func main() {

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

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
