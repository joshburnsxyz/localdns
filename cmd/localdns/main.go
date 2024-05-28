package main

import (
	"fmt"
	"net"

	"github.com/google/gopacket"
	layers "github.com/google/gopacket/layers"
	flag "github.com/spf13/pflag"
	"github.com/joshburnsxyz/localdns/pkg/record"
)

var (
	csvFileFlag string
	portFlag int
	hostFlag string
)

func init() {
	flag.StringVarP(&csvFileFlag, "csv", "c", "./dns.csv", "DNS map file")
	flag.IntVarP(&portFlag, "port", "p", 53, "Port to bind server too (default 53)")
	flag.StringVarP(&hostFlag, "host", "H", "0.0.0.0", "Host interface to bind too (default 0.0.0.0)")
}

func main() {
	flag.Parse()

	// initialize records
	recordsmap := record.NewRecordsFromCSV(csvFileFlag)

	// initialize UDP Server
	laddr := net.UDPAddr{
		Port: portFlag,
		IP: net.ParseIP(hostFlag),
	}
	u, _ := net.ListenUDP("udp", &laddr)

	// Listen for DNS requests
	for {
		tmp := make([]byte, 1024)
		_, addr, _ := u.ReadFrom(tmp)
		clientAddr := addr
		packet := gopacket.NewPacket(tmp, layers.LayerTypeDNS, gopacket.Default)
		dnsPacket := packet.Layer(layers.LayerTypeDNS)
		tcp, _ := dnsPacket.(*layers.DNS)
		serveDNS(u, clientAddr, tcp)
		
	}
}
