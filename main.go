package main

import (
	"log"
	"os"

	"github.com/codegangsta/cli"
)

var appName string = "mgfs"
var gridfsPrefix string

func main() {
	log.SetFlags(0)
	log.SetPrefix(appName + ": ")

	app := cli.NewApp()
	app.Name = appName
	app.Usage = "mount a mongodb database"
	app.Version = "0.1"

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "dbhost, d", Value: "localhost", Usage: "MongoDb host to connect to"},
		cli.StringFlag{Name: "prefix, p", Value: "fs", Usage: "GridFS prefix"},
	}

	app.Action = func(c *cli.Context) {

		if len(c.Args()) != 2 {
			log.Fatal("Usage: " + appName + " <dbname> <mountpoint>")
			os.Exit(2)
		}

		dbName := c.Args()[0]
		mountPoint := c.Args()[1]
		dbHost := c.String("dbhost")
		gridfsPrefix = c.String("prefix")

		// Connect to the database
		initDb(dbHost, dbName)

		// Mount the database
		err := mount(mountPoint, appName)
		checkErrorAndExit(err, 1)
	}

	app.Run(os.Args)
}
