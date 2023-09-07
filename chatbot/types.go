package chatbot

const (
	SFN_CHATBOT_NAME = "chatbot"
	SFN_SINK_NAME    = "sink"
	TAG_REQ          = uint32(0x10)
	TAG_RES          = uint32(0x11)
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Result struct {
	Answer string `json:"answer"`
	Error  string `json:"error"`
}
