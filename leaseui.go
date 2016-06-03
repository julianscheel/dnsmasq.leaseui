package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"github.com/kataras/iris"
)

type Lease struct {
	Expiry int
	Mac string
	Ip string
	Hostname string
	ClientId string
}

func parseLeases() []Lease {
	var csvFile, err = os.Open("/var/lib/misc/dnsmasq.leases")
	if err != nil {
		fmt.Println(err)
		return nil
	}

	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	reader.FieldsPerRecord = 5
	reader.Comma = ' '

	var csvData, cerr = reader.ReadAll()
	if cerr != nil {
		fmt.Println(err)
		return nil
	}

	var lease Lease
	var leases []Lease

	for _, entry := range csvData {
		lease.Expiry, err = strconv.Atoi(entry[0])
		lease.Mac = entry[1]
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

func main() {
	iris.Config().Render.Template.Directory = "./templates/web/default"

	iris.Get("/leases", func(c *iris.Context) {
		c.JSON(200, iris.Map{ "data": parseLeases() } )
	})
	iris.Static("/css", "./static/css", 1)
	iris.Static("/js", "./static/js", 1)

	iris.Get("/", func(ctx *iris.Context) {
		err := ctx.Render("leases.html", page{Title: "DHCP Leases"})
		if err != nil {
			println(err.Error())
		}
	})

	iris.Listen(":8080")
}
