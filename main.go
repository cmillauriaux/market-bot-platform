package main

import (
	"log"

	"git.icysoft.fr/cedric/kraken-bot/history"
)

func main() {

	_, err := history.LoadHistory(".krakenEUR.csv")
	if err != nil {
		log.Fatal(err)
	}
}
