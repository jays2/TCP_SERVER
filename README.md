# TCP TRANSFER FILE SERVER 

## Introduction
In this project, we're going to deploy a server to: chat, send and receive files through channels that are assigned to clients by subscriptions. Clients are connected via Telnet and are allowed to type different commands to access server faculties. Please be free to try a more secure connection since you wouldn't want someone sniffing your traffic on a production level.

Also use WSL2 for easy compatibility and keep in mind that the directories used in this project must be modified to your file system structure. Directories are created once run, so be careful where you deploy.

We're going to use 3 ports:
- 8888: TCP server, manages client-server connections and commands.
- 8080: Front-end client to visualize activity.
- 3000: API with updated data from the server.

First, let's define the scope of what we can do, and then take a tour to understand what's going on.

To deploy the server at 8888, we use "go run" as usual in "GENERAL" directory:

```
go run main.go
```

Deploy Vue.js client on "front-end" directory (This step is optional, you can interact with the server without visual aid):

```
npm run serve
```

Once the servers are running, a client can connect via Telnet. Many terminals are allowed at a time.
```
telnet localhost 8888
```
Now you're connected as a client and able to use these commands.

- /nick [user]: Creates a nickname to a client connection.
- /join [channel]: Joins a channel or creates it.
- /direct [channel] [message]: Sends a message to all subscribers in channel.
- /copy [channel] [file] : Copies a file from client to all subscribers in channel.
- /show: Displays channel subscriptions from a client.
- /kill [channel]: Unsubscribes a client from a channel. 
- /destroy: Kills client connection.

For every new client a default nickname is created as "Unknown", it has no channels nor directory inside ./Users folder. For this reason, no files can be copied to this user. Only if we assign a nickname its directory is created and full functionality is activated. Needless to say, the files transferred are copied to every user folder on ./Users/[user]. It's worth noting that for every nickname assigned, a [user].txt file it's created inside the client directory, you know... to avoid transferring files from your device to test and such.

Let's take a look for a few examples but not limited to:

## Example #1: Broadcast a message

Client Alice sends a message to Bob:
```
enidsierra@DESKTOP-3BV3HK1:/mnt/c/Users/enids$ telnet localhost 8888
Trying 127.0.0.1...
Connected to localhost.
Escape character is '^]'.
/nick Alice
> Welcome Alice
/join C1
> Welcome to C1
/join C2
> Welcome to C2
/direct C1 Hi!
> Message sent to Bob
```

Client Bob recieves a message from Alice:
```
enidsierra@DESKTOP-3BV3HK1:/mnt/c/Users/enids$ telnet localhost 8888
Trying 127.0.0.1...
Connected to localhost.
Escape character is '^]'.
/nick Bob
> Welcome Bob
/join C1
> Welcome to C1
> Alice: Hi!

```
## Example #2: Send a file

Alice joins channels and sends file to Bob:
```
enidsierra@DESKTOP-3BV3HK1:/mnt/c/Users/enids$ telnet localost 8888
telnet: could not resolve localost/8888: Name or service not known
enidsierra@DESKTOP-3BV3HK1:/mnt/c/Users/enids$ telnet localhost 8888
Trying 127.0.0.1...
Connected to localhost.
Escape character is '^]'.
/nick Alice
> Welcome Alice
/join G1
> Welcome to G1
> A file has been transfered to Bob
```
Bob joins the same channel and recieves file from Alice:
```
enidsierra@DESKTOP-3BV3HK1:/mnt/c/Users/enids$ telnet localost 8888
telnet: could not resolve localost/8888: Name or service not known
enidsierra@DESKTOP-3BV3HK1:/mnt/c/Users/enids$ telnet localhost 8888
Trying 127.0.0.1...
Connected to localhost.
Escape character is '^]'.
/nick Bob
> Welcome Bob
/join G1
> Welcome to G1
> A file has been received from Alice
```

## Example #3: Unsubscribe from a channel

Alice unsubscribes from a channel: 

```
enidsierra@DESKTOP-3BV3HK1:/mnt/c/Users/enids$ telnet localhost 8888
Trying 127.0.0.1...
Connected to localhost.
Escape character is '^]'.
/nick Alice
> Welcome Alice
/join G1
> Welcome to G1
/join G2
> Welcome to G2
/copy G1 Alice.txt
> You sent a file to Bob
/kill G1
> G1 deleted from your subscriptions
```

