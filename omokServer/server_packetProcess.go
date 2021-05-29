package omokServer

import (
	"fmt"
	"omokServer/protocol"
	"scommon"
)


func (svr *Server) packetProcess(sessionIndex int, packetData []byte) {
	packetID := protocol.PeekPacketID(packetData)
	_, bodyData := protocol.PeekPacketBody(packetData)

	if pfunc := svr.getPacketFunc(packetID); pfunc != nil {
		if user, ok := svr._userMgr.GetUser(sessionIndex); ok {
			pfunc(user, bodyData)
		} else {
			scommon.LogError(fmt.Sprintf("[packetProcess] invalid User. _sessionIndex: %d", sessionIndex))
		}
	} else {
		scommon.LogError(fmt.Sprintf("[packetProcess] invalid packetID: %d", packetID))
	}
}

func (svr *Server) getPacketFunc(packetID int16) func(*gameUser, []byte) int16 {
	for i, id := range svr._funcPackeIdlist {
		if id == packetID {
			return svr._funclist[i]
		}
	}

	return nil
}

func (svr *Server) settingPacketFunction() {
	maxFuncListCount := 16
	svr._funclist = make([]func(*gameUser, []byte) int16, 0, maxFuncListCount)
	svr._funcPackeIdlist = make([]int16, 0, maxFuncListCount)

	svr._addPacketFunction(protocol.PACKET_ID_LOGIN_REQ, svr.packetProcessLogin)
}

func (svr *Server) _addPacketFunction(packetID int16,
						packetFunc func(*gameUser, []byte) int16) {
	svr._funclist = append(svr._funclist, packetFunc)
	svr._funcPackeIdlist = append(svr._funcPackeIdlist, packetID)
}
