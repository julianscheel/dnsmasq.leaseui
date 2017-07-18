package leaseparsers

import (
	"github.com/klauspost/oui"
)

type Lease struct {
	Expiry int64
	Mac string
	MacVendor string
	Ip string
	Hostname string
	ClientId string
}

func GetVendorByMac(ouiDb oui.DynamicDB, mac string) string {
	if ouiDb == nil {
		return ""
	}

	macEntry, err := ouiDb.Query(mac)
	if err == nil {
		return macEntry.Manufacturer
	}

	return ""
}
