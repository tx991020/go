package main

import (
	"flag"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

var (
	pool          *redis.Pool
	redisServer   = flag.String("redisServer", "47.93.30.165:6379", "")
	redisPassword = flag.String("redisPassword", "111111", "")
)

//初始化一个pool
func newPool(server, password string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		MaxActive:   5,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if _, err := c.Do("AUTH", password); err != nil {
				c.Close()
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}
func main() {

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		flag.Parse()
		pool = newPool(*redisServer, *redisPassword)

		conn := pool.Get()
		defer conn.Close()
		v, err := redis.String(conn.Do("GET", "adsl"))
		if err == nil {
			// If no results send null

			c.String(http.StatusOK, "%s", v)
		} else {

			c.String(http.StatusOK, "%s", "60.169.216.244:8888")
		}

	})

	r.Run(":8000")

}
