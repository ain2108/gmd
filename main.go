package main

import (
	"flag"
	"log"
	"net"
	"os"
	"strconv"

	nu "github.com/ain2108/nashutils"
)

func main() {

	ip := "127.0.0.1"

	// Get the flags
	gmkeyp := flag.String("gmkey", "dead", "secret key of the gm address")
	bdkeyp := flag.String("bdkey", "dead", "secret key of the bd address")
	flag.Parse()

	gmkey := *gmkeyp
	bdkey := *bdkeyp

	if gmkey == "dead" || bdkey == "dead" {

		gmkey = os.Getenv("GMKEY")
		bdkey = os.Getenv("BDKEY")

		if gmkey == "" || bdkey == "" {

			flag.Usage()
			log.Fatalf("need to provide both gmkey and bdkey")
		}
	}

	// Initialzie the game dispatcher
	var bd nu.BotDispatcher
	e := bd.Init(ip, nu.BDPort, bdkey, false)
	if e != nil {
		log.Fatal(e)
	}
	bdAddr := ip + ":" + strconv.Itoa(nu.BDPort)
	defer bd.Kill()

	// Init the game master
	var gm nu.GM
	e = gm.Init(ip, nu.GMPort, gmkey, false, bdAddr)
	if e != nil {
		log.Fatal(e)
	}
	defer gm.Kill()
	log.Printf("INFO gmd: gamemaster initialized succesfully on %s:%s\n", ip, strconv.Itoa(nu.GMPort))

	defaultInit(ip)

	// Now we siply sit here indefinately, but we need a better way to terminate
	select {}

}

func defaultInit(ip string) error {

	// Initialize the clerk
	gmAddr := ip + ":" + strconv.Itoa(nu.GMPort)
	var clerk nu.Clerk
	e := clerk.Init(gmAddr)
	if e != nil {
		log.Fatal(e)
	}

	contractAddr := os.Getenv("ADDR1")

	e = clerk.ConnectGame(contractAddr)
	if e != nil {
		log.Fatal(e)
	}

	return nil
}

func checkInput(ip string, port string) {

	trial := net.ParseIP(ip)
	if trial.To4() == nil {
		log.Fatalln("Usage: provided ip address is invalid")
	}

	portnum, err := strconv.Atoi(port)
	if err != nil {
		log.Fatalf("Usage: %v\n", err)
	}

	if portnum < 49151 || portnum > 65535 {
		log.Fatalln("Usage: port has to be in 49151 – 65535 range")
	}
}
