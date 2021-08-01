package main

import (
	"redisDB"
)

func processRedisResTask_goroutine(svr *echoServer) {
	for {
		resTask := <- svr._resTaskChan

		if resTask.ID == redisDB.TaskID_ResLogin {
			_sendLoginResponse(svr, int(resTask.UID), resTask.Result)
		}
	}
}

func _sendLoginResponse(svr *echoServer, sessionIndex int, result int16) {
	res := LoginResPacket {
		Result: result,
	}

	outBuf := make([]byte, 32)
	resPkt, _ := res.EncodingPacket(outBuf)

	svr._serverNet.ISendToClient(sessionIndex, resPkt)
}