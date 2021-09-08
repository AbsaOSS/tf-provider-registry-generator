package main

import (
	"log"

	"github.com/k0da/tfreg-golang/internal/config"
	"github.com/k0da/tfreg-golang/internal/etl"
)

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	config, err := config.NewConfig("pages")
	process,err  := etl.NewEtl(config)
	checkError(err)
	err = process.Run()
	checkError(err)
}
