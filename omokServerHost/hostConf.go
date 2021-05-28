package main

import (
	"flag"
	"omokServer"
)

type hostConf struct {
	maxGameCount int
	startTcpPort int

	RedisAddress string
	RedisPoolSize int
	RedisReqTaskChanCapacity int
	RedisResTaskChanCapacity int

	omokConf omokServer.OmokConf
}

//-c_maxGameCount=100 -c_startTcpPort=11021 -c_network=tcp4 -c_ipAddress=127.0.0.1
func createHostConf() hostConf {
	var conf hostConf

	flag.IntVar(&conf.maxGameCount, "c_maxGameCount", 1, "int flag")
	flag.IntVar(&conf.startTcpPort, "c_startTcpPort", 11021, "int flag")

	flag.StringVar(&conf.RedisAddress, "c_redisAddress", "127.0.0.1:6379", "string flag")
	flag.IntVar(&conf.RedisPoolSize, "c_redisPoolSize", 8, "int flag")
	flag.IntVar(&conf.RedisReqTaskChanCapacity, "c_redisReqTaskChanCapacity", 32, "int flag")
	flag.IntVar(&conf.RedisResTaskChanCapacity, "c_redisResTaskChanCapacity", 32, "int flag")

	flag.StringVar(&conf.omokConf.Network, "c_network", "tcp4", "string flag")
	flag.StringVar(&conf.omokConf.IPAddress, "c_ipAddress", "127.0.0.1", "string flag")

	flag.IntVar(&conf.omokConf.MaxSessionCount, "c_maxSessionCount", 4, "int flag")
	flag.IntVar(&conf.omokConf.MaxPacketSize, "c_maxPacketSize", 1024, "int flag")
	flag.IntVar(&conf.omokConf.RecvPacketRingBufferMaxSize, "c_recvPacketRingBufferMaxSize", 1024 * 16, "int flag")
	flag.IntVar(&conf.omokConf.MaxNetMsgChanBufferCount, "c_maxNetMsgChanBufferCount", 128, "int flag")
	flag.IntVar(&conf.omokConf.RedisResChanCapacity, "c_redisResChanCapacity", 32, "int flag")
	flag.IntVar(&conf.omokConf.UserSendBufferSzie, "c_userSendBufferSzie", 1024 * 16, "int flag")
	flag.IntVar(&conf.omokConf.HeartbeatReqIntervalTimeMSec, "c_heartbeatReqIntervalTimeMSec", 3000, "int flag")
	flag.IntVar(&conf.omokConf.HeartbeatWaitTimeMSec, "c_heartbeatWaitTimeMSec", 3000, "int flag")

	flag.Parse()

	return conf
}