package cache

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	Pool *redis.Pool
)

// init() function is just like the main function, does not take any argument nor return anything.
// This function is present in every package and this function is called when the package is initialized.
// This function is declared implicitly,
// so you cannot reference it from anywhere
// and you are allowed to create multiple init() function in the same program and they execute in the order they are created.
func init() {
	redisHost := ":6379"
	Pool = newPool(redisHost)
	close()
}

func newPool(server string) *redis.Pool {

	return &redis.Pool{

		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,

		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			return c, err
		},

		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func close() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, syscall.SIGKILL)
	go func() {
		<-c
		Pool.Close()
		os.Exit(0)
	}()
}

func Get(key string) ([]byte, error) {

	conn := Pool.Get()
	defer conn.Close()

	var data []byte
	data, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return data, fmt.Errorf("error get key %s: %v", key, err)
	}
	return data, err
}

func Put(key string, value interface{}) error {
	conn := Pool.Get()
	defer conn.Close()
	data, err := conn.Do("SET", key, value)
	if err != nil {
		return err
	}
	_ = data
	return nil
}

func PutWithExp(key string, value interface{}, expirationTime int64) error {
	conn := Pool.Get()
	defer conn.Close()
	data, err := conn.Do("SETEX", key, expirationTime, value)
	if err != nil {
		return err
	}
	_ = data
	return nil
}
