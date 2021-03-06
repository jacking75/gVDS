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

	_isConnected           bool // tcp 연결여부
	_isEnableSend          bool // send 가능 여부
	_isEnablePacketProcess bool // 패킷 처리 가능 여부
	_isClientSession       bool // 클라이언트 세션 여부

	_recvBuffer *ringBuffer
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
	session._recvBuffer = newRingBuffer(config.RecvPacketRingBufferMaxSize, config.MaxPacketSize)
	if session._recvBuffer == nil {
		scommon.LogError(fmt.Sprintf("[initRingBuffer] Recv NewPacketRingBuffer. Session(%d)", session.getIndex()))
		return ringBufferRecvInitFail
	}

	return netLibErrNone
}

func (session *tcpSession) onConnect(conn *net.TCPConn) {
	host, _, err := net.SplitHostPort(conn.RemoteAddr().String())
	if err != nil {
		scommon.LogError(fmt.Sprintf("[onConnect] cannot get remote address. Session( %d ), %v", session.getIndex(), err))
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
}

func (session *tcpSession) getIndex() int {
	return session._index
}

func (session *tcpSession) getSocket() net.Conn {
	return session._tcpConn
}

func (session *tcpSession) sendPacket(data []byte) bool {
	if session.isStateConnect() == false || session._isEnableSend == false {
		return false
	}

	if _, ret := session._tcpConn.Write(data); ret != nil {
		scommon.LogError(fmt.Sprintf("[sendPacket] Error clientSession sendPacket. Session( %d ), %v", session.getIndex(), ret))
		return false
	}

	return true
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
	session._isConnected = false
}

func (session *tcpSession) setDisableSend() {
	session._isEnableSend = false
}

func (session *tcpSession) setEnableSend() {
	session._isEnableSend = true
}

func (session *tcpSession) isEnableSend() bool {
	return session._isConnected && session._isEnableSend
}

func (session *tcpSession) isEnablePacketProcess() bool {
	return session._isEnablePacketProcess
}

func (session *tcpSession) setDisablePacketProcess() {
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
	return session._isConnected
}

func (session *tcpSession) _setStateConnect() {
	session._isConnected = true
}

