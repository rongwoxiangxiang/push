package webservice

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"strconv"
	"time"
	"webs/webservice/common"
)

// 每隔1秒, 检查一次连接是否健康
func (wsConnection *WSConnection) heartbeatChecker() {
	var (
		timer *time.Timer
	)
	timer = time.NewTimer(time.Duration(G_config.WsHeartbeatInterval) * time.Second)
	for {
		select {
		case <-timer.C:
			if !wsConnection.IsAlive() {
				wsConnection.Close()
				goto EXIT
			}
			timer.Reset(time.Duration(G_config.WsHeartbeatInterval) * time.Second)
		case <-wsConnection.closeChan:
			timer.Stop()
			goto EXIT
		}
	}

EXIT:
	// 确保连接被关闭
}

// 处理PING请求
func (wsConnection *WSConnection) handlePing(bizReq *BizMessage) (message *WsMessage, err error) {
	var (
		buf []byte
	)

	wsConnection.KeepAlive()

	if buf, err = json.Marshal(BizPongData{}); err != nil {
		return
	}
	bizResp := &BizMessage{
		Type: MESSAGE_TYPE_PONG,
		Data: json.RawMessage(buf),
	}
	message, err = EncodeWSMessage(bizResp)
	return
}

// 处理JOIN请求
func (wsConnection *WSConnection) handleJoin(bizReq *BizMessage) (message *WsMessage, err error) {
	var (
		bizJoinData *BizJoinData
		existed     bool
	)
	bizJoinData = &BizJoinData{}
	if err = json.Unmarshal(bizReq.Data, bizJoinData); err != nil {
		return
	}
	if len(bizJoinData.Room) == 0 {
		err = common.ERR_ROOM_ID_INVALID
		return
	}
	if len(wsConnection.rooms) >= G_config.MaxJoinRoom {
		log.Println("rooms connection over max num")
		return
	}
	// 已加入过
	if _, existed = wsConnection.rooms[bizJoinData.Room]; existed {
		// 忽略掉这个请求
		return
	}
	// 建立房间 -> 连接的关系
	if err = G_connMgr.JoinRoom(bizJoinData.Room, wsConnection); err != nil {
		return
	}
	// 建立连接 -> 房间的关系
	wsConnection.rooms[bizJoinData.Room] = true
	return
}

// 处理LEAVE请求
func (wsConnection *WSConnection) handleLeave(bizReq *BizMessage) (message *WsMessage, err error) {
	var (
		bizLeaveData *BizLeaveData
		existed      bool
	)
	bizLeaveData = &BizLeaveData{}
	if err = json.Unmarshal(bizReq.Data, bizLeaveData); err != nil {
		return
	}
	if len(bizLeaveData.Room) == 0 {
		err = common.ERR_ROOM_ID_INVALID
		return
	}
	// 未加入过
	if _, existed = wsConnection.rooms[bizLeaveData.Room]; !existed {
		// 忽略掉这个请求
		return
	}
	// 删除房间 -> 连接的关系
	if err = G_connMgr.LeaveRoom(bizLeaveData.Room, wsConnection); err != nil {
		return
	}
	// 删除连接 -> 房间的关系
	delete(wsConnection.rooms, bizLeaveData.Room)
	return
}

func (wsConnection *WSConnection) handleMsg(bizReq *BizMessage) (message *WsMessage, err error) {
	msgData := &BizMsgData{}
	if err = json.Unmarshal(bizReq.Data, msgData); err != nil {
		return
	}
	msgData.FromUser = strconv.FormatUint(wsConnection.connId, 10)
	if msgData.FromUser == msgData.ToUser {
		log.Printf("handleMsg err[0] : same from and to user", msgData.ToUser)
		return
	}

	buf, err := json.Marshal(*msgData)
	if err != nil {
		return
	}
	bizResp := &BizMessage{
		Type: MESSAGE_TYPE_MSG,
		Data: json.RawMessage(buf),
	}
	message, err = EncodeWSMessage(bizResp)
	if err != nil {
		log.Printf("handleMsg err[1] : {}", err)
		return
	}
	to, err := strconv.ParseUint(msgData.ToUser, 10, 64)
	if err != nil {
		log.Printf("handleMsg err[2] : {}", err)
		return
	}
	G_connMgr.PushSingle(to, message)
	return
}

