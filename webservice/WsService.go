package webservice

import (
	"github.com/gorilla/websocket"
	"log"
	"net"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"
)

// 	WebSocket服务端
type WSServer struct {
	server    *http.Server
	curConnId uint64
}

var (
	G_wsServer *WSServer

	wsUpgrader = websocket.Upgrader{
		// 允许所有CORS跨域请求
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func handleConnect(resp http.ResponseWriter, req *http.Request) {
	var (
		err      error
		wsSocket *websocket.Conn
		connId   uint64
		wsConn   *WSConnection
	)

	// WebSocket握手
	if wsSocket, err = wsUpgrader.Upgrade(resp, req, nil); err != nil {
		log.Printf("WsService upgrade err", err)
		return
	}

	// 连接唯一标识
	connId = atomic.AddUint64(&G_wsServer.curConnId, 1)

	// 初始化WebSocket的读写协程
	wsConn = InitWSConnection(strconv.FormatUint(connId, 10), wsSocket)

	// 开始处理websocket消息
	wsConn.WSHandle()
}

func InitWSServer() (err error) {
	var (
		mux      *http.ServeMux
		server   *http.Server
		listener net.Listener
	)

	// 路由
	mux = http.NewServeMux()
	mux.HandleFunc("/connect", handleConnect)

	// HTTP服务
	server = &http.Server{
		ReadTimeout:  time.Duration(G_config.WsReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(G_config.WsWriteTimeout) * time.Millisecond,
		Handler:      mux,
	}

	// 监听端口
	if listener, err = net.Listen("tcp", ":"+strconv.Itoa(G_config.WsPort)); err != nil {
		log.Fatal("Push Application WsService err :", err)
	}
	log.Println("Push Application start at port :", G_config.WsPort)

	// 赋值全局变量
	G_wsServer = &WSServer{
		server:    server,
		curConnId: uint64(time.Now().Unix()),
	}

	// 拉起服务
	go server.Serve(listener)

	return
}
