# dnsmasq lease ui

A minimal webinterface to view dnsmasq leases. Data is served as JSON from a
small go application.
The frontend makes use of jQuery datatables (datatables.net)

## Getting started

Generate and run leaseui as following. You need to have `ragel` installed to
generated some of the parsers. The default configuration will work with
dnsmasq.
```
go generate
go run leaseui.go
```
And open the interface in your browser: [http://localhost:8080/]

If you use bind dhcpd select the dhcpd backend:
```
go run leaseui.go --backend bind
```
