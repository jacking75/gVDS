package smallNet

import (
	"fmt"
	"net"
	"runtime"
	"scommon"
)

func (session *tcpSession) handleReceive_goroutine(config NetworkConfig,
										pktRecvfunc PacketReceivceFunctors) {
	sessionIndex := session.getIndex()
	maxPacketSize := config.MaxPacketSize
	scommon.LogDebug(fmt.Sprintf("[handleClientTCPReceive goroutine] Start. Session( %d )", sessionIndex))

	defer pktRecvfunc.AddNetMsgOnCloseFunc(sessionIndex)
	//handleTCPReceive에서 패닉이 발생하면 여기에서 접속을 끊도록 한다. 일반적인 접속 종료이면 호출만 될 뿐이다.
	defer session.closeSocket(sessionCloseRecvGoroutineEnd)
	defer scommon.PrintPanicStack()

	err := session.handleTCPReceive(maxPacketSize, pktRecvfunc)
	session.closeSocket(err)

	scommon.LogDebug(fmt.Sprintf("[handleClientTCPReceive goroutine] End. Session(%d)", sessionIndex))
}

func (session *tcpSession) handleTCPReceive(maxPacketSize int,
									pktRecvfunc PacketReceivceFunctors) sessionError {
	packetHeaderSize := pktRecvfunc.PacketHeaderSize
	tcpConn := session.getSocket()

	for {
		receiveBuff := session._recvBuffer.getWriteBuffer(maxPacketSize)
		recvBytes, err := tcpConn.Read(receiveBuff)
		if recvBytes == 0 {
			return sessionCloseCloseRemote
		} else if serr := _checkReceiveError(session.getIndex(), err); serr != netLibErrNone {
			return serr
		} else {
			session._recvBuffer.aheadWriteCursor(recvBytes)
		}

		var dataBuff, readAbleByte = session._recvBuffer.readAbleBuffer()
		if serr := session._checkReadAfterEnablePacket(int16(readAbleByte), packetHeaderSize);
													serr == sessionDisablePacketProcess {
			session._recvBuffer.reset()
			continue
		} else if serr != netLibErrNone {
			return serr
		}

		if readSize, result := session._makePacketAndCallRecvEvent(readAbleByte,
											dataBuff, maxPacketSize,
											pktRecvfunc); result == netLibErrNone {
			session._recvBuffer.aheadReadCursor(readSize)
		} else {
			return result
		}
	}

	return netLibErrNone
}

func _checkReceiveError(sessionIndex int, err error) sessionError {
	if err == nil {
		return netLibErrNone
	}

	if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
		return sessionCloseSocketReadTimeout
	}

	if nErr, ok := err.(net.Error); ok && nErr.Temporary() {
		return netLibErrNone
	}

	scommon.LogError(fmt.Sprintf("[_checkReceiveError] Tcp Read. Session(%d0, %v", sessionIndex, err))
	return sessionCloseSocketError
}

func (session *tcpSession) _checkReadAfterEnablePacket(readAbleByte int16, packetHeaderSize int16) sessionError {
	// 헤더 크기 보다 작은 사이즈를 보낸 경우는 문제 있는 클라이언트. 연결을 짜르도록 한다
	if readAbleByte < packetHeaderSize {
		scommon.LogDebug(fmt.Sprintf("[handleTCPReceive] SESSION_CLOSE_CASE_DATA_SMALLER_THAN_HEADER_SIZE. Session( %d ), headerSizeROnly(%d)", session.getIndex(), packetHeaderSize))
		return sessionCloseRecvTooSmallData
	}

	if session.isEnablePacketProcess() == false {
		scommon.LogDebug("[handleTCPReceive] _tcpConn.Read - isEnablePacketProcess")
		return sessionDisablePacketProcess
	}

	return netLibErrNone
}

func (session *tcpSession) _makePacketAndCallRecvEvent(readAbleByte int,
												dataBuff []byte,
												maxPacketSizeRd int,
												pktRecvfunc PacketReceivceFunctors) (int, sessionError) {
	sessionIndex := session.getIndex()
	headerSizeROnly := int(pktRecvfunc.PacketHeaderSize)
	readPos := 0

	for {
		if readAbleByte < headerSizeROnly {
			break
		}

		requireDataSize := int(pktRecvfunc.PacketTotalSizeFunc(dataBuff[readPos:]))

		if (requireDataSize < headerSizeROnly) || (requireDataSize > readAbleByte) {
			scommon.LogDebug(fmt.Sprintf("[_makePacketAndCallRecvEvent] _tcpConn.Read - break. Session( %d ), requireDataSize(%d), headerSizeROnly(%d), readAbleByte(%d), readPos(%d)", session.getIndex(), requireDataSize, headerSizeROnly, readAbleByte, headerSizeROnly))
			break
		}

		//한번에 보내기로 한 패킷 보다 많이 보낸 경우. 클라이언트의 요청을 이후 무시한다. 즉 클라이언트가 종료하도록 유도
		if requireDataSize > maxPacketSizeRd {
			scommon.LogDebug(fmt.Sprintf("[_makePacketAndCallRecvEvent] Larger than maximum send Data. Session( %d ), requireDataSize(%d)", session.getIndex(), requireDataSize))
			return readPos, sessionCloseRecvMakePacketTooLargePacketSize
		}

		ltvPacket := dataBuff[readPos:(readPos + requireDataSize)]
		readPos += requireDataSize
		readAbleByte -= requireDataSize

		pktRecvfunc.AddNetMsgOnReceiveFunc(sessionIndex, ltvPacket)
	}

	return readPos, netLibErrNone
}


func (session *tcpSession) _realSendData(sendData []byte) bool {
	conn := session.getSocket()
	chunkSize := int(1024)
	sendDataLen := len(sendData)

	for i := 0; i < sendDataLen; i += chunkSize {
		end := i + chunkSize
		var splitSendData []byte
		if end >= sendDataLen {
			splitSendData = sendData[i:]
		} else {
			splitSendData = sendData[i:end]
		}

		_, err := conn.Write(splitSendData)
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				return false
			}

			if nErr, ok := err.(net.Error); ok && nErr.Temporary() {
				i -= chunkSize
				runtime.Gosched() // 잠깐 쉬는 것이 좋을 것 같으니 다른 고루틴에 양보한다.
				continue
			}

			return false
		}
	}

	return true
}