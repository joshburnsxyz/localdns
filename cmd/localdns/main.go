package main

import (
	"fmt"
	"net"
	"os"
	"github.com/google/gopacket"
	layers "github.com/google/gopacket/layers"
	flag "github.com/spf13/pflag"
	"github.com/joshburnsxyz/localdns/pkg/record"
)

var (
	csvFileFlag string
	portFlag int
	hostFlag string
	listRecordsFlag bool
	dnsAnswer layers.DNSResourceRecord
)

func init() {
	flag.StringVarP(&csvFileFlag, "csv", "c", "./dns.csv", "DNS map file")
	flag.IntVarP(&portFlag, "port", "p", 53, "Port to bind server too (default 53)")
	flag.StringVarP(&hostFlag, "host", "H", "0.0.0.0", "Host interface to bind too (default 0.0.0.0)")
	flag.BoolVarP(&listRecordsFlag, "list", "l", false, "List all records in the server")
}

func main() {
	flag.Parse()

	// initialize records
	recordsmap, err := record.NewRecordsFromCSV(csvFileFlag)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// List records if flag is true
	if listRecordsFlag {
		for k,v := range recordsmap {
			fmt.Printf("%s => %s\n", k, v)
		}
		os.Exit(0)
	}

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
		serveDNS(u, clientAddr, tcp, recordsmap)
		
	}
}

func serveDNS(u *net.UDPConn, clientAddr net.Addr, request *layers.DNS, recordsmap record.Records) {
	replyMsg := request
	dnsAnswer.Type = layers.DNSTypeA

	var ip string
	var ok bool

	ip, ok = recordsmap[string(request.Questions[0].Name)]
	if !ok {
		// TODO: Log no data present for the request
	}

	// Build DNS Answer
	a, _, _ := net.ParseCIDR(ip + "/24")
	dnsAnswer.Type = layers.DNSTypeA
	dnsAnswer.IP = a
	dnsAnswer.Name = []byte(request.Questions[0].Name)
	dnsAnswer.Class = layers.DNSClassIN

	fmt.Println(request.Questions[0].Name)

	// Build reply message
	replyMsg.QR = true
	replyMsg.ANCount = 1
	replyMsg.OpCode = layers.DNSOpCodeNotify
	replyMsg.AA = true
	replyMsg.Answers = append(replyMsg.Answers, dnsAnswer)
	replyMsg.ResponseCode = layers.DNSResponseCodeNoErr

	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{}
	err := replyMsg.SerializeTo(buf, opts)
	if err != nil {
		// TODO: Handle error
		fmt.Println(err)
	}
	u.WriteTo(buf.Bytes(), clientAddr)
}
