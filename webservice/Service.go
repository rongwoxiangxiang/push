package webservice

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
	"webs/webservice/common"
)

type Service struct {
	server *http.Server
}

var (
	G_service *Service
)

func handleStats(resp http.ResponseWriter, req *http.Request) {
	var (
		data []byte
		err error
	)

	if data, err = common.G_stats.Dump(); err != nil {
		return
	}

	resp.Write(data)
}

/**
 * items:[{"msg": "hi"},{"msg": "bye"}]
 */
func handlePushAll(resp http.ResponseWriter, req *http.Request) {
	var (
		err error
		items string
		msgIdx int
		msgArr []json.RawMessage
		bizMessage *BizMessage
		message *WsMessage
	)
	if err = req.ParseForm(); err != nil {
		log.Printf("Service: push all err[1]:  %v", err)
		return
	}

	items = req.PostForm.Get("items")
	if err = json.Unmarshal([]byte(items), &msgArr); err != nil {
		log.Printf("Service: push all err[2]:  %v", err)
		return
	}
	for msgIdx, _  = range msgArr {
		bizMessage = &BizMessage{
			Type: "PUSH",
			Data: json.RawMessage(msgArr[msgIdx]),
		}
		message, err = EncodeWSMessage(bizMessage)
		if err != nil {
			log.Printf("Service: push all err[3]:  %v", err)
			return
		}
		G_connMgr.PushAll(message)
	}
}

// 房间推送POST room=xxx&msg
func handlePushRoom(resp http.ResponseWriter, req *http.Request) {
	var (
		err error
		items string
		room string
		msgIdx int
		msgArr []json.RawMessage
		bizMessage *BizMessage
		message *WsMessage
	)
	if err = req.ParseForm(); err != nil {
		log.Printf("Service: push room err[1]:  %v", err)
		return
	}

	room = req.PostForm.Get("room")
	items = req.PostForm.Get("items")

	if err = json.Unmarshal([]byte(items), &message); err != nil {
		log.Printf("Service: push room err[2]:  %v", err)
		return
	}

	for msgIdx, _  = range msgArr {
		bizMessage = &BizMessage{
			Type: "PUSH",
			Data: json.RawMessage(msgArr[msgIdx]),
		}
		message, err = EncodeWSMessage(bizMessage)
		if err != nil {
			log.Printf("Service: push room err[3]:  %v", err)
			return
		}
		G_connMgr.PushRoom(room, message)
	}
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
	mux.HandleFunc("/stats", handleStats)


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

