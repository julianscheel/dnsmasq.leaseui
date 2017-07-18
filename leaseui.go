package main

import (
	"encoding/csv"
	"io/ioutil"
	"bytes"
	"fmt"
	"os"
	"strconv"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/klauspost/oui"
	"github.com/jasonlvhit/gocron"
)

type Lease struct {
	Expiry int
	Mac string
	MacVendor string
	Ip string
	Hostname string
	ClientId string
}

func parseLeases(ouiDb oui.DynamicDB) []Lease {
	csvFile, err := os.Open("/var/lib/misc/dnsmasq.leases")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	reader.FieldsPerRecord = 5
	reader.Comma = ' '

	csvData, err := reader.ReadAll()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	var lease Lease
	var leases []Lease

	for _, entry := range csvData {
		lease.Expiry, err = strconv.Atoi(entry[0])
		lease.Mac = entry[1]
		if ouiDb != nil {
			macEntry, err := ouiDb.Query(lease.Mac)
			if err == nil {
				lease.MacVendor = macEntry.Manufacturer
			}
		}
		lease.Ip = entry[2]
		lease.Hostname = entry[3]
		lease.ClientId = entry[4]

		leases = append(leases, lease)
	}

	return leases
}

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
		ctx.JSON(map[string]interface{}{ "data": parseLeases(db) } )
	})
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
