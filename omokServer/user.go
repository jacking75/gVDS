package omokServer

import (
	"omokServer/protocol"
	"scommon"
)

type userConf struct {
	userSendBufferSzie int
	heartbeatReqIntervalTimeMSec int64 // 몇 밀리세컨드 간격으로 클라에게 허트비트를 보낼지
	heartbeatWaitTimeMSec int64        // 허트비트 보낸 후 답변을 기다리는 최대 대기 시간. 밀리세컨드
}

type gameUser struct {
	_sessionIndex int
	_uid          uint64
	_id           string
	_isAuth       bool

	conf userConf

	_heartbeatNextReqTimeMSec int64
	_heartbeatResValue        int64
	_heartbeatResEndTimeMSec  int64

	_sendBuffer *sendRingBuffer

	netSend func(int, []byte) bool
}

func newGameUser(sessionIndex int, UID uint64, conf userConf) *gameUser {
	user := new(gameUser)
	user.conf = conf
	user._sessionIndex = sessionIndex
	user._uid = UID
	user._sendBuffer = newSendBuffer(conf.userSendBufferSzie)

	user._heartbeatNextReqTimeMSec = scommon.NextUnixMillSec(conf.heartbeatReqIntervalTimeMSec)
	user._heartbeatResValue = 0
	user._heartbeatResEndTimeMSec = 0
	return user
}

func (u *gameUser) netIndex() int {
	return u._sessionIndex
}

func (u *gameUser) isEnableLogin() bool {
	return u._isAuth == false
}

func (u *gameUser) setAuth(id string) {
	u._isAuth = true
	u._id = id
}

func (u *gameUser) getBuffer(requiredSize int) []byte {
	return u._sendBuffer.getBuffer(requiredSize)
}

func (u *gameUser) aheadWriteCursor(size int) {
	u._sendBuffer.aheadWriteCursor(size)
}

func (u *gameUser) isOverTimeHeartbeatResEndTime(curTime int64) bool {
	if u._heartbeatResEndTimeMSec != 0 && u._heartbeatResEndTimeMSec <= curTime {
		return true
	}
	return false
}

func (u *gameUser) enableHeartbeatReqTime(curTime int64) bool {
	if u._heartbeatNextReqTimeMSec <= curTime {
		return true
	}
	return false
}

func (u *gameUser) sendHeartBeatRes(expected int64) {
	req := protocol.HeartBeatReqPacket {
		Expected: expected,
	}

	outBuf := u.getBuffer(16)
	reqPkt, reqPktSize := req.EncodingPacket(outBuf)
	u.aheadWriteCursor(int(reqPktSize))
	u.netSend(u.netIndex(), reqPkt)

	u._heartbeatResEndTimeMSec = scommon.NextUnixMillSec(u.conf.heartbeatWaitTimeMSec)
	u._heartbeatResValue = expected
	u._heartbeatNextReqTimeMSec = 0
}

func (u *gameUser) checkHeartBeat(expected int64) bool {
	if expected < 1 || expected != u._heartbeatResValue {
		return false
	}

	u._heartbeatResEndTimeMSec = 0
	u._heartbeatResValue = 0
	u._heartbeatNextReqTimeMSec = scommon.NextUnixMillSec(u.conf.heartbeatReqIntervalTimeMSec)

	return true
}