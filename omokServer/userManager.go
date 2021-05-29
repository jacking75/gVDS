package omokServer

import (
	"math/rand"
	"scommon"
	"sync/atomic"
	"time"
)



type gameUserManager struct {
	_users map[int]*gameUser

	_seqNumber uint64

	_heartBeatRand *rand.Rand

	setDisableSendClient func(int)
	netSend func(int, []byte) bool
}

func newGameUserManager() *gameUserManager {
	userMgr := new(gameUserManager)
	userMgr._users = make(map[int]*gameUser)

	s := rand.NewSource(time.Now().UnixNano())
	userMgr._heartBeatRand = rand.New(s)
	return userMgr
}

func (mgr *gameUserManager) createUID() uint64 {
	newValue := atomic.AddUint64(&mgr._seqNumber, 1)
	return newValue
}

func (mgr *gameUserManager) addUser(sessionIndex int, conf userConf) bool {
	_, exists := mgr._users[sessionIndex]
	if exists {
		return false
	}

	user := newGameUser(sessionIndex, mgr.createUID(), conf)
	user.netSend = mgr.netSend

	mgr._users[sessionIndex] = user
	return true
}

func (mgr *gameUserManager) removeUser(sessionIndex int) {
	delete(mgr._users, sessionIndex)
}

func (mgr *gameUserManager) GetUser(sessionIndex int) (*gameUser, bool) {
	user, exists := mgr._users[sessionIndex]
	if exists {
		return user, true
	}
	return nil, false
}

func (mgr *gameUserManager) checkUserState() {
	curTime := scommon.CurrentUnixMillSec()

	for _, user := range mgr._users {
		if mgr.checkHeartBeat(user, curTime) {
			continue
		}
	}
}

func (mgr *gameUserManager) checkHeartBeat(user *gameUser, curTime int64) bool {
	// 허트 비트 조사
	if user.isOverTimeHeartbeatResEndTime(curTime) {
		//구현 하지 않았지만 클라이언트에게 send를 보내지 않게 했으면 애플리케이션에서 이 클라이언트에게 보내는 패킷을 따로 저장하고 있다가 send 가능하게 되면(허트비트 응답을 받으면) 저장하고 있던 데이터를 다 보내도록 한다
		mgr.setDisableSendClient(user.netIndex())

		scommon.LogInfo("isOverTimeHeartbeatResEndTime")
		return true
	}

	if user.enableHeartbeatReqTime(curTime) {
		expected := mgr._heartBeatRand.Int63()
		user.sendHeartBeatRes(expected)

		scommon.LogDebug("sendHeartBeatRes")
		return true
	}

	return false
}