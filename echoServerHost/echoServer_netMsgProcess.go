package main

import (
	"scommon"
	"smallNet"
)

func (svr *echoServer) Process_goroutine() {
	scommon.LogInfo("[Process_goroutine] Start")
	defer svr._serverNet.Stop()

LOOP:
	for {
		select {
		case netMsg := <-svr._serverNet.GetNetMsg():
			svr.processNetMsg(netMsg)
		case  <-svr._onStopNetMsgProcess:
			break LOOP
		}
	}

	scommon.LogInfo("[Process_goroutine] End")
}

func (svr *echoServer) processNetMsg(netMsg smallNet.NetMsg) {
	msg := svr._serverNet.PrepareNetMsg(netMsg)

	switch msg.Type {
	case smallNet.NetMsg_Receive:
		//packetID := peekPacketID(msg.Data)
		//bodySize, _ := peekPacketBody(msg.Data)
		//scommon.LogDebug(fmt.Sprintf("[OnReceive] packetID:%d, bodySize:%d", packetID, bodySize))

		incCount()
		svr._serverNet.ISendToClient(msg.SessionIndex, msg.Data)
	case smallNet.NetMsg_Connect:
		scommon.LogDebug("OnConnect")

	case smallNet.NetMsg_Close:
		scommon.LogDebug("OnClose")

	default:
		scommon.LogDebug("[Process_goroutine] none")
	}
}