## Example #4: Destroy client
```
enidsierra@DESKTOP-3BV3HK1:/mnt/c/Users/enids$ telnet localhost 8888
Trying 127.0.0.1...
Connected to localhost.
Escape character is '^]'.
/nick Alice
> Welcome Alice
/join G1
> Welcome to G1
/join G2
> Welcome to G2
/copy G1 Alice.txt
> You sent a file to Bob
/kill G1
> G1 deleted from your subscriptions
/destroy
> See you soon!
Connection closed by foreign host.
enidsierra@DESKTOP-3BV3HK1:/mnt/c/Users/enids$
```

## TCP Server

At the server manager, we have the following structure. Structs connect within each other with pointers, no database is used to store data, everything is store on run time except the files. If you want to delete data from a client, simply kill the connection or close the terminal instance. With every new run, the server will delete any data left behind.

```
//Channel: Struct to handle subscriptions
type Channel struct {
	Name    string
	Members map[net.Addr]*Client
	Payload int64
}
```

```
//Client: Struct for each connected client
type Client struct {
	Conn     net.Conn
	Nick     string
	Commands chan<- Command
	Channels []*Channel
}
```

```
//Command: Struct for input commands
type Command struct {
	Id     commandID
	Client *Client
	Args   []string
}

```

```
//Server: Struct to create a unique server instance
type Server struct {
	Channels map[string]*Channel
	Commands chan Command
	Nicks    []string
}
```


## Client-Server Arquitecture

Client data is kept in a loop, belonging to its instance. Commands are received on this side and sent to the server via channels (Golang channels).

```
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
		cmd := strings.TrimSpace(args[0]) //Root command

		//Sends input to server through chan
		switch cmd {
		case "/nick": //Creates a nickname
			c.Commands <- Command{
				Id:     CMD_NICK,
				Client: c,
				Args:   args,
			}
		case "/join": //Joins a client to a channel
			c.Commands <- Command{
				Id:     CMD_JOIN,
				Client: c,
				Args:   args,
			}
		case "/direct": //Sends a direct message to a channel members
			c.Commands <- Command{
				Id:     CMD_DIRECT,
				Client: c,
				Args:   args,
			}
		case "/copy": //Copies a file to channel members
			c.Commands <- Command{
				Id:     CMD_COPY,
				Client: c,
				Args:   args,
			}
		case "/show": //Displays client channels
			c.Commands <- Command{
				Id:     CMD_SHOW,
				Client: c,
				Args:   args,
			}
		case "/kill": //Unsubscribes a client from a channel
			c.Commands <- Command{
				Id:     CMD_KILL,
				Client: c,
				Args:   args,
			}
		case "/destroy": //Destroys client
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
```

On the other side the server it's listening for incoming commands typed by the clients. All done with the help of concurrency. For each command an action on the server it's performed; acting directly on the client physical address. This affects the overall data, check that the previous loop runs on the same instance; meaning that changes will be reflected globally.

```
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
```

## API feed
API data is collected on main.go with the help of concurrency, data is arranged due to map issues with randomness.

```
func handleRequests() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", returnServer)
	handler := cors.Default().Handler(mux)
	http.ListenAndServe(":3000", handler)
}

func returnServer(w http.ResponseWriter, r *http.Request) {
	vueVar = nil
	var membersContainer string
	var membersSorted []string
	for _, v := range s.Channels {
		membersSorted = nil
		membersContainer = ""
		//Fill slice with members, sort it, and assign it to members container
		for _, m := range v.Members {
			membersSorted = append(membersSorted, m.Nick)
		}
		sort.Strings(membersSorted)
		for _, m := range membersSorted {
			membersContainer += m + ", "
		}
		membersContainer := strings.TrimRight(membersContainer, ", ")
		//vueVar: Slice sent via API to Vue.js
		vueVar = append(vueVar, VueResponse{Channel: v.Name, Members: membersContainer, Payload: v.Payload})
	}

	if vueVar != nil {
		sort.Slice(vueVar, func(i, j int) bool {
			return vueVar[i].Channel < vueVar[j].Channel
		})
		json.NewEncoder(w).Encode(vueVar)
	}
}

```

## Final thoughts and improvements

As mentioned, security it's a big issue here. Encrypted connections client-server would be ideal, as well as better parsing at the terminal level, to avoid server collapse. Superusers are crucial to managing connections, this repo does the job, but it can improve at a production level. Even though; the structure is set to keep scaling and correct these issues. 























