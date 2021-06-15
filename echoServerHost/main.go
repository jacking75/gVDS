package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
)



func main() {
	ctx, done := signal.NotifyContext(context.Background(), os.Interrupt)
	defer done()

	hConf := createHostConf()

	echoServerList := make([]*echoServer, hConf.maxGameCount)
	for i := 0; i < hConf.maxGameCount; i++ {
		svr := new(echoServer)
		svr.Init(hConf.startTcpPort + i, hConf.conf)
		svr.StartServer()

		echoServerList[i] = svr
	}

	go printCountPerSec_goroutine()


	fmt.Println("Waiting for SIGINT.")
	<-ctx.Done()


	for i := 0; i < hConf.maxGameCount; i++ {
		echoServerList[i].Stop()
	}


	fmt.Println("END")
}

