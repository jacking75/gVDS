package smallNet

import (
	"fmt"
	"net"
	"scommon"
)

type tcpClientSessionManager struct {
	_maxSessionCount int
	_curSessionCount int

	_netConf     NetworkConfig
	_pktRecvFunc PacketReceivceFunctors

	_sessionList      []*tcpSession // 멀티스레드에서 호출된다
	_sessionIndexPool *scommon.Deque
}

func newClientSessionManager(config NetworkConfig,
							pktRecvFunc PacketReceivceFunctors) *tcpClientSessionManager {
	sessionMgr := new(tcpClientSessionManager)
	sessionMgr._pktRecvFunc = pktRecvFunc
	sessionMgr._initialize(config)

	return sessionMgr
}

func (mgr *tcpClientSessionManager) stop() {
	mgr._forceCloseAllSession()
}

func (mgr *tcpClientSessionManager) sendPacket(sessionIndex int, sendData []byte) bool {
	session, result := mgr.findSession(sessionIndex)
	if result == false || session.enableSend() == false {
		return false
	}

	session.sendPacket(sendData)
	return true
}

func (mgr *tcpClientSessionManager) sendPacketAllClient(sendData []byte) {
	for i := 0; i < mgr._maxSessionCount; i++ {
		session := mgr._sessionList[i]
		if session.enableSend() == false {
			continue
		}

		session.sendPacket(sendData)
	}
}

func (mgr *tcpClientSessionManager) forceDisconnectClient(sessionIndex int) {
	session, result := mgr.findSession(sessionIndex)
	if result == false || session.enableSend() == false {
		return
	}

	session.closeSocket(sessionCloseForce)
}

func (mgr *tcpClientSessionManager) disablePacketProcessClient(sessionIndex int) {
	session, result := mgr.findSession(sessionIndex)
	if result == false || session.enableSend() == false {
		return
	}

	session.disablePacketProcess()
}

func (mgr *tcpClientSessionManager) _initialize(config NetworkConfig) {
	mgr._netConf = config
	mgr._maxSessionCount = config.MaxSessionCount
	mgr._curSessionCount = 0

	mgr._createSessionPool()
}

func (mgr *tcpClientSessionManager) _createSessionPool() {
	mgr._sessionList = make([]*tcpSession, mgr._maxSessionCount)
	mgr._sessionIndexPool = scommon.NewCappedDeque((int)(mgr._maxSessionCount))

	for i := 0; i < (int)(mgr._maxSessionCount); i++ {
		mgr._sessionList[i] = new(tcpSession)
		mgr._sessionList[i].initialize(i, true, mgr._netConf)
		mgr._freeSessionObj(i)
	}
}

func (mgr *tcpClientSessionManager) _connectedSessionCount() int {
	return mgr._curSessionCount
}

func (mgr *tcpClientSessionManager) _incConnectedSessionCount() {
	mgr._curSessionCount++
}

func (mgr *tcpClientSessionManager) _decConnectedSessionCount() {
	mgr._curSessionCount--
}

func (mgr *tcpClientSessionManager) _validIndex(index int) bool {
	if index < 0 || index >= mgr._maxSessionCount {
		return false
	}
	return true
}

func (mgr *tcpClientSessionManager) findSession(sessionIndex int) (*tcpSession, bool) {
	if mgr._validIndex(sessionIndex) == false {
		return nil, false
	}

	session := mgr._sessionList[sessionIndex]

	if session == nil {
		return nil, false
	}

	return session, true
}

func (mgr *tcpClientSessionManager) newSession(tcpconn *net.TCPConn) int {
	newSession := mgr._allocSessionObj()
	if newSession == nil {
		scommon.LogError(fmt.Sprintf("[tcpClientSessionManager._newSession] empty SessionObj"))
		_ = tcpconn.Close()
		return -1
	}

	index := newSession.getIndex()

	newSession.clear()
	newSession.onConnect(tcpconn)
	newSession.settingTCPSocketOption(int(mgr._netConf.SockReadbuf), int(mgr._netConf.SockWritebuf))

	go newSession.handleReceive_goroutine(mgr._netConf, mgr._pktRecvFunc)

	return index
}

func (mgr *tcpClientSessionManager) deleteSession(sessionIndex int) {
	if session, ret := mgr.findSession(sessionIndex); ret == true {
		session.setDisableSend()

		mgr._decConnectedSessionCount()

		session.setStateClosed()

		mgr._freeSessionObj(sessionIndex)
	}
}

func (mgr *tcpClientSessionManager) _allocSessionObj() *tcpSession {
	item := mgr._sessionIndexPool.First()
	if item == nil {
		return nil
	}

	sessionIndex := item.(int)
	data := mgr._sessionList[sessionIndex]
	_ = mgr._sessionIndexPool.Shift()

	scommon.LogTrace(fmt.Sprintf("[tcpClientSessionManager] _allocSessionObj(%d)", sessionIndex))
	return data
}

func (mgr *tcpClientSessionManager) _freeSessionObj(sessionIndex int) {
	if useCount, ok := mgr._sessionIndexPool.Append(sessionIndex); ok {
		scommon.LogTrace(fmt.Sprintf("[tcpClientSessionManager] _freeSessionObj. sessionIndex(%d)UseCount( %d )",
			sessionIndex, useCount))
	} else {
		scommon.LogError(fmt.Sprintf("[tcpClientSessionManager] _freeSessionObj( %d )", sessionIndex))
		return
	}
}

func (mgr *tcpClientSessionManager) _forceCloseAllSession() {
	scommon.LogInfo("_forceCloseAllSession - start")

	for _, session := range mgr._sessionList {
		if session == nil {
			continue
		}

		session.closeSocket(sessionCloseAllSession)
	}

	scommon.LogInfo("_forceCloseAllSession - end")
}




