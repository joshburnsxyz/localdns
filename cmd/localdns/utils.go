package main

import (
	"fmt"
	"net"
	"github.com/google/gopacket"
	layers "github.com/google/gopacket/layers"
	"github.com/joshburnsxyz/localdns/pkg/record"
)

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
