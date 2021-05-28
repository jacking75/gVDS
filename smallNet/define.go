package smallNet

import "net"

type sessionError int16

const netLibErrNone = 0

const (
	sessionErrStart = sessionError(iota)

	sessionCloseForce
	sessionCloseAllSession
	sessionCloseRecvGoroutineEnd
	sessionCloseForceTerminateRecvGoroutine
	sessionCloseCloseRemote
	sessionCloseSocketError
	sessionCloseSocketReadTimeout
	sessionCloseRecvMakePacketTooLargePacketSize
	sessionCloseRecvTooSmallData
	sessionDisablePacketProcess
	ringBufferRecvInitFail
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
	sessionCloseCloseRemote:                 "session close - close remote",
	sessionCloseSocketError:                 "session close socket error",
	sessionCloseSocketReadTimeout:           "session close socket read timeout",
	sessionCloseRecvMakePacketTooLargePacketSize: "session close recv make packet too large packet size",
	sessionCloseRecvTooSmallData:                 "session close recv too small data",
	sessionDisablePacketProcess:				"session disable packet process",
	ringBufferRecvInitFail: "ringBufferRecvInitFail",
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