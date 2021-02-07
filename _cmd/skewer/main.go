package main

import (
	"flag"
	"log"
	"os"

	"github.com/secondarykey/skewer"
	"github.com/secondarykey/skewer/config"
)

var (
	verbose bool
)

func init() {
	flag.BoolVar(&verbose, "v", false, "Verbose Display")
}

func main() {

	flag.Parse()
	args := flag.Args()

	err := skewer.Listen(
		config.SetArgs(args),
		config.SetVerbose(verbose))
	if err != nil {
		log.Printf("Skewer Listen error ------ \n%+v", err)
		os.Exit(1)
	}
	log.Println("Terminated...")
}
