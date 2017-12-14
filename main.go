package main

import (
	"fmt"
	"os"
	"strings"
	"flag"

	"github.com/hpcloud/tail"
	"github.com/bogdanovich/dns_resolver"
)

var (
	filename string
	dnsServers string
	resolver *dns_resolver.DnsResolver
)

func main() {
	flag.StringVar(&filename, "filename", "/opt/bro/logs/current/dns.log", "")
	flag.StringVar(&dnsServers, "resolvers", "208.67.222.222,208.67.220.220","")
	flag.Parse()

	resolver = dns_resolver.New(strings.Split(dnsServers,","))

	t, err := tail.TailFile(filename, tail.Config{
		Location: &tail.SeekInfo{
			Whence: os.SEEK_END,
		},
		ReOpen: true,
		Follow: true,
	})
	if err != nil {
		fmt.Println("FATAL: ", err)
		os.Exit(1)
	}
	for line := range t.Lines {
		go query(line.Text)
	}
}

func query(line string) {
	// 	0 ts             1511964168.659199
	// 	1 uid            CPEF0D2NbVJOVCthPa
	// 	2 id.orig_h      192.168.151.252
	// 	3 id.orig_p      49833
	// 	4 id.resp_h      8.8.8.8
	// 	5 id.resp_p      53
	// 	6 proto          udp
	// 	7 trans_id       45603
	// 	8 rtt            -
	// 	9 query          elk-relay.skbkontur.ru
	//    10 qclass         1
	//    11 qclass_name    C_INTERNET
	//    12 qtype          28
	//    13 qtype_name     AAAA
	//    14 rcode          0
	//    15 rcode_name     NOERROR
	//    16 AA             F
	//    17 TC             F
	//    18 RD             T
	//    19 RA             F
	//    20 Z              0
	//    21 answers        -
	//    22 TTLs           -
	//    23 rejected       F
	//    24 total_answers  0
	//    25 total_replies  2

	items := strings.Split(line, "\t")
	if len(items) < 25 {
		fmt.Println("Bad line: line")
		return
	}
	if strings.HasSuffix(items[9], "kontur.ru") {
		return
	}
	resolver.LookupHost(items[9])
}
