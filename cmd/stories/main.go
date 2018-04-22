package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	stdlog "log"
	"os"
	"github.com/donutloop/chn/storiesservice"
)

func main() {

	stdlog.SetFlags(stdlog.Lshortfile | stdlog.Ldate | stdlog.Ltime)

	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "addr",
			Usage: "server is listing on addr",
		},
		cli.StringFlag{
			Name:  "env",
			Value: "DEV",
		},
		cli.StringFlag{
			Name: "etcd.addr",
			Usage: "etcd registry address",
		},
		cli.StringFlag{
			Name: "db.addr",
		},
	}

	app.Action = func(c *cli.Context) error {
		addr := c.GlobalString("addr")

		apiService := storiesservice.NewAPIService(c.GlobalString("etcd.addr"), c.GlobalString("env"), c.GlobalString("db.addr"))

		if err := apiService.Init(); err != nil {
			return err
		}

		if err := apiService.Start(addr); err != nil {
			return err
		}

		log.Infof("Running on %s", addr)
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.WithError(err).Error("couldn't start api service")
	}
}
