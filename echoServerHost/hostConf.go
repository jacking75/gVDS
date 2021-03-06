package main

import (
	"flag"
)

type hostConf struct {
	maxGameCount int // 논리적인 게임 서버의 수
	startTcpPort int // 게임 서버에서 사용할 포트 번호의 시작 수. 순차적으로 1씩 증가시킨다.

	conf echoConf
}

//-c_maxGameCount=1024 -c_startTcpPort=11021 -c_network=tcp4 -c_ipAddress=127.0.0.1 -c_maxSessionCount=4 -c_maxPacketSize=1024 -c_recvPacketRingBufferMaxSize=8012 -c_maxNetMsgChanBufferCount=128
func createHostConf() hostConf {
	var hConf hostConf

	flag.IntVar(&hConf.maxGameCount, "c_maxGameCount", 1, "int flag")
	flag.IntVar(&hConf.startTcpPort, "c_startTcpPort", 11021, "int flag")

	flag.StringVar(&hConf.conf.Network, "c_network", "tcp4", "string flag")
	flag.StringVar(&hConf.conf.IPAddress, "c_ipAddress", "127.0.0.1", "string flag")

	flag.IntVar(&hConf.conf.MaxSessionCount, "c_maxSessionCount", 4, "int flag")
	flag.IntVar(&hConf.conf.MaxPacketSize, "c_maxPacketSize", 1024, "int flag")
	flag.IntVar(&hConf.conf.RecvPacketRingBufferMaxSize, "c_recvPacketRingBufferMaxSize", 1024 * 16, "int flag")
	flag.IntVar(&hConf.conf.MaxNetMsgChanBufferCount, "c_maxNetMsgChanBufferCount", 128, "int flag")

	flag.Parse()

	return hConf
}