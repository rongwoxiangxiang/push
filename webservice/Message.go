package webservice

import (
	"encoding/json"
	"github.com/gorilla/websocket"
)

// 推送类型
var (
	PUSH_TYPE_ROOM = 1 // 推送房间
	PUSH_TYPE_ALL  = 2 // 推送在线
)

var (
	MESSAGE_TYPE_PING  = "PING"
	MESSAGE_TYPE_PONG  = "PONG"
	MESSAGE_TYPE_JOIN  = "JOIN"
	MESSAGE_TYPE_LEAVE = "LEAVE"
	MESSAGE_TYPE_PUSH  = "PUSH"
	MESSAGE_TYPE_MSG   = "MSG"
)

type WsMessage struct {
	MsgType int
	MsgData []byte
}

// 业务消息的固定格式(type+data)
type BizMessage struct {
	Type string          `json:"type"` // type消息类型: PING, PONG, JOIN, LEAVE, PUSH, MSG
	Data json.RawMessage `json:"data"` // data数据字段
}

// Data数据类型

// PUSH
type BizPushData struct {
	Items []*json.RawMessage `json:"items"`
}

// PING
type BizPingData struct{}

// PONG
type BizPongData struct{}

// JOIN
type BizJoinData struct {
	Room string `json:"room"`
}

// LEAVE
type BizLeaveData struct {
	Room string `json:"room"`
}

// Msg
type BizMsgData struct {
	Msg      string `json:"msg"`
	ToUser   string `json:"to"`   // 接受者
	FromUser string `json:"from"` // 发送者
}

func BuildWSMessage(msgType int, msgData []byte) (wsMessage *WsMessage) {
	return &WsMessage{
		MsgType: msgType,
		MsgData: msgData,
	}
}

func EncodeWSMessage(bizMessage *BizMessage) (wsMessage *WsMessage, err error) {
	var (
		buf []byte
	)
	if MESSAGE_TYPE_MSG == bizMessage.Type {
		wsMessage = &WsMessage{websocket.TextMessage, bizMessage.Data}
		return
	}
	if buf, err = json.Marshal(*bizMessage); err != nil {
		return
	}
	wsMessage = &WsMessage{websocket.TextMessage, buf}
	return
}

// 解析{"type": "PING", "data": {...}}的包
func DecodeBizMessage(buf []byte) (bizMessage *BizMessage, err error) {
	var (
		bizMsgObj BizMessage
	)

	if err = json.Unmarshal(buf, &bizMsgObj); err != nil {
		return
	}

	bizMessage = &bizMsgObj
	return
}
