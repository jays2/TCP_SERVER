package mypackage

import (
	"io"
	"log"
	"net"
	"os"
)

//Channel: Struct to handle subscriptions
type Channel struct {
	Name    string
	Members map[net.Addr]*Client
	Payload int64
}

//Broadcast: Handles broadcast messages to subscribers on channel
func (r *Channel) Broadcast(sender *Client, msg string) {
	var n = 0
	for addr, m := range r.Members {
		if sender.Conn.RemoteAddr() != addr {
			n = n + 1
			m.Msg(msg)
			sender.Msg("Message sent to " + m.Nick)
		}
	}

	if n == 0 {
		sender.Msg("Not able to deliver any message")
	}

}

//BroadcastFiles: Handles file transfer from client directory to subscribers on channel
func (r *Channel) BroadcastFiles(sender *Client, msg string) {
	var n bool
	openFileRoute := Current_files + sender.Nick + "/" + msg
	for addr, reciever := range r.Members {
		if sender.Conn.RemoteAddr() != addr {
			n = true
			createFileRoute := Current_files + reciever.Nick + "/" + msg
			CopyFiles(openFileRoute, createFileRoute, reciever, sender, r)

		}
	}

	if !n {
		sender.Msg("Not a single client to send file")
	}
}

//CopyFiles: Open, creates and copies bytes on destination channel
func CopyFiles(origin string, destination string, reciever *Client, sender *Client, r *Channel) {
	sourceFile, err := os.Open(origin)
	if err != nil {
		sender.Msg("Cannot open file")
		return
	}
	defer sourceFile.Close()

	newFile, err := os.Create(destination)
	if err != nil {
		sender.Msg("Cannot create file, check directory")
		return
	}
	defer newFile.Close()

	bytesCopied, err := io.Copy(newFile, sourceFile)
	if err != nil {
		sender.Msg("Cannot copy bytes")
		return
	}

	log.Printf("Bytes %d copied", bytesCopied)
	r.Payload += bytesCopied
	reciever.Msg("A file has been received from " + sender.Nick)
	sender.Msg("You sent a file to " + reciever.Nick)
}
