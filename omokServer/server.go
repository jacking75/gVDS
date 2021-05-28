package omokServer

import (
	"encoding/binary"
	"fmt"
	"redisDB"
	"smallNet"
)

//TODO redis를 사용한 로그인
//TODO 들어옴 통보하기. 로그인 하면
//TODO 나감 통보하기. 로그인 유저가 접속이 끊긴 경우
//TODO 오목 게임 플레이

type Server struct {
	serverNet *smallNet.ServrNet
	onStopNetMsgProcess chan struct{}

	port int
	conf OmokConf

	_funcPackeIdlist []int16
	_funclist []func(*gameUser, []byte) int16

	userMgr *gameUserManager

	reqTaskChanRef chan redisDB.ReqTask
	resTaskChan chan redisDB.ResTask
}

func (svr *Server) Init(port int, conf OmokConf, reqTaskChan chan redisDB.ReqTask) {
	svr.conf = conf
	svr.port = port
	svr.userMgr = newGameUserManager()
	svr.settingPacketFunction()

	svr.reqTaskChanRef = reqTaskChan
	svr.resTaskChan = make(chan redisDB.ResTask, conf.RedisResChanCapacity)
}

func (svr *Server) StartServer() {
	packetHeaderSize := int16(5)
	svr.serverNet = smallNet.StartNetwork(confToNetConf(svr.port, svr.conf), packetHeaderSize, packetTotalSize)

	go svr.Process_goroutine()
}

func (svr *Server) Stop() {
	svr.serverNet.Stop()
}

func confToNetConf(port int, conf OmokConf) smallNet.NetworkConfig {
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


