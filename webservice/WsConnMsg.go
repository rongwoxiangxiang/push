package webservice

import (
	"sync"
	"webs/webservice/common"
)

var (
	G_connMgr *WsConnMgr
)

type WsConnMgr struct {
	rwMutex sync.RWMutex
	connections map[uint64]*WSConnection
	rooms map[string]*Room
}

func InitConnMgr() (err error) {
	var (
		connMgr *WsConnMgr
	)

	connMgr = &WsConnMgr{
		connections: make(map[uint64]*WSConnection),
		rooms: make(map[string]*Room),
	}

	G_connMgr = connMgr
	return
}



func (connMgr *WsConnMgr) AddConn(wsConn *WSConnection) {
	connMgr.rwMutex.Lock()
	defer connMgr.rwMutex.Unlock()
	connMgr.connections[wsConn.connId] = wsConn
}

func (connMgr *WsConnMgr) DelConn(wsConn *WSConnection) {
	connMgr.rwMutex.Lock()
	defer connMgr.rwMutex.Unlock()

	delete(connMgr.connections, wsConn.connId)
}

func (connMgr *WsConnMgr) JoinRoom(roomId string, wsConn *WSConnection) (err error) {
	var (
		room *Room
		existed bool
	)
	if room, existed = connMgr.rooms[roomId]; !existed {
		room = InitRoom(roomId)
		connMgr.rooms[roomId] = room
	}
	err = room.Join(wsConn)
	return
}

func (connMgr *WsConnMgr) LeaveRoom(roomId string, wsConn *WSConnection) (err error) {
	var (
		room *Room
		existed bool
	)
	if room, existed = connMgr.rooms[roomId]; !existed {
		err = common.ERR_ROOM_ID_INVALID
		return
	}
	err = room.Leave(wsConn)
	return
}

// 向所有在线用户发送消息
func (connMgr *WsConnMgr) PushAll(msg *WsMessage) {
	for _, connect := range connMgr.connections {
		connect.SendMessage(msg)
	}
}

// 向指定房间发送消息
func (connMgr *WsConnMgr) PushRoom(roomId string, msg *WsMessage) (err error) {
	var (
		room *Room
		existed bool
	)
	if room, existed = connMgr.rooms[roomId]; !existed {
		err = common.ERR_ROOM_ID_INVALID
		return
	}
	room.Push(msg)
	return
}