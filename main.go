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
	app.Version = "0.2"

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "addr, a", Value: "localhost", Usage: "MongoDb host to connect to"},
		cli.IntFlag{Name: "port, p", Value: 27017, Usage: "MongoDb port to connect to"},
		cli.StringFlag{Name: "user, u", Value: "", Usage: "username to access MongoDb"},
		cli.StringFlag{Name: "password, P", Value: "", Usage: "password to access MongoDb"},
		cli.StringFlag{Name: "gridfs, g", Value: "fs", Usage: "GridFS prefix"},
	}

	app.Action = func(c *cli.Context) {

		if len(c.Args()) != 2 {
			log.Fatal("Usage: " + appName + " <dbname> <mountpoint>")
			os.Exit(2)
		}

		dbName := c.Args()[0]
		mountPoint := c.Args()[1]
		dbHost := c.String("addr")
		dbPort := string(c.String("port"))
		dbUser := c.String("user")
		dbPassword := c.String("password")
		gridfsPrefix = c.String("gridfs")
		credentials := dbUser + ":" + dbPassword

		// Connect to the database
		initDb(dbName, dbHost, dbPort, credentials)

		// Mount the database
		err := mount(mountPoint, appName)
		checkErrorAndExit(err, 1)
	}

	app.Run(os.Args)
}
