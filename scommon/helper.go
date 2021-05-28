package scommon

import "time"



func CurrnetUnixTime() int64 {
	return time.Now().Unix()
}

func CurrentUnixMillSec() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func NextUnixMillSec(addTime int64) int64 {
	curTime := CurrentUnixMillSec()
	return (curTime + addTime)
}

func StringSliceToBytes(beginIndx int, endIndex int, stringSlice []string, outBufferMaxSize int, outBuffer []byte) (writeOutBuffer []byte, netxIndex int) {
	writePos := 0
	netxIndex = beginIndx

	for i := beginIndx; i < endIndex; i++ {
		buffer := []byte(stringSlice[i])
		bufferSize := len(buffer)

		if (writePos + bufferSize) >= outBufferMaxSize {
			break
		}

		copy(outBuffer[writePos:], buffer)
		writePos += bufferSize
		netxIndex++
	}

	writeOutBuffer = outBuffer[0:writePos]
	return writeOutBuffer, netxIndex
}




