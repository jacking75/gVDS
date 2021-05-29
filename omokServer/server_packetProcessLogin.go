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

	userID := string(bytes.Trim(reqPkt.UserID[:], "\x00"))

	reqTaskBody := redisDB.ReqTaskLogin {
		UserID: userID,
		AuthCode: reqPkt.AUthCode,
	}
	buf := bytes.NewBuffer(nil)
	_ = gob.NewEncoder(buf).Encode(&reqTaskBody)

	reqTask := redisDB.ReqTask {
		ResChan: svr._resTaskChan,
		UID: user._uid,
		ID: redisDB.TaskID_ReqLogin,
		Data: buf.Bytes(),
	}
	svr._reqTaskChanRef <- reqTask


	// 동기로 답변을 받는다
	resTask := <- svr._resTaskChan

	if resTask.Result == redisDB.TaskResult_Success {
		user.setAuth(userID)
	}

	_sendLoginResponse(svr, user, resTask.Result)

	return protocol.ERROR_CODE_NONE
}

func _sendLoginResponse(svr *Server, user *gameUser, result int16) {
	res := protocol.LoginResPacket {
		Result: result,
	}

	outBuf := user.getBuffer(32)
	resPkt, resPktSize := res.EncodingPacket(outBuf)
	user.aheadWriteCursor(int(resPktSize))

	svr._serverNet.ISendToClient(user.netIndex(), resPkt)
}
