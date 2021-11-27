package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	tm "github.com/buger/goterm"

	"github.com/go-routeros/routeros"
)

var (
	useTLS     = flag.Bool("tls", false, "use tls")
	ip         = flag.String("address", "192.168.1.1:8728", "mikrotik ip")
	user       = flag.String("username", "admin", "mikrotik username")
	pass       = flag.String("password", "admin", "mikrotik password")
	properties = flag.String("properties", "server,user,address,uptime,bytes-in,bytes-out", "Properties")
	interval   = flag.Duration("interval", 1*time.Second, "Interval")
)

func dial() (*routeros.Client, error) {
	if *useTLS {
		return routeros.DialTLS(*ip, *user, *pass, nil)
	}

	return routeros.Dial(*ip, *user, *pass)
}

func main() {
	flag.Parse()

	client, err2 := dial()
	if err2 != nil {
		log.Fatal(err2)
		os.Exit(1)
	}

	for {
		reply, runErr := client.Run("/ip/hotspot/active/print")

		if runErr != nil {
			log.Fatal(runErr)
			os.Exit(1)
		}

		for _, re := range reply.Re {
			tm.Clear()
			for _, p := range strings.Split(*properties, ",") {
				fmt.Print(re.Map[p], "\t")
				tm.Flush()
			}
		}
		fmt.Print("\n")

		time.Sleep(*interval)
	}
}
