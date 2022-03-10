package main

import (
	"log"

	"github.com/hkloudou/mqx/packets"
	"github.com/hkloudou/mqx/transport"
)

type hander struct {
}

func (m *hander) Auth(sock transport.Socket, req *packets.ConnectPacket) byte {
	if req.ProtocolVersion != 4 {
		return packets.ErrRefusedBadProtocolVersion
	}
	state := sock.ConnectState()
	if state != nil {
		// log.Println("state.PeerCertificates", state.PeerCertificates)
		for index, v := range state.PeerCertificates {
			// if !v.IsCA {
			log.Println(index, v.Subject.CommonName, v.IsCA)
			// }
		}
	}
	// c, exist := sock.Session().Get("s.peercert.0")
	// log.Println(c, exist)
	// if !exist {
	// 	return packets.ErrRefusedNotAuthorised
	// }
	return packets.Accepted
}
