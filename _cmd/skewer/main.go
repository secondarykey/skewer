package main

import (
	"flag"
	"log"
	"os"

	"github.com/secondarykey/skewer"
	"github.com/secondarykey/skewer/config"
)

func main() {

	flag.Parse()
	args := flag.Args()

	err := skewer.Listen(config.SetArgs(args))
	if err != nil {
		log.Printf("Skewer Listen error ------ \n%+v", err)
		os.Exit(1)
	}
	log.Println("Terminated...")
}
