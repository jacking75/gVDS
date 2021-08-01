package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"redisDB"
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
		svr.packetProcess(msg.SessionIndex, msg.Data)
	case smallNet.NetMsg_Connect:
		scommon.LogDebug("OnConnect")

	case smallNet.NetMsg_Close:
		scommon.LogDebug("OnClose")

	default:
		scommon.LogDebug("[Process_goroutine] none")
	}
}


func (svr *echoServer) packetProcess(sessionIndex int, packetData []byte) {
	packetID := peekPacketID(packetData)
	_, bodyData := peekPacketBody(packetData)

	if packetID == PACKET_ID_LOGIN_REQ {
		svr.packetProcessLogin(sessionIndex, bodyData)
	} else {
		scommon.LogError(fmt.Sprintf("[packetProcess] invalid packetID: %d", packetID))
	}
}

func (svr *echoServer) packetProcessLogin(sessionIndex int, bodyData []byte) int16 {
	var reqPkt LoginReqPacket
	reqPkt.Decoding(bodyData)

	userID := string(bytes.Trim(reqPkt.UserID[:], "\x00"))

	reqTaskBody := redisDB.ReqTaskLogin {
		UserID: userID,
		AuthCode: reqPkt.AUthCode,
	}
	buf := bytes.NewBuffer(nil)
	_ = gob.NewEncoder(buf).Encode(&reqTaskBody)

	reqTask := redisDB.ReqTask {
		ResChan: svr._resTaskChan,
		UID: uint64(sessionIndex),
		ID: redisDB.TaskID_ReqLogin,
		Data: buf.Bytes(),
	}
	svr._reqTaskChanRef <- reqTask

	return ERROR_CODE_NONE
}

