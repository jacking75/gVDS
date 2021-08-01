package redisDB

import (
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

type Conf struct {
	Address string
	PoolSize int
	ReqTaskChanCapacity int
}

type Client struct {
	_rc *redis.Client

	ReqTaskChan chan ReqTask

	_workerCount int
	_onDone      chan struct{}

	_funcTaskIDList []int16
	_funcList       []func(ReqTask)
}

func (c *Client) Init(conf Conf) bool {
	c.connect(conf.Address, conf.PoolSize)

	c._workerCount = conf.PoolSize
	c.ReqTaskChan = make(chan ReqTask, conf.ReqTaskChanCapacity)
	c._onDone = make(chan struct{})

	c.settingPacketFunction()

	return true
}

func (c *Client) Start() {
	for i := 0; i < c._workerCount; i++ {
		go c.worker_goroutine()
	}
}

func (c *Client) Stop() {
	fmt.Println("[redisDB.Stop] Start")
	close(c._onDone)
	time.Sleep(time.Millisecond * 100)
	fmt.Println("[redisDB.Stop] End")
}


func (c *Client) worker_goroutine() {
	fmt.Println("[worker_goroutine] Start")
LOOP_EXIT:
	for {
		select {
		case task := <- c.ReqTaskChan:
			c.processTask(task)
		case <- c._onDone:
			break LOOP_EXIT
		}
	}

	fmt.Println("[worker_goroutine] End")
}

func (c *Client) connect(address string, pool int) {
	c._rc = redis.NewClient(&redis.Options{
		Addr:     address,
		PoolSize: pool,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	/*if err := c._rc.Ping(); err != nil {
		panic(err)
	}*/
}

func (c *Client) getTaskFunc(packetID int16) func(ReqTask) {
	for i, id := range c._funcTaskIDList {
		if id == packetID {
			return c._funcList[i]
		}
	}

	return nil
}

func (c *Client) settingPacketFunction() {
	maxFuncListCount := 16
	c._funcList = make([]func(ReqTask), 0, maxFuncListCount)
	c._funcTaskIDList = make([]int16, 0, maxFuncListCount)

	c.addTaskFunction(TaskID_ReqLogin, c.processTaskReqLogin)
}

func (c *Client) addTaskFunction(packetID int16,	taskFunc func(ReqTask)) {
	c._funcList = append(c._funcList, taskFunc)
	c._funcTaskIDList = append(c._funcTaskIDList, packetID)
}




/*
err := rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

	val2, err := rdb.Get(ctx, "key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}
	// Output: key value
	// key2 does not exist
 */


/*
func _RPush(rdb *redis.Client, battleKey string, playData []byte) {
	if err := rdb.RPush(ctx, battleKey, playData).Err(); err != nil {
		fmt.Println("[fail] _RPush: ", err)
	}
}

func _LRange(rdb *redis.Client, battleKey string, start int64, last int64) {
	if values, err := rdb.LRange(ctx, battleKey, start, last).Result(); err != nil {
		fmt.Println("[fail] _LRange: ", err)
	} else {
		for _, value := range values {
			fmt.Println(value)
		}
	}
}
 */