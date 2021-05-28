package redisDB

import "fmt"

func redisKey_UserAuth(userID string) string {
	return 	fmt.Sprintf("%s_%s", "auth", userID)
}
