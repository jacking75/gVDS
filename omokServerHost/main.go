package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"redisDB"

	"omokServer"
)

func main() {
	ctx, done := signal.NotifyContext(context.Background(), os.Interrupt)
	defer done()

	conf := createConfig()

	redisC := createRedis()
	redisC.Start()

	omokServ := omokServer.Server{}
	omokServ.Init(conf, redisC.ReqTaskChan)
	omokServ.StartServer()

	fmt.Println("Waiting for SIGINT.")
	<-ctx.Done()

	fmt.Println("We are done here.")
	omokServ.Stop()
	redisC.Stop()

	fmt.Println("END")
}

func createConfig() omokServer.OmokConf {
	conf := omokServer.OmokConf{}
	conf.Network = "tcp4"
	conf.BindAddress = ":11021"
	conf.MaxSessionCount = 32
	conf.MaxPacketSize = 1024
	conf.RecvPacketRingBufferMaxSize = 1024 * 16
	conf.SendPacketRingBufferMaxSize = 1024 * 16
	conf.MaxNetMsgChanBufferCount = 128

	return conf
}

func createRedis() *redisDB.Client {
	conf := redisDB.Conf{
		Address: "127.0.0.1:6379",
		PoolSize: 8,
		ReqTaskChanCapacity: 128,
		ResTaskChanCapacity: 128,
	}

	client := new(redisDB.Client)
	client.Init(conf)
	return client
}