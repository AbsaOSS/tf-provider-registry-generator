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
	checkError(err)
	process, err := etl.NewEtl(etl.NewEtlFactory(config))
	checkError(err)
	err = process.Run()
	checkError(err)
}
