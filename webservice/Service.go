package webservice

import (
	"encoding/json"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Service struct {
	server *http.Server
}

var (
	G_service *Service
)

// 全量推送POST msg={}
func handlePushAll(resp http.ResponseWriter, req *http.Request) {
	var (
		err error
		items string
		message *WsMessage
	)
	if err = req.ParseForm(); err != nil {
		return
	}

	items = req.PostForm.Get("items")
	if err = json.Unmarshal([]byte(items), &message); err != nil {
		return
	}
	G_connMgr.PushAll(message)
}

// 房间推送POST room=xxx&msg
func handlePushRoom(resp http.ResponseWriter, req *http.Request) {
	var (
		err error
		room string
		items string
		message *WsMessage
	)
	if err = req.ParseForm(); err != nil {
		return
	}

	room = req.PostForm.Get("room")
	items = req.PostForm.Get("items")

	if err = json.Unmarshal([]byte(items), &message); err != nil {
		return
	}

	G_connMgr.PushRoom(room, message)
}

func InitService() (err error) {
	var (
		mux *http.ServeMux
		server *http.Server
		listener net.Listener
	)

	// 路由
	mux = http.NewServeMux()
	mux.HandleFunc("/push/all", handlePushAll)
	mux.HandleFunc("/push/room", handlePushRoom)

	// TLS证书解析验证
	//if _, err = tls.LoadX509KeyPair(G_config.ServerPem, G_config.ServerKey); err != nil {
	//	return common.ERR_CERT_INVALID
	//}

	server = &http.Server{
		ReadTimeout: time.Duration(G_config.ServiceReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(G_config.ServiceWriteTimeout) * time.Millisecond,
		Handler: mux,
	}

	// 监听端口
	if listener, err = net.Listen("tcp", ":" + strconv.Itoa(G_config.ServicePort)); err != nil {
		return
	}

	// 赋值全局变量
	G_service = &Service{
		server: server,
	}

	// 拉起服务
	go server.Serve(listener)
	//go server.ServeTLS(listener, G_config.ServerPem, G_config.ServerKey)

	for {
		time.Sleep(1 * time.Second)
	}
	os.Exit(0)

	return
}

