# leases-cli

This is a command line viewer for DHCP leases provided by the leaseui backend
service.

## Getting started

Usage is as simple as
```
./leases-cli --uri http://server:8080/leases hostname-filter
```

For convenience you may specify an alias in your environment to avoid
specifying the uri on each call.
```
alias leases "leases-cli --uri http://server:8080/leases"
```
