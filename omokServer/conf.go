package omokServer

type OmokConf struct {
	// 네트워크쪽
	Network              string // tcp4(ipv4 only), tcp6(ipv6 only), tcp(ipv4, ipv6)
	IPAddress          string // 127.0.0.1
	MaxSessionCount      int    // 최대 클라이언트 세션 수. 넉넉하게 많이 해도 괜찮다
	MaxPacketSize        int    // 최대 패킷 크기
	RecvPacketRingBufferMaxSize int
	MaxNetMsgChanBufferCount     int // 네트워크 이벤트 메시지 채널 버퍼의 최대 크기


	// 애플리케이션쪽
	RedisResChanCapacity         int
	UserSendBufferSzie           int
	HeartbeatReqIntervalTimeMSec int // 몇 밀리세컨드 간격으로 클라에게 허트비트를 보낼지
	HeartbeatWaitTimeMSec        int // 허트비트 보낸 후 답변을 기다리는 최대 대기 시간. 밀리세컨드
}
