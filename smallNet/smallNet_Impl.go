package smallNet

import (
	"fmt"
	"net"
	"os"
	"scommon"
)

func startNetwork_impl(netConfig NetworkConfig,
						packetHeaderSize int16,
						packetTotalSizeFunc func([]byte) int16) *ServrNet {
	serverNet := &ServrNet{}

	serverNet._IsEnd = false
	serverNet._netConfig = netConfig
	serverNet._NetMsgChan = make(chan NetMsg, netConfig.MaxNetMsgChanBufferCount)

	serverNet._initNetworkSendFunction()

	pktRecvFunc := PacketReceivceFunctors {
		AddNetMsgOnCloseFunc: serverNet._AddNetMsgOnClose,
		AddNetMsgOnReceiveFunc: serverNet._AddNetMsgOnReceive,
		PacketTotalSizeFunc: packetTotalSizeFunc,
		PacketHeaderSize: packetHeaderSize,
		IsClientSession: true,
	}

	serverNet._tcpSessionManager = newClientSessionManager(netConfig,pktRecvFunc)
	go serverNet._startTCPServer()

	return serverNet
}

func (snet *ServrNet) Stop_impl() {
	scommon.LogInfo("[ServerNet] Stop begin")

	snet._IsEnd = true

	_ = snet._mClientListener.Close()
	scommon.LogInfo("Stop TCPServer Accept")

	snet._tcpSessionManager.stop()

	snet.stopServerSessionManager()

	scommon.LogInfo("[ServerNet] Stop end")
}

func (snet *ServrNet) stopServerSessionManager() {
	if snet._tcpServerSessionManager == nil {
		return
	}

	scommon.LogInfo("Stop All ServerSession")
	snet._tcpServerSessionManager.stop()
}

func (snet *ServrNet) getNetConfig() NetworkConfig {
	return snet._netConfig
}

func (snet *ServrNet) prepareNetMsg_impl(msg NetMsg) NetMsg {
	if msg.Type == NetMsg_Connect {
		if sessionIndex := snet._tcpSessionManager.newSession(msg.TcpConn); sessionIndex >= 0 {
			msg.SessionIndex = sessionIndex
		} else {
			msg.Type = NetMsg_None
		}
	} else if msg.Type == NetMsg_Close {
		snet._tcpSessionManager.deleteSession(msg.SessionIndex)
	}

	return msg
}

func (snet *ServrNet) _startTCPServer() {
	scommon.LogInfo("[_startTCPServer] - start")

	defer scommon.PrintPanicStack()

	config := snet._netConfig
	snet._mClientListener = _newListener(config.Network, config.BindAddress)

	scommon.LogInfo("tcpServerStart - Accept Wait...")

	for {
		tcpConn, err := snet._mClientListener.AcceptTCP()
		if err != nil {
			if snet._IsEnd {
				break;
			}

			if nErr, ok := err.(net.Error); ok {
				if nErr.Temporary() {
					continue
				}
			}

			scommon.LogError(fmt.Sprintf("[_startTCPServer] accept error %v", err))
			break
		}

		snet._AddNetMsgOnConnect(tcpConn)
	}

	scommon.LogInfo("[_startTCPServer] - end")
}

func _newListener(network string, clientAddress string) *net.TCPListener {
	tcpAddr, err := net.ResolveTCPAddr(network, clientAddress)
	if err != nil {
		os.Exit(-1)
	}

	var listener *net.TCPListener
	listener, err = net.ListenTCP(network, tcpAddr)
	if err != nil {
		os.Exit(-1)
	}

	scommon.LogInfo(fmt.Sprintf("[_newListener] listening on %s", clientAddress))
	return listener
}

// 보내기 함수(선언만 있는. 일종의 인터페이스)에 실제 동작함수를 연결한다
func (snet *ServrNet) _initNetworkSendFunction() {
	snet.ISendToClient = snet._sendToClient
	snet.ISendToAllClient = snet._sendToAllClient
}

func (snet *ServrNet) _sendToClient(sessionIndex int, data []byte) bool {
	return snet._tcpSessionManager.sendPacket(sessionIndex, data)
}

func (snet *ServrNet) _sendToAllClient(sendData []byte) {
	snet._tcpSessionManager.sendPacketAllClient(sendData)
}

func (snet *ServrNet) _sendPacketToServer(sessionIndex int, sendData []byte) bool {
	return snet._tcpServerSessionManager.sendPacket(sessionIndex, sendData)
}


func (snet *ServrNet) _AddNetMsgOnConnect(tcpConn *net.TCPConn) {
	msg := NetMsg{
		Type: NetMsg_Connect,
		TcpConn: tcpConn,
	}

	snet._NetMsgChan <- msg
}

func (snet *ServrNet) _AddNetMsgOnClose(sessionIndex int) {
	if snet._IsEnd {
		return
	}

	msg := NetMsg{
		Type: NetMsg_Close,
		SessionIndex: sessionIndex,
	}

	snet._NetMsgChan <- msg
}

func (snet *ServrNet) _AddNetMsgOnReceive(sessionIndex int, data []byte) {
	msg := NetMsg{
		Type: NetMsg_Receive,
		SessionIndex: sessionIndex,
		Data: data,
	}

	snet._NetMsgChan <- msg
}

