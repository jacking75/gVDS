package omokServer

import (
	"encoding/binary"
	"redisDB"
	"smallNet"
)

type Server struct {
	serverNet *smallNet.ServrNet
	onStopNetMsgProcess chan struct{}
	conf OmokConf

	_funcPackeIdlist []int16
	_funclist []func(*gameUser, []byte) int16

	userMgr *gameUserManager

	reqTaskChanRef chan redisDB.ReqTask
	resTaskChan chan redisDB.ResTask
}

func (svr *Server) Init(conf OmokConf, reqTaskChan chan redisDB.ReqTask) {
	svr.conf = conf
	svr.userMgr = newGameUserManager()
	svr.settingPacketFunction()

	svr.reqTaskChanRef = reqTaskChan
	svr.resTaskChan = make(chan redisDB.ResTask, conf.redisResChanCapacity)
}

func (svr *Server) StartServer() {
	packetHeaderSize := int16(5)
	svr.serverNet = smallNet.StartNetwork(confToNetConf(svr.conf), packetHeaderSize, packetTotalSize)

	go svr.Process_goroutine()
}

func (svr *Server) Stop() {
	svr.serverNet.Stop()
}

func confToNetConf(conf OmokConf) smallNet.NetworkConfig {
	netConf := smallNet.NetworkConfig {
		Network: conf.Network,
		BindAddress: conf.BindAddress,
		MaxSessionCount: conf.MaxSessionCount,
		MaxPacketSize: conf.MaxPacketSize,
		RecvPacketRingBufferMaxSize: conf.RecvPacketRingBufferMaxSize,
		SendPacketRingBufferMaxSize: conf.SendPacketRingBufferMaxSize,
		MaxNetMsgChanBufferCount: conf.MaxNetMsgChanBufferCount,
	}

	return netConf
}

func packetTotalSize(data []byte) int16 {
	totalsize := binary.LittleEndian.Uint16(data[:2])
	return int16(totalsize)
}


