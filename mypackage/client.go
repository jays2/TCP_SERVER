package mypackage

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

//Client: Struct for each connected client
type Client struct {
	Conn     net.Conn
	Nick     string
	Commands chan<- Command
	Channels []*Channel
}

//ReadInput: Reads client inputs
func (c *Client) ReadInput() {
	//Loop mantains updated client input data
	for {
		msg, err := bufio.NewReader(c.Conn).ReadString('\n')
		if err != nil {
			fmt.Println("User " + c.Nick + " has left")
			msg = "/destroy"
		}

		msg = strings.Trim(msg, "\r\n")
		args := strings.Split(msg, " ")
		cmd := strings.TrimSpace(args[0])

		//Sends input to server through chan
		switch cmd {
		case "/nick":
			c.Commands <- Command{
				Id:     CMD_NICK,
				Client: c,
				Args:   args,
			}
		case "/join":
			c.Commands <- Command{
				Id:     CMD_JOIN,
				Client: c,
				Args:   args,
			}
		case "/direct":
			c.Commands <- Command{
				Id:     CMD_DIRECT,
				Client: c,
				Args:   args,
			}
		case "/copy":
			c.Commands <- Command{
				Id:     CMD_COPY,
				Client: c,
				Args:   args,
			}
		case "/show":
			c.Commands <- Command{
				Id:     CMD_SHOW,
				Client: c,
				Args:   args,
			}
		case "/kill":
			c.Commands <- Command{
				Id:     CMD_KILL,
				Client: c,
				Args:   args,
			}
		case "/destroy":
			c.Commands <- Command{
				Id:     CMD_DESTROY,
				Client: c,
				Args:   nil,
			}
			return
		default:
			c.Msg("Unknown command")
		}
	}
}

//Msg: Displays a message to client
func (c *Client) Msg(msg string) {
	c.Conn.Write([]byte("> " + msg + "\n"))
}

//BroadcastDirect: Sends a message to everybody except the sender
func (r *Client) BroadcastDirect(sender *Client, medium string, msg string) {
	if medium == "" || msg == "" {
		sender.Msg("Command incomplete")
	}

	var n = 0
	for _, can := range sender.Channels {
		if can.Name == medium {
			n = n + 1
			can.Broadcast(sender, msg)
		}
	}

	if n == 0 {
		sender.Msg("Not logged?, please join the channel")
	}
}

//BroadcastDirectFiles: Sends a file to everybody except the sender
func (r *Client) BroadcastDirectFiles(sender *Client, channel string, file string) {
	var n = 0

	for _, can := range sender.Channels {
		if can.Name == channel {
			n = n + 1
			can.BroadcastFiles(sender, file)
		}
	}

	if n == 0 {
		sender.Msg("Not logged?, please join the channel")
	}
}
