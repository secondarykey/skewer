package main

import (
	"log"
	"os"

	"github.com/secondarykey/skewer"
)

func main() {
	err := skewer.Listen()
	if err != nil {
		log.Printf("Skewer Listen error ------ \n%+v", err)
		os.Exit(1)
	}
}
