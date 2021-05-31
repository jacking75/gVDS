//TODO 사용할 때는 코드 수정이 필요
package smallNet

import (
	"fmt"
	"scommon"
	"sync"
	"sync/atomic"
	"time"
)

type tcpServerSessionInfo struct {
	isFree          bool
	disabledTimeSec int64
	session         *remoteSession
}

type tcpServerSessionManager struct {
	_netConfig NetworkConfig

	_pktRecvFunc PacketReceivceFunctors

	_curSessionCount int32 // 멀티스레드에서 호출된다

	_sessionAllocFreeLock *sync.Mutex
	_sessionList          []tcpServerSessionInfo // 멀티스레드에서 호출된다
}


func newServerSessionManager(config NetworkConfig,
							pktRecvFunc PacketReceivceFunctors) *tcpServerSessionManager {
	sessionMgr := new(tcpServerSessionManager)
	sessionMgr._netConfig = config
	sessionMgr._pktRecvFunc = pktRecvFunc
	sessionMgr._curSessionCount = 0
	sessionMgr._initSessionList()

	scommon.LogInfo(fmt.Sprintf("New tcpServerSessionManager. maxSessionCount(%d)", sessionMgr._netConfig.MaxSessionCount))

	return sessionMgr
}

func (smgr *tcpServerSessionManager) connectAndReceive(serverName string, remoteAddress string) bool {
	newSession := smgr._allocSession(serverName)

	if newSession == nil {
		return false
	}

	go newSession.connectAndReceive_goroutine(remoteAddress)

	return true
}

func (smgr *tcpServerSessionManager) stop() {
	scommon.LogInfo(fmt.Sprintf("Stop ServerSessionManager. sessionCount(%d)", smgr._curSessionCount))

	for i := 0; i < smgr._netConfig.MaxSessionCount; i++ {
		if smgr._sessionList[i].isFree {
			continue
		}

		smgr._sessionList[i].session.stop(sessionCloseAllSession)
	}
}

func (smgr *tcpServerSessionManager) waitAllSessionStop() bool {
	for {
		stopCount := int32(0)

		for i := 0; i < smgr._netConfig.MaxSessionCount; i++ {
			if smgr._sessionList[i].isFree {
				continue
			}

			if smgr._sessionList[i].session.isRunningToStoped() {
				stopCount++
			}
		}

		if stopCount == smgr._curSessionCount {
			return true
		} else {
			scommon.LogInfo("[waitAllSessionStop] end")
			time.Sleep(time.Second)
		}
	}

	return false
}

func (smgr *tcpServerSessionManager) sendPacket(sessionIndex int, data []byte) bool {
	session, result := smgr._findSession(sessionIndex)
	if result {
		session.sendPacket(data)
		return true
	}

	return false
}

func (smgr *tcpServerSessionManager) sendPacketAll(data []byte) {
	for i := 0; i < smgr._netConfig.MaxSessionCount; i++ {
		if smgr._sessionList[i].isFree || smgr._sessionList[i].disabledTimeSec != 0 {
			continue
		}

		smgr._sessionList[i].session.sendPacket(data)
	}
}

func (smgr *tcpServerSessionManager) _incConnectedSessionCount() {
	atomic.AddInt32(&smgr._curSessionCount, 1)
}

func (smgr *tcpServerSessionManager) _decConnectedSessionCount() {
	atomic.AddInt32(&smgr._curSessionCount, -1)
}

func (smgr *tcpServerSessionManager) _initSessionList() {
	smgr._sessionAllocFreeLock = new(sync.Mutex)
	smgr._sessionList = make([]tcpServerSessionInfo, smgr._netConfig.MaxSessionCount)

	for i := 0; i < smgr._netConfig.MaxSessionCount; i++ {
		smgr._sessionList[i].isFree = true
		smgr._sessionList[i].disabledTimeSec = 1
	}
}

func (smgr *tcpServerSessionManager) _allocSession(serverName string) *remoteSession {
	smgr._sessionAllocFreeLock.Lock()
	defer smgr._sessionAllocFreeLock.Unlock()

	// 할당되지 않은 세션이라도 비 사용으로 된 시간과 지금 시간이 지정 시간이 지난 뒤에만 사용 가능하다.
	// 이유는 멀티스레드의 미묘한 동기화 문제를 피하기 위해 타이밍을 지연시킨다
	waitTimeSec := int64(3)
	curTime := time.Now().Unix()

	for i := 0; i < smgr._netConfig.MaxSessionCount; i++ {
		if smgr._sessionList[i].isFree == false {
			continue
		}

		diffTime := curTime - smgr._sessionList[i].disabledTimeSec
		if diffTime < waitTimeSec {
			continue
		}

		newSession := new(remoteSession)
		newSession.initialize(i, serverName, smgr._netConfig, smgr._pktRecvFunc)

		smgr._sessionList[i].session = nil
		smgr._sessionList[i].session = newSession
		smgr._sessionList[i].disabledTimeSec = 0
		smgr._sessionList[i].isFree = false

		smgr._curSessionCount++

		return newSession
	}

	return nil
}

func (smgr *tcpServerSessionManager) _validIndex(index int) bool {
	maxSessionCount := smgr._netConfig.MaxSessionCount
	if index < 0 || index >= maxSessionCount {
		return false
	}

	return true
}

func (smgr *tcpServerSessionManager) _findSession(sessionIndex int) (*remoteSession, bool) {
	if smgr._validIndex(sessionIndex) == false {
		return nil, false
	}

	sessionInfo := smgr._sessionList[sessionIndex]

	if sessionInfo.isFree {
		return nil, false
	}

	return sessionInfo.session, true
}

func (smgr *tcpServerSessionManager) forceDisconnectServerSession(sessionIndex int) bool {
	if smgr._validIndex(sessionIndex) == false {
		return false
	}
	sessionInfo := smgr._sessionList[sessionIndex]

	if sessionInfo.isFree {
		return false
	}

	sessionInfo.isFree = true
	sessionInfo.disabledTimeSec = 1

	sessionInfo.session.setStateClosed()
	_ = sessionInfo.session._tcpConn.Close()

	scommon.LogInfo(fmt.Sprintf("[forceDisconnectServerSession] TcpServerSession Close. sessionIndex: %d", sessionIndex))
	return true
}