package webservice

import (
	"log"
	"sync"
	"webs/webservice/common"
)

var (
	G_connMgr *WsConnMgr
)

type Map map[string]*WSConnection

type WsConnMgr struct {
	rwMutex     sync.RWMutex
	connections map[string]*WSConnection
	rooms       map[string]*Room
}

func InitConnMgr() (err error) {
	var (
		connMgr *WsConnMgr
	)

	connMgr = &WsConnMgr{
		connections: make(map[string]*WSConnection),
		rooms:       make(map[string]*Room),
	}

	G_connMgr = connMgr
	return
}

func (connMgr *WsConnMgr) AddConn(wsConn *WSConnection) {
	connMgr.rwMutex.Lock()
	defer connMgr.rwMutex.Unlock()

	connMgr.connections[wsConn.connId] = wsConn
	common.OnlineConnections_INCR()
}

func (connMgr *WsConnMgr) DelConn(wsConn *WSConnection) {
	connMgr.rwMutex.Lock()
	defer connMgr.rwMutex.Unlock()

	delete(connMgr.connections, wsConn.connId)
	common.OnlineConnections_DESC()
}

func (connMgr *WsConnMgr) JoinRoom(roomId string, wsConn *WSConnection) (err error) {
	var (
		room    *Room
		existed bool
	)
	connMgr.rwMutex.Lock()
	defer connMgr.rwMutex.Unlock()
	if room, existed = connMgr.rooms[roomId]; !existed {
		room = InitRoom(roomId)
		connMgr.rooms[roomId] = room
		common.RoomCount_INCR()
	}
	err = room.Join(wsConn)
	return
}

func (connMgr *WsConnMgr) LeaveRoom(roomId string, wsConn *WSConnection) (err error) {
	var (
		room    *Room
		existed bool
	)
	connMgr.rwMutex.Lock()
	defer connMgr.rwMutex.Unlock()
	if room, existed = connMgr.rooms[roomId]; !existed {
		err = common.ERR_ROOM_ID_INVALID
		return
	}
	err = room.Leave(wsConn)
	if room.Count() == 0 {
		delete(connMgr.rooms, roomId)
		common.RoomCount_DESC()
	}
	return
}

// 向所有在线用户发送消息
func (connMgr *WsConnMgr) PushSingle(toUser string, msg *WsMessage) {
	var (
		connection *WSConnection
		ok         bool
		err        error
	)
	connMgr.rwMutex.RLock()
	defer connMgr.rwMutex.RUnlock()
	if connection, ok = connMgr.connections[toUser]; ok == false {
		log.Printf("PushSingle user not exist, user {}", toUser)
		return
	}
	err = connection.SendMessage(msg)
	if err != nil {
		log.Printf("PushSingle err: {}", err)
	}
}

// 向所有在线用户发送消息
func (connMgr *WsConnMgr) PushAll(msg *WsMessage) {
	connMgr.rwMutex.RLock()
	defer connMgr.rwMutex.RUnlock()
	for _, connect := range connMgr.connections {
		connect.SendMessage(msg)
	}
}

// 向指定房间发送消息
func (connMgr *WsConnMgr) PushRoom(roomId string, msg *WsMessage) (err error) {
	var (
		room    *Room
		existed bool
	)
	connMgr.rwMutex.RLock()
	defer connMgr.rwMutex.RUnlock()
	if room, existed = connMgr.rooms[roomId]; !existed {
		err = common.ERR_ROOM_ID_INVALID
		return
	}
	room.Push(msg)
	return
}
