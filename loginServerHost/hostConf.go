package main

import (
	"flag"
)

type hostConf struct {
	maxGameCount int // 논리적인 게임 서버의 수
	startTcpPort int // 게임 서버에서 사용할 포트 번호의 시작 수. 순차적으로 1씩 증가시킨다.

	RedisAddress string
	RedisPoolSize int
	RedisReqTaskChanCapacity int

	appConf echoConf
}

//-c_maxGameCount=100 -c_startTcpPort=11021 -c_network=tcp4 -c_ipAddress=127.0.0.1
func createHostConf() hostConf {
	var hConf hostConf

	flag.IntVar(&hConf.maxGameCount, "c_maxGameCount", 1, "int flag")
	flag.IntVar(&hConf.startTcpPort, "c_startTcpPort", 11021, "int flag")

	flag.StringVar(&hConf.RedisAddress, "c_redisAddress", "127.0.0.1:6379", "string flag")
	flag.IntVar(&hConf.RedisPoolSize, "c_redisPoolSize", 8, "int flag")
	flag.IntVar(&hConf.RedisReqTaskChanCapacity, "c_redisReqTaskChanCapacity", 32, "int flag")

	flag.StringVar(&hConf.appConf.Network, "c_network", "tcp4", "string flag")
	flag.StringVar(&hConf.appConf.IPAddress, "c_ipAddress", "127.0.0.1", "string flag")

	flag.IntVar(&hConf.appConf.MaxSessionCount, "c_maxSessionCount", 4, "int flag")
	flag.IntVar(&hConf.appConf.MaxPacketSize, "c_maxPacketSize", 1024, "int flag")
	flag.IntVar(&hConf.appConf.RecvPacketRingBufferMaxSize, "c_recvPacketRingBufferMaxSize", 1024 * 16, "int flag")
	flag.IntVar(&hConf.appConf.MaxNetMsgChanBufferCount, "c_maxNetMsgChanBufferCount", 128, "int flag")
	flag.IntVar(&hConf.appConf.RedisResChanCapacity, "c_redisResChanCapacity", 32, "int flag")

	flag.Parse()

	return hConf
}