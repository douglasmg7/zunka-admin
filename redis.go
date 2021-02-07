package main

import (
	"time"

	"github.com/go-redis/redis/v8"
)

// Get.
func redisGet(key string) string {
	val, err := redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return ""
		// log.Printf("Key not exist")
	} else if err != nil {
		checkError(err)
		return ""
	} else {
		return val
		// log.Printf("Val: %+v\n", val)
	}
}

// Set.
func redisSet(key string, val string, exp time.Duration) error {
	err := redisClient.Set(ctx, key, val, exp).Err()
	if err != nil {
		return err
	}
	return nil
}

// Del.
func redisDel(key string) error {
	return redisClient.Del(ctx, key).Err()
}

//****************************************************************************
//	Mercado Livre
//****************************************************************************
// Set ML user code.
func setMLUserCode(code string) {
	key := "ml-user-code"
	// Save for one wekeend.
	_ = redisSet(key, code, time.Hour*720)
}

// Get ML user code.
func getMLUserCode() string {
	key := "ml-user-code"
	return redisGet(key)
}
