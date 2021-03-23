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
	mode        string
	args        string
)

func init() {
	flag.IntVar(&port, "p", 8080, `Application Port Number.(give priority to -e.If "0" is specified,port monitoring is not performed)`)
	flag.BoolVar(&portEnv, "e", false, "Get the application port number from an environment variable.")
	flag.BoolVar(&verbose, "v", false, "Show verbose.")

	flag.StringVar(&binName, "n", "skewer-bin", "Name of the file to create.")
	flag.StringVar(&ignoreFiles, "i", ".*", `Specifying files to exclude monitoring(glob pattern,multiple can be specified by "|")P`)
	flag.StringVar(&mode, "m", "http", `Plan to implement test mode etc...`)
	flag.StringVar(&args, "args", "", `go run arguments`)
}

func main() {

	flag.Parse()
	files := flag.Args()

	err := skewer.Patrol(
		config.SetVerbose(verbose),
		config.SetFiles(files),
		config.SetMode(mode, port, portEnv),
		config.SetBin(binName),
		config.SetArgs(args),
		config.SetIgnoreFiles(ignoreFiles))

	if err != nil {
		if verbose {
			log.Printf("%+v\n", err)
		}
		log.Fatalf("skewer error:\n%s", err)
	}

	log.Println("skewer terminated.")
}

//Usage
