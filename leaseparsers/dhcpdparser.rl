package leaseparsers

import (
	"bytes"
	"time"
	"fmt"
	"io"
	"os"
	"github.com/klauspost/oui"
)

%%{
	machine lease;
	write data;
}%%

func ParseDhcpdLeases(ouiDb oui.DynamicDB) []Lease {
	file, err := os.Open("/var/lib/dhcpd/dhcpd.leases")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	buf := bytes.NewBuffer(nil)
	io.Copy(buf, file)
	data := string(buf.Bytes())

	cs, p, pe, eof := 0, 0, len(data), len(data)
	var second, hour, minute int
	var day, month, year int
	var str string
	var date time.Time
	var val int

	var entry Lease
	var leases []Lease

	%%{
		action date_complete { date = time.Date(year, time.Month(month), day, hour, minute, second, 0, time.UTC) }
		action entry_complete {
			entry.MacVendor = GetVendorByMac(ouiDb, entry.Mac)
			leases = append(leases, entry)
			entry = Lease{}
		}

		crlf = '\r'? '\n';
		pad = ( crlf | space );

		num = ( digit @{ val = val * 10 + (int(fc) - '0') } );
		action stringify { str = str + string(fc) }
		strval = ( ( alnum | ( punct - [";] ) ) @stringify )+ >{ str = "" };

		second = num+ % { second = val; val = 0; };
		minute = num+ % { minute = val; val = 0; };
		hour = num+ % { hour = val; val = 0; };
		day = num+ % { day = val; val = 0; };
		month = num+ % { month = val; val = 0; };
		year = num+ % { year = val; val = 0; };
		weekday = alnum+;

		date = ( weekday space+ year '/' month '/' day space+ hour ':' minute ':' second ) % date_complete;

		start = ( "starts" space+ date );
		end = ( "ends" space+ date ) % { entry.Expiry = int64(date.Unix()) };
		tstp = ( "tstp" space+ date );
		tsfp = ( "tsfp" space+ date );
		atsfp = ( "atsfp" space+ date );
		cltt = ( "cltt" space+ date );
		times = ( start | end | tstp | tsfp | atsfp | cltt );

		hostname = ( "client-hostname" space+ '"' strval '"' ) %{ entry.Hostname = str };
		ethernet = ( "hardware" space+ "ethernet" space+ strval ) %{ entry.Mac = str };
		uid = ( "uid" space+ '"' strval '"' ) %{ entry.ClientId = str };
		unknown = any+;

		parameter = ( space* ( times | hostname | ethernet | uid | unknown ) ';' space* );

		address = strval %{ entry.Ip = str };

		entry = ( "lease" pad+ address pad* '{' pad* parameter+ pad*'}' ) % entry_complete;
		comment = ( '#' any* crlf );
		main := ( pad* ( entry | comment ) pad* )*;

		write init;
		write exec;
	}%%

	if cs < lease_first_final {
		fmt.Println("dhcpd.leases parse error")
	}

	return leases;
}
