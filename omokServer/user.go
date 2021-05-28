package omokServer

import "scommon"

type userConf struct {
	userSendBufferSzie int
	heartbeatReqIntervalTimeMSec int64 // 몇 밀리세컨드 간격으로 클라에게 허트비트를 보낼지
	heartbeatWaitTimeMSec int64        // 허트비트 보낸 후 답변을 기다리는 최대 대기 시간. 밀리세컨드
}

type gameUser struct {
	sessionIndex int
	UID uint64
	ID string
	isAuth bool

	conf userConf

	heartbeatNextReqTimeMSec int64
	heartbeatResValue int64
	heartbeatResEndTimeMsec int64

	sendBuffer *sendRingBuffer
}

func newGameUser(sessionIndex int, UID uint64, conf userConf) *gameUser {
	user := new(gameUser)
	user.conf = conf
	user.sessionIndex = sessionIndex
	user.UID = UID
	user.sendBuffer = newSendBuffer(conf.userSendBufferSzie)

	user.heartbeatNextReqTimeMSec = scommon.NextUnixMillSec(conf.heartbeatReqIntervalTimeMSec)
	user.heartbeatResValue = 0
	user.heartbeatResEndTimeMsec = 0
	return user
}

func (u *gameUser) netIndex() int {
	return u.sessionIndex
}

func (u *gameUser) isEnableLogin() bool {
	return u.isAuth == false
}

func (u *gameUser) setAuth() {
	u.isAuth = true
}

func (u *gameUser) getBuffer(requiredSize int) []byte {
	return u.sendBuffer.getBuffer(requiredSize)
}

func (u *gameUser) aheadWriteCursor(size int) {
	u.sendBuffer.aheadWriteCursor(size)
}
