package omokServer

import "sync/atomic"

var _seqNumber uint64

func SeqNumIncrement() uint64 {
	newValue := atomic.AddUint64(&_seqNumber, 1)
	return newValue
}


type gameUser struct {
	sessionIndex int
	UID uint64
	ID string
	isAuth bool
}

func (u *gameUser) isEnableLogin() bool {
	return u.isAuth == false
}

func (u *gameUser) setAuth() {
	u.isAuth = true
}



type gameUserManager struct {
	users map[int]*gameUser
}

func newGameUserManager() *gameUserManager {
	userMgr := new(gameUserManager)
	userMgr.users = make(map[int]*gameUser)

	return userMgr
}

func (mgr *gameUserManager) addUser(sessionIndex int) bool {
	_, exists := mgr.users[sessionIndex]
	if exists {
		return false
	}

	user := new(gameUser)
	user.sessionIndex = sessionIndex
	user.UID = SeqNumIncrement()

	mgr.users[sessionIndex] = user
	return true
}

func (mgr *gameUserManager) removeUser(sessionIndex int) {
	delete(mgr.users, sessionIndex)
}

func (mgr *gameUserManager) GetUser(sessionIndex int) (*gameUser, bool) {
	user, exists := mgr.users[sessionIndex]
	if exists {
		return user, true
	}
	return nil, false
}