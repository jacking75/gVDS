//TODO 사용할 때는 코드 수정이 필요
package smallNet

import (
	"fmt"
	"net"
	"scommon"
	"time"
)

type remoteSession struct {
	tcpSession

	_remoteServerName    string
	_remoteServerAddress string
	_maxConnectTryCount  int

	_maxPacketSize              int16
	_maxReceiveBufferSize       int
	_sockReadbuf, _sockWritebuf int // 소켓 버퍼 크기-받기, 소켓 버퍼 크기-보내기

	_pktRecvfunc PacketReceivceFunctors

	_closing         bool
	_runningToStoped bool
}


func (session *remoteSession) initialize(index int,
						remoteServerName string,
						config NetworkConfig,
						pktRecvfunc PacketReceivceFunctors) sessionError {

	if err := session.tcpSession.initialize(index, false, config); err != netLibErrNone {
		return err
	}

	session._remoteServerName = remoteServerName
	session._runningToStoped = false
	session._maxPacketSize = (int16)(config.MaxPacketSize)
	session._sockReadbuf = (int)(config.SockReadbuf)
	session._sockWritebuf = (int)(config.SockWritebuf)
	session._maxConnectTryCount = 0 //무한 반복하도록 일단 0 으로 고정한다
	session._pktRecvfunc = pktRecvfunc
	session._closing = false

	return netLibErrNone
}

func (session *remoteSession) onConnect() {
	session._setStateConnect()

	session._isEnableSend = true
	session._isEnablePacketProcess = true
}

func (session *remoteSession) stop(reason sessionError) {
	session._closing = true
	session.closeSocket(reason)
}

func (session *remoteSession) isRunningToStoped() bool {
	return session._runningToStoped
}

func (session *remoteSession) connectAndReceive_goroutine(remoteAddress string) {
	scommon.LogInfo(fmt.Sprintf("[connectAndReceive_goroutine] Start. ServerSession( %d)", session.getIndex()))

	session._remoteServerAddress = remoteAddress

	for {
		if session.connectAndReceive_goroutine_impl(remoteAddress) {
			scommon.LogInfo(fmt.Sprintf("[connectAndReceive_goroutine] Wanted Stop . ServerSession( %d)", session.getIndex()))
			break
		}
	}

	session._runningToStoped = true
	scommon.LogInfo(fmt.Sprintf("[connectAndReceive_goroutine] Stop . ServerSession( %d)", session.getIndex()))
}

func (session *remoteSession) connectAndReceive_goroutine_impl(remoteAddress string) bool {
	IsWantedTermination := false
	defer scommon.PrintPanicStack()
	defer session.closeSocket(sessionCloseRecvGoroutineEnd)

	reTryCount := 0

	for {
		if IsWantedTermination {
			break
		}

		reTryCount++
		if session._maxConnectTryCount > 0 && reTryCount > session._maxConnectTryCount {
			break
		}

		if session._tcpConnectToRemote() {
			scommon.LogInfo(fmt.Sprintf("[connectAndReceive_goroutine_impl] Connected: Remote Server. ServerSession( %d)", session.getIndex()))
		} else {
			IsWantedTermination = session._closing
			time.Sleep(time.Second * 5)
			continue
		}

		session.onConnect()
		session.settingTCPSocketOption(session._sockReadbuf, session._sockWritebuf)

		//sessionIndex := session.getIndex()
		//session._networkFuncDelegate.OnConnect(sessionIndex)

		//TOOD 옵션에 의해서 자동으로 재연결(현재 자동 재연결. 끊어짐도 알리지 않음)을 선택할 수 있게 한다
		err := session.handleTCPReceive(int(session._maxPacketSize),
										session._pktRecvfunc)
		if err != netLibErrNone {
			//sessionID, uniqueIDNum := session.getIndexUniqueID()
			session.closeSocket(err)

			//session._networkFuncDelegate.OnClose(sessionID, uniqueIDNum)

			if session._closing == false {
				scommon.LogError(fmt.Sprintf("[connectAndReceive_goroutine_impl] Close ServerSeession. Session( %d ), Reason: %s", session.getIndex(), err.Error()))
			}
		}

		//동기화를 위해 좀 대기한 후 다시 접속 시도한다
		if session._closing == false {
			session.setDisableSend()

			time.Sleep(time.Second * 5)

			session._tcpConn = nil
		} else {
			IsWantedTermination = true
		}
	}

	scommon.LogInfo(fmt.Sprintf("[connectAndReceive_goroutine_impl] end. Session( %d ), result(%t)", session.getIndex(), IsWantedTermination))
	return IsWantedTermination
}

func (session *remoteSession) _tcpConnectToRemote() bool {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", session._remoteServerAddress)
	if err != nil {
		scommon.LogError(fmt.Sprintf("[_tcpConnectToRemote] fail ResolveTCPAddr Remote. Session( %d ), Address(%s), %v", session.getIndex(), session._remoteServerAddress, err))
		return false
	}

	session._tcpConn, err = net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		scommon.LogError(fmt.Sprintf("[_tcpConnectToRemote] fail connect Remote. Session( %d ), %v", session.getIndex(), err))
		return false
	}

	return true
}

func (session *remoteSession) sendPacket(data []byte) {
	if session.isStateConnect() && session._isEnableSend {
		_, err := session._tcpConn.Write(data)
		if err != nil {
			if nErr, ok := err.(net.Error); ok && nErr.Temporary() {
				scommon.LogError(fmt.Sprintf("[sendPacket] Temporary Error remoteSession Send. Session( %d ), %v", session.getIndex(), err))
				return
			}

			scommon.LogError(fmt.Sprintf("E[sendPacket] rror remoteSession Send. Session( %d ), %v", session.getIndex(), err))
		}
	}
}

