# localdns
A simple DNS server designed to run locally. Say goodbye to localhost

## Usage

1. Write a `dns.csv` file that matches a domain to an ip address

```csv
domain,ip
mylocalapp.local,192.168.10.10
otherlocalapp.local,192.168.10.20
```

2. Run the server under superuser, you may also have to point your systems DNS server settings to `localhost`. You can 
use the `--port <PORT>` option to bind to a non-priviliged port.

## Installation

1. Clone the repository and run `make` to build the binary.
