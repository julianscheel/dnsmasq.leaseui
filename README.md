# dnsmasq lease ui

A minimal webinterface to view DHCP leases managed by dnsmasq or bind dhcpd.
Data is served as JSON by a small go service.
The frontend makes use of [jQuery datatables](datatables.net)

## Getting started

Generate and run leaseui as following. You need to have
[ragel](http://www.colm.net/open-source/ragel/) installed to generated some of
the parsers. The default configuration will work with
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

## Command-line

A command-line client is provided in the `leases-cli` subfolder.
