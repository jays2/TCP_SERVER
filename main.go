package main

import (
	"log"
	"net"
	"os"

	pk "github.com/jays2/general/mypackage"
)

func main() {
	s := pk.NewServer()
	go s.Run()

	//Deletes users directory if existent
	os.RemoveAll(pk.Current_dir)

	//Creates users directory for storage
	if err := os.Mkdir(pk.Current_dir, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	//Listens for incoming conections
	listener, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatalf("Sever was not able to start: %s", err.Error())
	}

	defer listener.Close()

	log.Printf("Server deployed on port :8888")

	//Accepts incoming connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Communication has failed: %s", err.Error())
			continue
		}

		go s.NewClient(conn)
	}
}
