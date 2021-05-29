package omokServer

import "sync/atomic"



type gameUserManager struct {
	_users map[int]*gameUser

	_seqNumber uint64
}

func newGameUserManager() *gameUserManager {
	userMgr := new(gameUserManager)
	userMgr._users = make(map[int]*gameUser)

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
	//TODO 접속된 유저들의 상태를 조사한다

	//TODO 허트비트 구현하기.
	// 서버에서 랜던 값을 보내면 클라는 이것보다 1 큰 값을 보내야 한다.
	// 허트비트가 오지 않으면 네트워크 동작을 멈춘다.
	// 유저에게 보내는 send를 고루틴으로 만들지 않기 때문에 특정 유저가 send 버퍼가 다 차서 보내지 못하면 블럭킹이 발생한다. 이것을 방지하기 위해서 허트비트는 필수이다.
}