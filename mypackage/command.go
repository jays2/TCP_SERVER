package mypackage

type commandID int

//Constants for every possible commands
const (
	CMD_NICK commandID = iota
	CMD_JOIN
	CMD_DIRECT
	CMD_COPY
	CMD_SHOW
	CMD_KILL
	CMD_DESTROY
)

//Command: Struct for input commands
type Command struct {
	Id     commandID
	Client *Client
	Args   []string
}
