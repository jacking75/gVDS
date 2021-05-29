package omokServer

import (
	"scommon"
	"smallNet"
	"time"
)

func (svr *Server) Process_goroutine() {
	scommon.LogInfo("[Process_goroutine] Start")

	timeTicker := time.NewTicker(time.Millisecond * 100)

	defer timeTicker.Stop()
	defer svr._serverNet.Stop()

LOOP:
	for {
		select {
		case _ = <-timeTicker.C:
			svr._userMgr.checkUserState()
		case netMsg := <-svr._serverNet.GetNetMsg():
			svr.processNetMsg(netMsg)
		case  <-svr._onStopNetMsgProcess:
			break LOOP
		}
	}

	scommon.LogInfo("[Process_goroutine] End")
}

func (svr *Server) processNetMsg(netMsg smallNet.NetMsg) {
	msg := svr._serverNet.PrepareNetMsg(netMsg)

	switch msg.Type {
	case smallNet.NetMsg_Receive:
		scommon.LogDebug("OnReceive")

		svr.packetProcess(msg.SessionIndex, msg.Data)
	case smallNet.NetMsg_Connect:
		scommon.LogDebug("OnConnect")

		conf := userConf{
			heartbeatReqIntervalTimeMSec: int64(svr.conf.HeartbeatReqIntervalTimeMSec),
			heartbeatWaitTimeMSec: int64(svr.conf.HeartbeatWaitTimeMSec),
		}
		svr._userMgr.addUser(msg.SessionIndex, conf)
	case smallNet.NetMsg_Close:
		scommon.LogDebug("OnClose")

		svr._userMgr.removeUser(msg.SessionIndex)
	default:
		scommon.LogDebug("[Process_goroutine] none")
	}
}




