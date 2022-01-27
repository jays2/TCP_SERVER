package mypackage

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

//Server: Struct to create unique instance of server
type Server struct {
	Channels map[string]*Channel
	Commands chan Command
	Nicks    []string
}

//NewServer: Returns a server address
func NewServer() *Server {
	return &Server{
		Channels: make(map[string]*Channel),
		Commands: make(chan Command),
	}
}

//NewClient: Creates a new client instance
func (s *Server) NewClient(conn net.Conn) {
	log.Printf("A new client join the session: %s", conn.RemoteAddr().String())

	c := &Client{
		Conn:     conn,
		Nick:     "Unknown",
		Commands: s.Commands, //Connects with the server channel to transfer data
	}
	c.ReadInput()
}

//Run: Keeps reading chan for incoming commands from clients
func (s *Server) Run() {
	for cmd := range s.Commands {
		switch cmd.Id {
		case CMD_NICK:
			s.Nick(cmd.Client, cmd.Args)
		case CMD_JOIN:
			s.Join(cmd.Client, cmd.Args)
		case CMD_DIRECT:
			s.DirectMessage(cmd.Client, cmd.Args)
		case CMD_COPY:
			s.CopyJunk(cmd.Client, cmd.Args)
		case CMD_SHOW:
			s.Show(cmd.Client)
		case CMD_KILL:
			s.Unsub(cmd.Client, cmd.Args)
		case CMD_DESTROY:
			s.Destroy(cmd.Client)
		}
	}
}

//Nick: Defines a new client nickname
func (s *Server) Nick(c *Client, args []string) {
	nick := args[1]

	if nick == "" {
		c.Msg("Please insert a valid nickname")
		return
	}

	for _, v := range s.Nicks {
		if v == nick {
			c.Msg("This user already exists, cannot create directory")
			return
		}
	}

	//Deletes from server old nickname, in case it already exists
	for i, v := range s.Nicks {
		if v == c.Nick {
			//shift the slice
			s.Nicks = append(s.Nicks[:i], s.Nicks[i+1:]...)
			//Removes client directory
			os.RemoveAll(Current_dir + "/" + c.Nick)
			break
		}
	}

	c.Nick = nick
	c.Msg(fmt.Sprintf("Welcome %s", nick))

	//Creates directory for user
	userDirectory := Current_dir + "/" + nick
	if err := os.Mkdir(userDirectory, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	//Defines a simple text file for future use
	userDirectoryWrite := userDirectory + "/" + nick + ".txt"
	f, err := os.Create(userDirectoryWrite)
	if err != nil {
		log.Fatal(err)
		return
	}

	//Writes text into file
	l, err := f.WriteString("Hi I'm " + nick)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}

	fmt.Println(l, "A file has been saved in server")
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	s.Nicks = append(s.Nicks, nick)
}

//Join: Joins a channel to transfer data
func (s *Server) Join(c *Client, args []string) {
	chanName := args[1]

	if chanName == "" {
		c.Msg("Channel must have a valid name, try again!")
		return
	}

	//Verifies if client is already subscribed
	for _, v := range c.Channels {
		if v.Name == chanName {
			c.Msg("You already belong to this channel")
			return
		}
	}

	//Verifies if channel exists in server, if not we create it
	r, ok := s.Channels[chanName]
	if !ok {
		r = &Channel{
			Name:    chanName,
			Members: make(map[net.Addr]*Client),
		}
		s.Channels[chanName] = r
	}

	//Adds members to channel, if already exists we add a new member
	r.Members[c.Conn.RemoteAddr()] = c

	c.Channels = append(c.Channels, r)

	c.Msg(fmt.Sprintf("Welcome to %s", chanName))

}

//DirectMessage: Sends a chat message to a channel
func (s *Server) DirectMessage(c *Client, args []string) {
	msg := strings.Join(args[2:], " ") //Concatenate and separate by ""
	medium := args[1]
	c.BroadcastDirect(c, medium, c.Nick+": "+msg)
}

//CopyJunk: Copies files from a client to subscribers in defined channel
func (s *Server) CopyJunk(c *Client, args []string) {
	channelCopy := args[1]
	fileCopy := args[2]
	c.BroadcastDirectFiles(c, channelCopy, fileCopy)
}

//Show: Shows client channels
func (s *Server) Show(c *Client) {
	for _, v := range c.Channels {
		c.Msg(v.Name)
	}
}

//Unsub: Unsubscribes from a channel
func (s *Server) Unsub(c *Client, args []string) {
	chanName := args[1]

	if chanName == "" {
		c.Msg("Channel must have a name")
		return
	}

	r, ok := s.Channels[chanName]
	if !ok {
		c.Msg("Channel does not exist")
	}

	//Deletes member on specified channel from server
	delete(r.Members, c.Conn.RemoteAddr())

	//Deletes channel on client
	for i, v := range c.Channels {
		if v.Name == chanName {
			//shift the slice
			c.Channels = append(c.Channels[:i], c.Channels[i+1:]...)
			c.Msg(chanName + " deleted from your subscriptions")
			break
		}
	}

}

//Destroy: Client gets erased for good
func (s *Server) Destroy(c *Client) {
	c.Msg("See you soon!")
	c.Conn.Close()

	//Deletes client from every channel in server
	for _, v := range s.Channels {
		delete(v.Members, c.Conn.RemoteAddr())
	}

	//Deletes nickname in server
	for i, v := range s.Nicks {
		if v == c.Nick {
			//shift the slice
			s.Nicks = append(s.Nicks[:i], s.Nicks[i+1:]...)
			break
		}
	}

	//Removes client directory
	os.RemoveAll(Current_dir + "/" + c.Nick)
}
