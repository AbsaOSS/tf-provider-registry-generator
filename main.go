package main

import (
	"log"

	"github.com/AbsaOSS/tf-provider-registry-generator/internal/config"
	"github.com/AbsaOSS/tf-provider-registry-generator/internal/etl"
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
