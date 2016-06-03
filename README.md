# dnsmasq lease ui

A minimal webinterface to view dnsmasq leases. Data is served as JSON from a
small go application.
The frontend makes use of jQuery datatables (datatables.net)

## Getting started

Assuming you have a running dnsmasq installation just run:
```
go run leaseui.go
```
And open the interface in your browser: [http://localhost:8080/]
