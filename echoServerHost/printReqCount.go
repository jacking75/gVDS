package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

// 초당 처리량 표시
func printCountPerSec_goroutine() {
	_onDoneCountPerSec = make(chan struct{})
	status_ticker := time.NewTicker(time.Second * 1)

	for {
		select {
		case <-_onDoneCountPerSec:
			return
		case <-status_ticker.C:
			curTime := time.Now().Format("2006-01-02 15:04:05")

			count := atomic.LoadUint64(&_countPerSec)
			atomic.StoreUint64(&_countPerSec, 0)

			fmt.Println(fmt.Sprintf("%s,   %d\n", curTime, count))
		}
	}
}

/*func _stopCountPerSec() {
	close(_onDoneCountPerSec)
}*/

func incCount() {
	atomic.AddUint64(&_countPerSec, 1)
}

// private
var _onDoneCountPerSec chan struct{}

var _countPerSec uint64