package omokServer

import (
	"scommon"
	"smallNet"
)

func (svr *Server) Process_goroutine() {
	scommon.LogInfo("[Process_goroutine] Start")

	defer svr.serverNet.Stop()

LOOP:
	for {
		select {
		case netMsg := <-svr.serverNet.GetNetMsg():
			svr.processNetMsg(netMsg)
		case  <-svr.onStopNetMsgProcess:
			break LOOP
		}
	}

	scommon.LogInfo("[Process_goroutine] End")
}

func (svr *Server) processNetMsg(netMsg smallNet.NetMsg) {
	msg := svr.serverNet.PrepareNetMsg(netMsg)

	switch msg.Type {
	case smallNet.NetMsg_Receive:
		scommon.LogDebug("OnReceive")
		svr.packetProcess(msg.SessionIndex, msg.Data)
	case smallNet.NetMsg_Connect:
		scommon.LogDebug("OnConnect")
		svr.userMgr.addUser(msg.SessionIndex)
	case smallNet.NetMsg_Close:
		scommon.LogDebug("OnClose")
		svr.userMgr.removeUser(msg.SessionIndex)
	default:
		scommon.LogDebug("[Process_goroutine] none")
	}
}




