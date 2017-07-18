package leaseparsers

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"github.com/klauspost/oui"
)

func ParseDnsmasqLeases(ouiDb oui.DynamicDB) []Lease {
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
		lease.Expiry, err = strconv.ParseInt(entry[0], 10, 64)
		lease.Mac = entry[1]
		lease.MacVendor = GetVendorByMac(ouiDb, lease.Mac);
		lease.Ip = entry[2]
		lease.Hostname = entry[3]
		lease.ClientId = entry[4]

		leases = append(leases, lease)
	}

	return leases
}
