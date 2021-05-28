package redisDB

type ReqTask struct {
	ResChan chan ResTask
	UID uint64
	ID int16
	Data []byte
}

type ReqTaskLogin struct {
	UserID string
	AuthCode uint64
}












