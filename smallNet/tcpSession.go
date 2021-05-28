package smallNet

import (
	"fmt"
	"net"
	"scommon"
)

type tcpSession struct {
	_index int

	_ip          net.IP
	_tcpConn     *net.TCPConn

	_isTcpConnected                 bool // tcp 연결여부
	_isEnableSend                   bool // send 가능 여부
	_isEnablePacketProcess          bool // 패킷 처리 가능 여부
	_isClientSession                bool // 클라이언트 세션 여부

	_recvBuffer *ringBuffer
	_sendBuffer *ringBuffer
}

func (session *tcpSession) initialize(index int,
									isClientSession bool,
									conf NetworkConfig) sessionError {
	session._index = index
	session._isClientSession = isClientSession
	err := session.initRingBuffer(conf)
	if err != netLibErrNone {
		return err
	}

	return netLibErrNone
}

// 받기, 보내기 모두 링버퍼를 사용하도록 한다
func (session *tcpSession) initRingBuffer(config NetworkConfig) sessionError {
	var buffErr ringbufferErr

	recvPacketRingBufferMaxSize := config.RecvPacketRingBufferMaxSize
	session._recvBuffer, buffErr = newRingBuffer(recvPacketRingBufferMaxSize, config.MaxPacketSize)
	if buffErr != err_ringbuffer_none {
		scommon.LogError(fmt.Sprintf("[initRingBuffer] Recv NewPacketRingBuffer. Session(%d), %d", session.getIndex(), buffErr))
		return ringBufferRecvInitFail
	}

	sendPacketRingBufferMaxSize := config.SendPacketRingBufferMaxSize
	session._sendBuffer, buffErr = newRingBuffer(sendPacketRingBufferMaxSize, config.MaxPacketSize)
	if buffErr != err_ringbuffer_none {
		scommon.LogError(fmt.Sprintf("[initRingBuffer] Send NewPacketRingBuffer. Session(%d), %d", session.getIndex(), buffErr))
		return ringBufferSendInitFail
	}

	return netLibErrNone
}

func (session *tcpSession) onConnect(conn *net.TCPConn) {
	host, _, err := net.SplitHostPort(conn.RemoteAddr().String())
	if err != nil {
		scommon.LogError(fmt.Sprintf("[onConnect] cannot get remote address. Session( %d ), %v", session.getIndex(), err))
		return
	}

	session._tcpConn = conn
	session._ip = net.ParseIP(host)

	session._setStateConnect()

	session._isEnableSend = true
	session._isEnablePacketProcess = true
}

func (session *tcpSession) clear() {
	session.setDisableSend()
	session._tcpConn = nil
	session._ip = nil
	session._recvBuffer.reset()
	session._sendBuffer.reset()
}

func (session *tcpSession) getIndex() int {
	return session._index
}

func (session *tcpSession) getSocket() net.Conn {
	return session._tcpConn
}

func (session *tcpSession) sendPacket(data []byte) {
	if session.isStateConnect() == false || session._isEnableSend == false {
		return
	}

	if _, ret := session._tcpConn.Write(data); ret != nil {
		scommon.LogError(fmt.Sprintf("[sendPacket] Error clientSession sendPacket. Session( %d ), %v", session.getIndex(), ret))
	}
}

func (session *tcpSession) closeSocket(err sessionError) {
	if session.isStateConnect() == false {
		return
	}

	session.setStateClosed()
	_ = session._tcpConn.Close()
	scommon.LogInfo(fmt.Sprintf("[closeSocket] tcpSession Close. Session( %d ), %v", session.getIndex(), err.Error()))
}

func (session *tcpSession) setStateClosed() {
	session._isTcpConnected = false
}

func (session *tcpSession) forceTerminateReceiveGoroutine() {
	session.closeSocket(sessionCloseForceTerminateRecvGoroutine)
}

func (session *tcpSession) setDisableSend() {
	session._isEnableSend = false
}

func (session *tcpSession) enableSend() bool {
	return session._isTcpConnected && session._isEnableSend
}

func (session *tcpSession) isEnablePacketProcess() bool {
	return session._isEnablePacketProcess
}

func (session *tcpSession) disablePacketProcess() {
	session._isEnablePacketProcess = false
}

func (session *tcpSession) settingTCPSocketOption(readBufSize int, writeBufSize int) {
	if readBufSize > 0 {
		_ = session._tcpConn.SetReadBuffer(readBufSize)
	}

	if writeBufSize > 0 {
		_ = session._tcpConn.SetWriteBuffer(writeBufSize)
	}
}

func (session *tcpSession) isStateConnect() bool {
	return session._isTcpConnected
}

func (session *tcpSession) _setStateConnect() {
	session._isTcpConnected = true
}

