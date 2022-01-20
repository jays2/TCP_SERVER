package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"

	pk "github.com/jays2/general/mypackage"
	"github.com/rs/cors"
)

var s *pk.Server

type VueResponse struct {
	Channel string `json:"channel"`
	Members string `json:"members"`
	Payload int64  `json:"payload"`
}

var vueVar []VueResponse

func main() {
	s = pk.NewServer()

	go handleRequests()

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

func handleRequests() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", returnServer)
	handler := cors.Default().Handler(mux)
	http.ListenAndServe(":3000", handler)
}

func returnServer(w http.ResponseWriter, r *http.Request) {
	vueVar = nil
	var membersContainer string
	for _, v := range s.Channels {
		membersContainer = ""
		for _, m := range v.Members {
			membersContainer += m.Nick + " "
		}
		vueVar = append(vueVar, VueResponse{Channel: v.Name, Members: membersContainer, Payload: v.Payload})
	}
	json.NewEncoder(w).Encode(vueVar)
}
