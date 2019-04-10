package ws

type Message struct {
	Sender    string `json:"sender,omitempty"`
	Recipient string `json:"to,omitempty"`
	Content   string `json:"content,omitempty"`
	Type	  string `json:"type,omitempty"`
}

const (
	MESSAGE_TYPE_TO_SINGLE_USER = "single"	//单发
	MESSAGE_TYPE_TO_GROUP_USER  = "group"	//组发
	MESSAGE_TYPE_TO_MANY_USER   = "many" //群发
	MESSAGE_TYPE_TO_BROADCAST   = "broadcast" //广播
)