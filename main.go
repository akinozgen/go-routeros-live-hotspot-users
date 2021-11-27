package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	mysql "github.com/go-sql-driver/mysql"

	"github.com/go-routeros/routeros"
)

var (
	useTLS              = flag.Bool("tls", false, "use tls")
	ip                  = flag.String("address", os.Getenv("ROUTER_IP_PORT"), "mikrotik ip")
	user                = flag.String("username", os.Getenv("ROUTER_USER"), "mikrotik username")
	pass                = flag.String("password", os.Getenv("ROUTER_PWD"), "mikrotik password")
	properties          = flag.String("properties", os.Getenv("PRINT_PARAMETERS"), "Properties")
	interval            = flag.Duration("interval", 1*time.Second, "Interval")
	testMysqlConnection = flag.Bool("test-mysql", false, "Test mysql")
)

var db *sql.DB

func dial() (*routeros.Client, error) {
	if *useTLS {
		return routeros.DialTLS(*ip, *user, *pass, nil)
	}

	return routeros.Dial(*ip, *user, *pass)
}

func main() {

	flag.Parse()
	// Testing for mysql connection
	if *testMysqlConnection {
		testMysql()
	}

	client, err := dial()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// Print hotspot active users forevah
	for {
		reply, runErr := client.Run("/ip/hotspot/active/print")

		if runErr != nil {
			log.Fatal(runErr)
			os.Exit(1)
		}

		for _, re := range reply.Re {
			clear()

			for _, p := range strings.Split(*properties, ",") {
				fmt.Print(re.Map[p], "\t")
			}
			fmt.Print("\n")
		}

		time.Sleep(*interval)
	}
}

func testMysql() {

	config := mysql.Config{
		User:                 os.Getenv("MYSQL_USER"),
		Passwd:               os.Getenv("MYSQL_PWD"),
		Net:                  "tcp",
		Addr:                 os.Getenv("MYSQL_HOST_PORT"),
		DBName:               os.Getenv("MYSQL_DB"),
		AllowNativePasswords: true,
	}

	db, dbErr := sql.Open("mysql", config.FormatDSN())

	if dbErr != nil {
		log.Fatal(dbErr)
		os.Exit(1)
	}

	fmt.Println(db.Ping())

	os.Exit(0)
}

func clear() {
	cmdName := "clear"
	cmd := exec.Command(cmdName)
	cmd.Stdout = os.Stdout
	cmd.Run()
}