func (wsConnection *WSConnection) leaveAll() {
	var (
		roomId string
	)
	// 从所有房间中退出
	for roomId, _ = range wsConnection.rooms {
		G_connMgr.LeaveRoom(roomId, wsConnection)
		delete(wsConnection.rooms, roomId)
	}
}

// 处理websocket请求
func (wsConnection *WSConnection) WSHandle() {
	var (
		message *WsMessage
		bizReq  *BizMessage
		bizResp *BizMessage
		err     error
	)

	// 连接加入管理器, 可以推送端查找到
	G_connMgr.AddConn(wsConnection)

	log.Printf("WsHandler: new client connect : %v", wsConnection)
	// 心跳检测线程
	go wsConnection.heartbeatChecker()

	// 请求处理协程
	for {
		if message, err = wsConnection.ReadMessage(); err != nil {
			if err == common.ERR_CONNECTION_LOSS {
				log.Printf("WsHandler: ReadMessage close : %v", wsConnection.connId)
			}
			log.Printf("WsHandler: ReadMessage err : %v", err)
			goto ERR
		}

		// 只处理文本消息
		if message.MsgType != websocket.TextMessage {
			continue
		}

		// 解析消息体
		if bizReq, err = DecodeBizMessage(message.MsgData); err != nil {
			log.Printf("WsHandler: DecodeBizMessage err : %v", err)
			goto ERR
		}

		bizResp = nil

		// 1,收到PING则响应PONG: {"type": "PING"}, {"type": "PONG"}
		// 2,收到JOIN则加入ROOM: {"type": "JOIN", "data": {"room": "chrome-plugin"}}
		// 3,收到LEAVE则离开ROOM: {"type": "LEAVE", "data": {"room": "chrome-plugin"}}
		// 4,收到MSG信息返回并发送msg到对应clinet: {"type": "MSG", "data": {"to": "111","msg":"11wasdda"}}

		// 请求串行处理
		switch bizReq.Type {
		case MESSAGE_TYPE_PING:
			if message, err = wsConnection.handlePing(bizReq); err != nil {
				log.Printf("WsHandler: ping err : %v", err)
				goto ERR
			}
		case MESSAGE_TYPE_JOIN:
			if message, err = wsConnection.handleJoin(bizReq); err != nil {
				log.Printf("WsHandler: join err : %v", err)
				goto ERR
			}
			log.Printf("WsHandler: new client [%v] jion room [%s]", wsConnection.connId, bizResp)
		case MESSAGE_TYPE_LEAVE:
			if message, err = wsConnection.handleLeave(bizReq); err != nil {
				log.Printf("WsHandler: leave err : %v", err)
				goto ERR
			}
			log.Printf("WsHandler: one client [%v] leave room [%v]", wsConnection.connId, bizResp)
		case MESSAGE_TYPE_MSG:
			if message, err = wsConnection.handleMsg(bizReq); err != nil {
				log.Printf("WsHandler: send msg err : %v", err)
				goto ERR
			}
			log.Printf("WsHandler: client [%v] send msg [%v]", wsConnection.connId, bizReq.Data)
		default:
			message = nil
			log.Printf("WsHandler: send type err : %v", err)
			goto ERR
		}

		if message != nil {
			if err = wsConnection.SendMessage(message); err != nil {
				if err != common.ERR_SEND_MESSAGE_FULL {
					log.Printf("WsHandler SendMessage err:%v ", err)
					goto ERR
				} else {
					err = nil
				}
			}
		}
	}

ERR:
	// 确保连接关闭
	wsConnection.Close()

	// 离开所有房间
	wsConnection.leaveAll()

	// 从连接池中移除
	G_connMgr.DelConn(wsConnection)
	return
}
