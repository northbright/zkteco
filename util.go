package zkteco

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
)

// GetRedisConn gets the Redis connection.
func GetRedisConn(redisAddr, redisPassword string) (redis.Conn, error) {
	var err error
	var c redis.Conn
	pongStr := ""

	if c, err = redis.Dial("tcp", redisAddr); err != nil {
		return c, err
	}

	if len(redisPassword) != 0 {
		if _, err = c.Do("AUTH", redisPassword); err != nil {
			return c, err
		}
	}

	if pongStr, err = redis.String(c.Do("PING")); err != nil {
		return c, err
	}

	if pongStr != "PONG" {
		err = fmt.Errorf("Redis PING != PONG(%v)", pongStr)
		return c, err
	}

	return c, nil
}
