package webservice

import (
	"github.com/gorilla/websocket"
	"sync"
	"time"
	"webs/webservice/common"
)

type WSConnection struct {
	mutex             sync.Mutex
	connId            string
	wsSocket          *websocket.Conn
	inChan            chan *WsMessage
	outChan           chan *WsMessage
	closeChan         chan byte
	isClosed          bool
	lastHeartbeatTime time.Time       // 最近一次心跳时间
	rooms             map[string]bool // 加入了哪些房间,单个ws不存在并发
}

func InitWSConnection(connId string, wsSocket *websocket.Conn) (wsConnection *WSConnection) {
	wsConnection = &WSConnection{
		wsSocket:          wsSocket,
		connId:            connId,
		inChan:            make(chan *WsMessage, G_config.WsInChannelSize),
		outChan:           make(chan *WsMessage, G_config.WsOutChannelSize),
		closeChan:         make(chan byte),
		lastHeartbeatTime: time.Now(),
		rooms:             make(map[string]bool), //同一个用户不能在同一时刻加入两个房间，不存在并发
	}

	go wsConnection.readLoop()
	go wsConnection.writeLoop()

	return
}

// 读websocket
func (wsConnection *WSConnection) readLoop() {
	var (
		msgType int
		msgData []byte
		message *WsMessage
		err     error
	)
	for {
		if msgType, msgData, err = wsConnection.wsSocket.ReadMessage(); err != nil {
			goto ERR
		}
		message = BuildWSMessage(msgType, msgData)

		select {
		case wsConnection.inChan <- message:
		case <-wsConnection.closeChan:
			goto CLOSED
		}
	}

ERR:
	wsConnection.Close()
CLOSED:
}

// 写websocket
func (wsConnection *WSConnection) writeLoop() {
	var (
		message *WsMessage
		err     error
	)
	for {
		select {
		case message = <-wsConnection.outChan:
			if err = wsConnection.wsSocket.WriteMessage(message.MsgType, message.MsgData); err != nil {
				goto ERR
			}
		case <-wsConnection.closeChan:
			goto CLOSED
		}
	}
ERR:
	wsConnection.Close()
CLOSED:
}

// 发送消息
func (wsConnection *WSConnection) SendMessage(message *WsMessage) (err error) {
	select {
	case wsConnection.outChan <- message:
		common.SendMessageTotal_INCR()
	case <-wsConnection.closeChan:
		err = common.ERR_CONNECTION_LOSS
	default:
		err = common.ERR_SEND_MESSAGE_FULL
		common.SendMessageFail_INCR()
	}
	return
}

// 读取消息
func (wsConnection *WSConnection) ReadMessage() (message *WsMessage, err error) {
	select {
	case message = <-wsConnection.inChan:
	case <-wsConnection.closeChan:
		err = common.ERR_CONNECTION_LOSS
	}
	return
}

// 关闭连接
func (wsConnection *WSConnection) Close() {
	wsConnection.wsSocket.Close()

	wsConnection.mutex.Lock()
	defer wsConnection.mutex.Unlock()

	if !wsConnection.isClosed {
		wsConnection.isClosed = true
		close(wsConnection.closeChan)
	}
}

// 检查心跳（不需要太频繁）
func (wsConnection *WSConnection) IsAlive() bool {
	var (
		now = time.Now()
	)

	wsConnection.mutex.Lock()
	defer wsConnection.mutex.Unlock()

	// 连接已关闭 或者 太久没有心跳
	if wsConnection.isClosed || now.Sub(wsConnection.lastHeartbeatTime) > time.Duration(G_config.WsHeartbeatInterval)*time.Second {
		return false
	}
	return true
}

// 更新心跳
func (WSConnection *WSConnection) KeepAlive() {
	var (
		now = time.Now()
	)

	WSConnection.mutex.Lock()
	defer WSConnection.mutex.Unlock()

	WSConnection.lastHeartbeatTime = now
}
