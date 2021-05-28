// 애플리케이션에서 네트워크 라이브러리에 접근할 함수는 모두 여기에만 정의한다.
package smallNet

import "net"

type ServrNet struct {
	_IsEnd bool

	_mClientListener *net.TCPListener
	_tcpSessionManager *tcpClientSessionManager
	_tcpServerSessionManager *tcpServerSessionManager
	_netConfig NetworkConfig
	_NetMsgChan chan NetMsg

	ISendToClient func(int, []byte) bool
	ISendToAllClient func([]byte)
}

// 네트워크 시작
func StartNetwork(clientConfig NetworkConfig,
					packetHeaderSize int16,
					packetTotalSizeFunc func([]byte) int16) *ServrNet {

	return startNetwork_impl(clientConfig, packetHeaderSize, packetTotalSizeFunc)
}

func (snet *ServrNet) Stop() {
	snet.stop_impl()
}

func (snet *ServrNet) GetNetMsg() <-chan NetMsg {
	return snet._NetMsgChan
}

func (snet *ServrNet) PrepareNetMsg(msg NetMsg) NetMsg {
	return snet.prepareNetMsg_impl(msg)
}

// 지정한 클라이언트를 강제 종료 시킨다
func (snet *ServrNet) ForceDisconnectClient(sessionIndex int) {
	snet._tcpSessionManager.forceDisconnectClient(sessionIndex)
}

// 지정한 클라이언트의 패킷 처리를 중단 시킨다.
func (snet *ServrNet) DisablePacketProcessClient(sessionIndex int) {
	snet._tcpSessionManager.disablePacketProcessClient(sessionIndex)
}










// Remote Network
/*
///<<< 서버 세션 관련 네트워크 API
func (snet *ServrNet) InitServerSessionManager(serverConfig NetworkConfig,
								packetHeaderSize int16,
								packetTotalSizeFunc func([]byte) int16) {
	_tcpServerSessionManager = newServerSessionManager(serverConfig,
													packetHeaderSize,
													packetTotalSizeFunc)
}

func (snet *ServrNet) stopServerSessionManager() {
	snet.stopServerSessionManager_impl()
}

func ConnectAndReceiveToRemoteServer(serverName string, remoteAddress string) {
	_tcpServerSessionManager.connectAndReceive(serverName, remoteAddress)
}

func ForceDisconnectServerSession(sessionIndex int) {
	_tcpServerSessionManager.forceDisconnectServerSession(sessionIndex)
}

func StopAllServerSession() {
	_tcpServerSessionManager.stop()
}

// Send Interface Function
var ISendPacketToServer func(int, []byte) bool
///>>> 서버 세션
*/