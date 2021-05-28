module omokServer

go 1.16

require (
	redisDB v0.0.1
	scommon v0.0.1
	smallNet v0.0.1
)

replace scommon v0.0.1 => ../scommon

replace smallNet v0.0.1 => ../smallNet

replace redisDB v0.0.1 => ../redisDB
