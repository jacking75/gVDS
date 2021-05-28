package smallNet

import "net"

type sessionError int16

const netLibErrNone = 0

const (
	sessionErrStart = sessionError(iota) // DUMMY

	sessionCloseForce
	sessionCloseAllSession
	sessionCloseRecvGoroutineEnd
	sessionCloseForceTerminateRecvGoroutine
	sessionCloseForceTerminateSendGoroutine
	sessionCloseCloseRemote
	sessionCloseSocketError
	sessionCloseSendChannelIsFull
	sessionCloseSocketReadTimeout
	sessionCloseSocketFailCallReadDeadLine

	sessionCloseRecvMakePacketTooLargePacketSize
	sessionCloseRecvTooSmallData
	sessionDisablePacketProcess
	sessionCloseRingBufferErr

	ringBufferRecvInitFail
	ringBufferSendInitFail
	ringBufferWriteGatherPacket
	ringBufferMaxSizeGreaterMaxPacketSize
)

func (s sessionError) Error() string {
	return _closeCaseMessage[s]
}

var _closeCaseMessage = [...]string{
	sessionErrStart: "",

	sessionCloseForce:                       "session close force",
	sessionCloseAllSession:                  "session close all session",
	sessionCloseRecvGoroutineEnd:            "session close recv goroutine end",
	sessionCloseForceTerminateRecvGoroutine: "session close force terminate recv goroutine",
	sessionCloseForceTerminateSendGoroutine: "session close force terminate send goroutine",
	sessionCloseCloseRemote:                 "session close - close remote",
	sessionCloseSocketError:                 "session close socket error",
	sessionCloseSendChannelIsFull:           "session close send channel is full",
	sessionCloseSocketReadTimeout:           "session close socket read timeout",
	sessionCloseSocketFailCallReadDeadLine:  "session close socket fail call read deadline",

	sessionCloseRecvMakePacketTooLargePacketSize: "session close recv make packet too large packet size",
	sessionCloseRecvTooSmallData:                 "session close recv too small data",
	sessionDisablePacketProcess:				"session disable packet process",
	sessionCloseRingBufferErr:                    "session close ring buffer error",

	ringBufferRecvInitFail: "ringBufferRecvInitFail",
	ringBufferSendInitFail: "ringBufferSendInitFail",
	ringBufferWriteGatherPacket: "ring buffer write gather packet",
	ringBufferMaxSizeGreaterMaxPacketSize:     "ring buffer max size greater max packet size",

}


var	NetMsg_None int8 = 0
var	NetMsg_Connect int8 = 1
var	NetMsg_Close int8 = 2
var	NetMsg_Receive int8 = 3

type NetMsg struct {
	Type 			int8
	SessionIndex    int
	Data            []byte
	TcpConn			*net.TCPConn
}


type PacketReceivceFunctors struct {
	AddNetMsgOnCloseFunc func(int)
	AddNetMsgOnReceiveFunc func(int, []byte)

	// 데이터를 분석하여 패킷 크기를 반환한다.
	PacketTotalSizeFunc func([]byte) int16

	// 패킷 헤더의 크기
	PacketHeaderSize int16

	// true 이면 client와 연결한 세션이다.
	IsClientSession bool
}