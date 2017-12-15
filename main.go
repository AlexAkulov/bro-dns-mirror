package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

var (
	ds                             string
	offset                         string
	countOK, countErr, countIgnore int64
	dnsServers                     []string
)

func main() {
	dnsServers = strings.Split(ds, ",")

	stdinReader := bufio.NewReader(os.Stdin)
	var line bytes.Buffer
	for {
		part, isPrefix, err := stdinReader.ReadLine()
		if err == nil {
			line.Write(part)
			if !isPrefix {
				l := make([]byte, line.Len())
				copy(l, line.Bytes())
				parts := strings.Split(string(l), "\t")
				if len(parts) == 2 {
					go resolve(parts[0], parts[1])
				}
				line.Reset()
			}
			continue
		}
		if err != io.EOF {
			fmt.Printf("can't read: %v\n", err)
			os.Exit(1)
		}
		break
	}
}

func resolve(query, qtype string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx,
		"/usr/bin/nslookup",
		"-timeout=4",
		fmt.Sprintf("-type=%s", qtype),
		query,
		dnsServers[0])
	out, err := cmd.CombinedOutput()
	oneline := strings.Split(string(out), "\n")
	if err != nil {
		fmt.Printf("FAIL: type: %s\tquery: %s\toutput: %s\n", qtype, query, oneline[0])
		return
	}
	fmt.Printf("OK  : type: %s\tquery: %s\n", qtype, query)
}
