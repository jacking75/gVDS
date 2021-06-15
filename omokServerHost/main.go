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

	conf := createHostConf()

	redisC := createRedis(conf)
	redisC.Start()

	omokServerList := make([]*omokServer.Server, conf.maxGameCount)
	for i := 0; i < conf.maxGameCount; i++ {
		omokServ := new(omokServer.Server)
		omokServ.Init(conf.startTcpPort + i, conf.omokConf, redisC.ReqTaskChan)
		omokServ.StartServer()

		omokServerList[i] = omokServ
	}

	fmt.Println("Waiting for SIGINT.")
	<-ctx.Done()


	for i := 0; i < conf.maxGameCount; i++ {
		omokServerList[i].Stop()
	}

	redisC.Stop()

	fmt.Println("END")
}

func createRedis(hconf hostConf) *redisDB.Client {
	conf := redisDB.Conf{
		Address: hconf.RedisAddress,
		PoolSize: hconf.RedisPoolSize,
		ReqTaskChanCapacity: hconf.RedisReqTaskChanCapacity,
		ResTaskChanCapacity: hconf.RedisResTaskChanCapacity,
	}

	client := new(redisDB.Client)
	client.Init(conf)
	return client
}