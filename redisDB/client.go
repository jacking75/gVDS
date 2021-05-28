package redisDB

import (
	"fmt"
	"github.com/go-redis/redis"
)

type Conf struct {
	Address string
	PoolSize int
	ReqTaskChanCapacity int
	ResTaskChanCapacity int
}

type Client struct {
	rc *redis.Client

	ReqTaskChan chan ReqTask

	workerCount int
	onDone chan struct{}

	funcTaskIDList []int16
	funcList []func(ReqTask)
}

func (c *Client) Init(conf Conf) bool {
	c.connect(conf.Address, conf.PoolSize)

	c.workerCount = conf.PoolSize
	c.ReqTaskChan = make(chan ReqTask, conf.ReqTaskChanCapacity)
	c.onDone = make(chan struct{})

	c.settingPacketFunction()

	return true
}

func (c *Client) Start() {
	for i := 0; i < c.workerCount; i++ {
		go c.worker_goroutine()
	}
}

func (c *Client) Stop() {
	fmt.Println("[redisDB.Stop] Start")
	close(c.onDone)
	fmt.Println("[redisDB.Stop] End")
}


func (c *Client) worker_goroutine() {
	fmt.Println("[worker_goroutine] Start")
LOOP_EXIT:
	for {
		select {
		case task := <- c.ReqTaskChan:
			c.processTask(task)
		case <- c.onDone:
			break LOOP_EXIT
		}
	}

	fmt.Println("[worker_goroutine] End")
}

func (c *Client) connect(address string, pool int) {
	c.rc = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		PoolSize: pool,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	if err := c.rc.Ping(); err != nil {
		panic(err)
	}
}

func (c *Client) getTaskFunc(packetID int16) func(ReqTask) {
	for i, id := range c.funcTaskIDList {
		if id == packetID {
			return c.funcList[i]
		}
	}

	return nil
}

func (c *Client) settingPacketFunction() {
	maxFuncListCount := 16
	c.funcList = make([]func(ReqTask), 0, maxFuncListCount)
	c.funcTaskIDList = make([]int16, 0, maxFuncListCount)

	c.addTaskFunction(TaskID_ReqLogin, c.processTaskReqLogin)
}

func (c *Client) addTaskFunction(packetID int16,	taskFunc func(ReqTask)) {
	c.funcList = append(c.funcList, taskFunc)
	c.funcTaskIDList = append(c.funcTaskIDList, packetID)
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