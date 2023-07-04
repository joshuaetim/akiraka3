package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/joshuaetim/akiraka3/route"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(route.RunAPI(":4000"))
}
