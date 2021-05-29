package omokServer

import (
	"omokServer/protocol"
)

func (svr *Server) packetProcessHeartBeat(user *gameUser, bodyData []byte) int16 {
	var resPkt protocol.HeartBeatResPacket
	resPkt.Decoding(bodyData)

	user.checkHeartBeat(resPkt.Expected)
	return protocol.ERROR_CODE_NONE
}


