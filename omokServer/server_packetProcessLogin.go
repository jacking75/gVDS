package omokServer

import (
	"bytes"
	"encoding/gob"
	"omokServer/protocol"
	"redisDB"
)

func (svr *Server) packetProcessLogin(user *gameUser, bodyData []byte) int16 {
	var reqPkt protocol.LoginReqPacket
	reqPkt.Decoding(bodyData)

	//TODO 이미 로그인한 유저인지 체크한다
	if user.isEnableLogin() == false {
		return	1
	}

	userID := string(reqPkt.UserID)

	reqTaskBody := redisDB.ReqTaskLogin {
		UserID: userID,
		AuthCode: reqPkt.AUthCode,
	}
	buf := bytes.NewBuffer(nil)
	_ = gob.NewEncoder(buf).Encode(&reqTaskBody)

	reqTask := redisDB.ReqTask {
		ResChan: svr.resTaskChan,
		UID: user.UID,
		ID: redisDB.TaskID_ReqLogin,
		Data: buf.Bytes(),
	}
	svr.reqTaskChanRef <- reqTask


	// 비동기로 답변을 받는다
	resTask := <- svr.resTaskChan

	if resTask.Result != redisDB.TaskResult_Success {

	} else {

	}

	return protocol.ERROR_CODE_NONE
}

func (svr *Server) _sendLoginResponse(sessionIndex int, result int16) {
	res := protocol.LoginResPacket {
		Result: result,
	}

	outBuf := svr.serverNet.GetWBuffer(sessionIndex,32)
	resPkt, resPktSize := res.EncodingPacket(outBuf)
	svr.serverNet.WBufferAheadWCursor(sessionIndex, int(resPktSize))

	svr.serverNet.ISendToClient(sessionIndex, resPkt)
}
