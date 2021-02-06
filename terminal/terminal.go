package terminal

import (
	"log"
	"os"
)

func Start() error {
	log.Println("skewer start.")
	return nil
}

func End() {
	log.Println("skewer terminated.")
	os.Exit(0)
}
