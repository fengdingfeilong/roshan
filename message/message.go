package message

//CmdType is the Type of Command
type CmdType int32

//Message is the payload of the packet
type Message struct {
	Name    string `json:"-"`
	ID      string `json:"id"`
	Version string `json:"version"`
}

const MessageVersion string = "1.0.0.0"

func (message Message) String() string {
	return message.Name
}
