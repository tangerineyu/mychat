package websocket

import "encoding/json"

// 消息类型
type Action string

const (
	ActionHeartbeat   Action = "heartbeat"    //心跳
	ActionLogin       Action = "login"        //登录
	ActionChatMessage Action = "chat_message" //聊天消息
	ActionReCall      Action = "recall"       //撤回
	ActionAck         Action = "ack"
)

type Message struct {
	Action Action `json:"action"`
	//不使用string，因为内容会被转义
	Content json.RawMessage `json:"content"`
	//追踪日志
	TraceId string `json:"trace_id,omitempty"`
}
type ChatMessageContent struct {
	SendId     string `json:"send_id"`     //发送者
	ReceiverId string `json:"receiver_id"` //接收者
	Type       int    `json:"type"`        //1:文本， 2：图片
	Content    string `json:"content"`     //文本内容 or 图片内容
	Uuid       string `json:"uuid"`        //ACK
}
type AckMessage struct {
	MsgId  string `json:"msg_id"`
	UserId string `json:"user_id"`
}
