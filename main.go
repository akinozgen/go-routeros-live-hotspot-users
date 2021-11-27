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
	useTLS = flag.Bool("tls", false, "use tls")

	routerIpPort    = flag.String("router-ip-port", os.Getenv("ROUTER_IP_PORT"), "mikrotik ip")
	routerUser      = flag.String("router-user", os.Getenv("ROUTER_USER"), "mikrotik username")
	routerPwd       = flag.String("router-pwd", os.Getenv("ROUTER_PWD"), "mikrotik password")
	printProperties = flag.String("print-parameters", os.Getenv("PRINT_PARAMETERS"), "Properties")

	interval            = flag.Duration("interval", 1*time.Second, "Interval")
	testMysqlConnection = flag.Bool("test-mysql", false, "Test mysql")

	mysqlUser     = flag.String("mysql-user", os.Getenv("MYSQL_USER"), "mysql username")
	mysqlPwd      = flag.String("mysql-pwd", os.Getenv("MYSQL_PWD"), "mysql password")
	mysqlDb       = flag.String("mysql-db", os.Getenv("MYSQL_DB"), "mysql database")
	mysqlHostPort = flag.String("mysql-host-port", os.Getenv("MYSQL_HOST_PORT"), "mysql host port")
)

var db *sql.DB

func dial() (*routeros.Client, error) {
	if *useTLS {
		return routeros.DialTLS(*routerIpPort, *routerUser, *routerPwd, nil)
	}

	return routeros.Dial(*routerIpPort, *routerUser, *routerPwd)
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

			for _, p := range strings.Split(*printProperties, ",") {
				fmt.Print(re.Map[p], "\t")
			}
			fmt.Print("\n")
		}

		time.Sleep(*interval)
	}
}

func testMysql() {

	config := mysql.Config{
		User:                 *mysqlUser,
		Passwd:               *mysqlPwd,
		Net:                  "tcp",
		Addr:                 *mysqlHostPort,
		DBName:               *mysqlDb,
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
