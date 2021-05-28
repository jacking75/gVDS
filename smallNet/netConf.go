package smallNet


type NetworkConfig struct {
	Network              string // tcp4(ipv4 only), tcp6(ipv6 only), tcp(ipv4, ipv6)
	BindAddress          string // 만약 IP와 포트번호 결합이면 localhost:19999
	MaxSessionCount      int    // 최대 클라이언트 세션 수. 넉넉하게 많이 해도 괜찮다
	MaxPacketSize        int    // 최대 패킷 크기

	RecvPacketRingBufferMaxSize int

	MaxNetMsgChanBufferCount     int // 네트워크 이벤트 메시지 채널 버퍼의 최대 크기


	// golang은 기본은 nagle's 알고리즘은 off 시킨다.  https://golang.org/pkg/net/#TCPConn.SetNoDelay
	IsNoDelay            bool
	SockReadbuf          int    // 소켓 버퍼 크기-받기
	SockWritebuf         int    // 소켓 버퍼 크기-보내기
}