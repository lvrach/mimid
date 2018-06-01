package main

import (
	"log"
	"net/http"
	"os"

	"github.com/lvrach/mimid"
)

func main() {

	args := os.Args[1:]
	switch args[0] {
	case "proxy":
		endpoint := args[1]
		proxy, err := mimid.NewProxy("http://localhost:9922", endpoint, "./testdata")
		if err != nil {
			log.Fatalln(err)
		}

		log.Println("listening on http://localhost:9922")
		log.Println("proxing to ", endpoint)
		log.Println("capturing on", "./testdata")

		log.Fatal(http.ListenAndServe(":9922", proxy))

	case "mock":
		fake := mimid.NewFakeServer("./testdata")

		log.Println("listening on localhost:9922")
		log.Fatal(http.ListenAndServe(":9922", fake))
	}

}
