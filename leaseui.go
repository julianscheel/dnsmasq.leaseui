package main

import (
	"io/ioutil"
	"bytes"
	"flag"
	"fmt"
	"os"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/klauspost/oui"
	"github.com/jasonlvhit/gocron"
	"./leaseparsers"
)

type page struct {
	Title string
}

func updateOuiDb(db oui.DynamicDB) {
	fmt.Println("Start oui database update...")
	oui.UpdateHttp(db, "http://standards-oui.ieee.org/oui.txt")
	fmt.Println("Done.")
}

func main() {
	db, err := oui.OpenFile("oui.txt")
	if err != nil {
		/* No local cache exists, create empty database */
		io := ioutil.NopCloser(bytes.NewReader(nil))
		db, err = oui.Open(io)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	/* Start oui update in background */
	go updateOuiDb(db)

	/* Schedule weekly auto update */
	gocron.Every(1).Day().Do(updateOuiDb, 1, db)
	_, time := gocron.NextRun()
	fmt.Println(time)
	gocron.Start()

	app := iris.New()
	app.RegisterView(iris.HTML("./templates/web/default", ".html"))
	app.Get("/leases", func(ctx context.Context) {
		ctx.JSON(map[string]interface{}{ "data": leaseparsers.ParseDnsmasqLeases(db) } )
	app.StaticWeb("/css", "./static/css")
	app.StaticWeb("/js", "./static/js")

	app.Get("/", func(ctx context.Context) {
		ctx.ViewData("Title", "DHCP Leases")
		ctx.View("leases.html")

		if err != nil {
			println(err.Error())
		}
	})

	app.Run(iris.Addr(":8080"))
}
