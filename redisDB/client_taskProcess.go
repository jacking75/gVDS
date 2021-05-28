package redisDB

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"strconv"
)

func (c *Client) processTask(task ReqTask) {
	if taskFunc := c.getTaskFunc(task.ID); taskFunc != nil {
		taskFunc(task)
	} else {
		fmt.Println("[processTask] invalid task id. ", task.ID)
	}
}

func (c *Client) processTaskReqLogin(reqTask ReqTask) {
	var req ReqTaskLogin
	buf := bytes.NewBuffer(reqTask.Data)
	if err := gob.NewDecoder(buf).Decode(&req); err  != nil {
		fmt.Println("[processTaskReqLogin] err ", err)
	}


	var res ResTask
	res.UID = reqTask.UID
	res.ID = TaskID_ResLogin
	res.Result = taskResult_None

	key := redisKey_UserAuth(req.UserID)
	val, err := c.rc.Get(key).Result()
	if err != nil {
		res.Result = taskResult_EmptyAuth
	} else {
		if auth, _ := strconv.ParseUint(val, 10, 64); auth != req.AuthCode {
			res.Result = taskResult_FailAuth
		}
	}

	reqTask.ResChan <- res
}

/*
package main

import (
    "bytes"
    "encoding/gob"
    "fmt"
)

type Hoge struct {
    F1 string
    F2 int64
}

func main() {
    encoded := encode()
    fmt.Printf("encoded: %d bytes\n", len(encoded))

    decoded := decode(encoded)
    fmt.Printf("decoded: %+v\n", decoded)
}

func encode() []byte {
    h := Hoge{F1: "hoge", F2: 123}
    buf := bytes.NewBuffer(nil)
    _ = gob.NewEncoder(buf).Encode(&h)
    return buf.Bytes()
}

func decode(data []byte) *Hoge {
    var h Hoge
    buf := bytes.NewBuffer(data)
    _ = gob.NewDecoder(buf).Decode(&h)
    return &h
}
 */