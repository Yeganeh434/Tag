package main

import (
	"flag"
	"log"
	"tag_project/internal/adapters/controllers/http"
	"tag_project/internal/adapters/databases/mysql"
)

var port = flag.Int("port", 8080, "Port to run the HTTP server")

func main() {
	mysql.InitialDatabase()

	flag.Parse()
	err := http.RunWebServer(*port)
	if err != nil {
		log.Printf("could not start server:%v", err)
	}
}
