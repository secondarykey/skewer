package main

import (
	"flag"
	"log"

	"github.com/secondarykey/skewer"
	"github.com/secondarykey/skewer/config"
)

var (
	port        int
	portEnv     bool
	verbose     bool
	binName     string
	ignoreFiles string
)

func init() {
	flag.IntVar(&port, "p", 8080, `Application Port Number.(give priority to -e.If "0" is specified,port monitoring is not performed)`)
	flag.BoolVar(&portEnv, "e", false, "Get the application port number from an environment variable.")
	flag.BoolVar(&verbose, "v", false, "Show verbose.")

	flag.StringVar(&binName, "n", "skewer-bin", "Name of the file to create.")
	flag.StringVar(&ignoreFiles, "i", ".*", `Specifying files to exclude monitoring(glob pattern,multiple can be specified by "|")P`)
}

func main() {

	flag.Parse()
	args := flag.Args()

	err := skewer.Listen(
		config.SetVerbose(verbose),
		config.SetArgs(args),
		config.SetPort(port, portEnv),
		config.SetBin(binName),
		config.SetIgnoreFiles(ignoreFiles))

	if err != nil {
		if verbose {
			log.Printf("%+v\n", err)
		}
		log.Fatalf("skewer error:\n%s", err)
	}

	log.Println("Terminated...")
}

//Usage
