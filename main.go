package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ankeshnirala/go/go_services/api"
	"github.com/ankeshnirala/go/go_services/storage"
	"github.com/joho/godotenv"
)

func main() {
	logger := log.New(os.Stdout, "prractice_backend ", log.LstdFlags)

	err := godotenv.Load("app.env")
	if err != nil {
		logger.Println(err)
	}

	appPort := os.Getenv("PORT")

	listenAddr := flag.String("listenaddr", appPort, "the server address")
	flag.Parse()

	mongoStore, err := storage.NewMongoStore()
	if err != nil {
		logger.Println("MongoDB Connection Error: ", err.Error())
		return
	}
	logger.Println("MongoDB Connected Successfully")

	server := api.NewServer(logger, *listenAddr, *mongoStore)
	msg := fmt.Sprintf("started server on [::]:%s, url: http://localhost:%s", *listenAddr, *listenAddr)
	logger.Println(msg)

	logger.Fatal(server.Start())
}
