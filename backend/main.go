package main

import (
	_ "github.com/lib/pq"
	"log"
	"random_numbers/internal/server"
)

func main() {
	srv := server.New()
	err := srv.Start(":8000")
	if err != nil {
		log.Fatalln(err.Error())
	}
}
