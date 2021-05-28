package omokServer

/*
var _seqNumber uint64

func SeqNumIncrement() uint64 {
	newValue := atomic.AddUint64(&_seqNumber, 1)
	return newValue
}
*/

type gameUser struct {
	sessionIndex int
	ID string
	isAuth bool
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

	mgr.users[sessionIndex] = user
	return true
}

func (mgr *gameUserManager) removeUser(sessionIndex int) {
	delete(mgr.users, sessionIndex)
}