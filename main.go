package main

import (
	"bytes"
	"context"
	"log"
	"os"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/hkloudou/mqx"
	"github.com/hkloudou/mqx/packets"
	"github.com/hkloudou/mqx/transport"
)

var _connPool = &sync.Map{}

// var _wildcardLists = []string{
// 	"/${uid}/",
// }

// var _topicLists = []string{
// 	"/cfg/platform.json",
// 	"/cfg/platform.json",
// }

// var _sub =
const prefix = "mqx_"

func main() {

	han := &hander{}
	app := mqx.NewMQX()
	app.OnDispose(func(sock transport.Socket) {
		se := sock.Session()
		cid := se.GetString("cid")
		log.Println("sock closed", cid, se.GetString("client_id"))
		_connPool.Delete(cid)
	})
	app.OnConnect(func(sock transport.Socket, req *packets.ConnectPacket) {
		res := packets.NewControlPacket(packets.Connack).(*packets.ConnackPacket)
		res.Qos = req.Qos
		res.ReturnCode = han.Auth(sock, req)
		if res.ReturnCode != packets.Accepted {
			sock.Send(res)
			return
		}

		cid := uuid.New().String()
		sock.Session().Set("cid", cid)
		sock.Session().Set("user_name", req.Username)
		sock.Session().Set("client_id", req.ClientIdentifier)
		_connPool.Store(cid, sock)
		sock.Send(res)
	})
	app.OnSubcribe(func(sock transport.Socket, req *packets.SubscribePacket) {
		res := packets.NewControlPacket(packets.Suback).(*packets.SubackPacket)
		res.MessageID = req.MessageID
		res.ReturnCodes = make([]byte, len(req.Topics))
		sock.Send(res)
	})
	app.OnPublish(func(sock transport.Socket, req *packets.PublishPacket) {
		if req.Qos != 0 {
			return
		}
		if req.Retain {
			rdb := redis.NewClient(&redis.Options{
				Addr:     "localhost:6379",
				Password: "", // no password set
				DB:       10, // use default DB
			})
			defer rdb.Close()
			var buf bytes.Buffer
			if err := req.Write(&buf); err != nil {
				sock.Close()
				return
			}
			if err := rdb.Set(context.TODO(), prefix+req.TopicName, buf.Bytes(), 0).Err(); err != nil {
				sock.Close()
				return
			}
		}
	})
	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
	<-make(chan bool)
}
