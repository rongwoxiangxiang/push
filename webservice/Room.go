package webservice

import (
	"sync"
	"webs/webservice/common"
)

type Room struct {
	rwMutex     sync.RWMutex
	roomId      string
	connections map[uint64]*WSConnection
}

func InitRoom(roomId string) (room *Room) {
	room = &Room{
		roomId:      roomId,
		connections: make(map[uint64]*WSConnection),
	}
	return
}

func (room *Room) Join(wsConn *WSConnection) (err error) {
	var (
		existed bool
	)

	room.rwMutex.Lock()
	defer room.rwMutex.Unlock()

	if _, existed = room.connections[wsConn.connId]; existed {
		err = common.ERR_JOIN_ROOM_TWICE
		return
	}

	room.connections[wsConn.connId] = wsConn
	return
}

func (room *Room) Leave(wsConn *WSConnection) (err error) {
	var (
		existed bool
	)

	room.rwMutex.Lock()
	defer room.rwMutex.Unlock()

	if _, existed = room.connections[wsConn.connId]; !existed {
		err = common.ERR_NOT_IN_ROOM
		return
	}

	delete(room.connections, wsConn.connId)
	return
}

func (room *Room) Count() int {
	room.rwMutex.RLock()
	defer room.rwMutex.RUnlock()

	return len(room.connections)
}

func (room *Room) Push(wsMsg *WsMessage) {
	var (
		wsConn *WSConnection
	)
	room.rwMutex.RLock()
	defer room.rwMutex.RUnlock()

	for _, wsConn = range room.connections {
		wsConn.SendMessage(wsMsg)
	}
}
