package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/joshuaetim/frontdesk/route"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(route.RunAPI(":5000"))
}
