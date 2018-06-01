package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/lvrach/mimid"
)

//var dataDir = flag.String("data", "./testdata", "specify the directory for saving/loading mock data")

func help() {
	fmt.Println("commands:")
	fmt.Println("	mimid proxy <http_endpoint>		Start a http reverse-proxy that captures request as mock data")
	fmt.Println("	mimid mock				Start a http server that serve mock data")
	fmt.Println("	mimid help				Shows this help message")

	fmt.Println("args:")

	fmt.Println("	-data <dir>				Specify the directory for saving/loading mock data")

}

func main() {
	var dataDir string
	flag.StringVar(&dataDir, "data", "./testdata", "specify the directory for saving/loading mock data")
	flag.Parse()

	switch flag.Arg(0) {
	case "proxy":
		if flag.NArg() < 2 {
			fmt.Println("invalid usage: missing endpoint argument")
			help()
			os.Exit(1)
		}

		endpoint := flag.Arg(1)
		proxy, err := mimid.NewProxy("http://localhost:9922", endpoint, dataDir)
		if err != nil {
			log.Fatalln(err)
		}

		log.Println("listening on http://localhost:9922")
		log.Println("proxing to ", endpoint)
		log.Println("saving data on", dataDir)

		log.Fatal(http.ListenAndServe(":9922", proxy))
	case "mock":
		fake := mimid.NewFakeServer(dataDir)

		log.Println("listening on localhost:9922")
		log.Fatal(http.ListenAndServe(":9922", fake))
	case "help":
		help()
	default:
		fmt.Println("invalid usage")
		help()
		os.Exit(1)
	}
}
