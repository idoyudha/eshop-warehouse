package main

import (
	"log"

	"github.com/idoyudha/eshop-warehouse/config"
	"github.com/idoyudha/eshop-warehouse/internal/app"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	app.Run(cfg)
}
