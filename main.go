package main

import (
	"log"
	"os"

	"github.com/codegangsta/cli"
)

var appName string = "mgfs"

func main() {
	app := cli.NewApp()
	app.Name = appName
	app.Usage = "mount a mongodb database"
	app.Version = "0.1"

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "dbhost, d", Value: "localhost", Usage: "the mongodb host to connect to"},
	}

	app.Action = func(c *cli.Context) {

		if len(c.Args()) != 2 {
			log.Fatal("Usage: " + appName + " <dbname> <mountpoint>")
		}

		dbName := c.Args()[0]
		dbHost := c.String("dbhost")
		mountPoint := c.Args()[1]

		initDb(dbHost, dbName)
		mount(mountPoint)
	}

	app.Run(os.Args)
}
