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
		pfunc(nil, bodyData)
	} else {
		scommon.LogError(fmt.Sprintf("invalid packetID: %d", packetID))
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

	//svr._addPacketFunction(protocol.PACKET_ID_ROOM_ENTER_REQ, svr._packetProcess_EnterUser)
	//svr._addPacketFunction(protocol.PACKET_ID_ROOM_LEAVE_REQ, svr._packetProcess_LeaveUser)
	//svr._addPacketFunction(protocol.PACKET_ID_ROOM_CHAT_REQ, svr._packetProcess_Chat)
	svr._addPacketFunction(protocol.PACKET_ID_ROOM_RELAY_REQ, svr._packetProcess_Relay)
}

func (svr *Server) _addPacketFunction(packetID int16,
						packetFunc func(*gameUser, []byte) int16) {
	svr._funclist = append(svr._funclist, packetFunc)
	svr._funcPackeIdlist = append(svr._funcPackeIdlist, packetID)
}






func (svr *Server) _packetProcess_Relay(user *gameUser, bodyData []byte) int16 {
	/*var relayNotify protocol.RoomRelayNtfPacket
	relayNotify.RoomUserUniqueId = user.RoomUniqueId
	relayNotify.Data = packet.Data
	notifySendBuf, packetSize := relayNotify.EncodingPacket(packet.DataSize)
	room.broadcastPacket(packetSize, notifySendBuf, 0)

	NTELIB_LOG_DEBUG("Room Relay", zap.String("Sender", string(user.ID[:])))*/
	return protocol.ERROR_CODE_NONE
}