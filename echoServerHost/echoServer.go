package main

import (
	"encoding/binary"
	"fmt"
	"smallNet"
)

type echoServer struct {
	_serverNet           *smallNet.ServrNet
	_onStopNetMsgProcess chan struct{}

	_port int
	conf  echoConf
}

type echoConf struct {
	// 네트워크쪽
	Network              string // tcp4(ipv4 only), tcp6(ipv6 only), tcp(ipv4, ipv6)
	IPAddress          string // 127.0.0.1
	MaxSessionCount      int    // 최대 클라이언트 세션 수. 넉넉하게 많이 해도 괜찮다
	MaxPacketSize        int    // 최대 패킷 크기
	RecvPacketRingBufferMaxSize int
	MaxNetMsgChanBufferCount     int // 네트워크 이벤트 메시지 채널 버퍼의 최대 크기
}


func (svr *echoServer) Init(port int, conf echoConf) {
	svr.conf = conf
	svr._port = port
}

func (svr *echoServer) StartServer() {
	packetHeaderSize := int16(5)
	svr._serverNet = smallNet.StartNetwork(confToNetConf(svr._port, svr.conf), packetHeaderSize, packetTotalSize)

	go svr.Process_goroutine()
}

func (svr *echoServer) Stop() {
	svr._serverNet.Stop()
}

func confToNetConf(port int, conf echoConf) smallNet.NetworkConfig {
	bindAddress := fmt.Sprintf("%s:%d", conf.IPAddress, port)

	netConf := smallNet.NetworkConfig {
		Network: conf.Network,
		BindAddress: bindAddress,
		MaxSessionCount: conf.MaxSessionCount,
		MaxPacketSize: conf.MaxPacketSize,
		RecvPacketRingBufferMaxSize: conf.RecvPacketRingBufferMaxSize,
		MaxNetMsgChanBufferCount: conf.MaxNetMsgChanBufferCount,
	}

	return netConf
}

func packetTotalSize(data []byte) int16 {
	totalsize := binary.LittleEndian.Uint16(data[:2])
	return int16(totalsize)
}
