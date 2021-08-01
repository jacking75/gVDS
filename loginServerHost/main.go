package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"redisDB"
)



func main() {
	ctx, done := signal.NotifyContext(context.Background(), os.Interrupt)
	defer done()

	hConf := createHostConf()

	redisC := createRedis(hConf)
	redisC.Start()

	echoServerList := make([]*echoServer, hConf.maxGameCount)
	for i := 0; i < hConf.maxGameCount; i++ {
		svr := new(echoServer)
		svr.Init(hConf.startTcpPort, hConf, redisC.ReqTaskChan)
		svr.StartServer()

		echoServerList[i] = svr
	}

	fmt.Println("Waiting for SIGINT.")
	<-ctx.Done()


	for i := 0; i < hConf.maxGameCount; i++ {
		echoServerList[i].Stop()
	}


	fmt.Println("END")
}

func createRedis(hconf hostConf) *redisDB.Client {
	conf := redisDB.Conf{
		Address: hconf.RedisAddress,
		PoolSize: hconf.RedisPoolSize,
		ReqTaskChanCapacity: hconf.RedisReqTaskChanCapacity,
	}

	client := new(redisDB.Client)
	client.Init(conf)
	return client
}