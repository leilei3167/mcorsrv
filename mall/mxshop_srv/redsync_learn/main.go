package main

import (
	"fmt"
	"sync"

	goredislib "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
)

func main() {
	// Create a pool with go-redis (or redigo) which is the pool redisync will
	// use while communicating with Redis. This can also be any pool that
	// implements the `redis.Pool` interface.
	client := goredislib.NewClient(&goredislib.Options{
		Addr: "localhost:6379",
	})
	pool := goredis.NewPool(client) // or, pool := redigo.NewPool(...)

	// Create an instance of redisync to be used to obtain a mutual exclusion
	// lock.
	rs := redsync.New(pool)

	// Obtain a new mutex by using the same name for all instances wanting the
	// same lock.

	mutexname := "test_mu"
	mutex := rs.NewMutex(mutexname)
	rawNum := 100
	m := map[int]int{}

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i1 int) {
			defer wg.Done()
			if err := mutex.Lock(); err != nil {
				panic(err)
			}

			// fmt.Printf("goroutine:%d 获取到锁!\n", i1)
			rawNum--
			m[i1] = i1 + 1
			fmt.Println("rawNum:", rawNum)

			if ok, err := mutex.Unlock(); !ok || err != nil {
				panic("unlock failed")
			}
		}(i)
	}

	wg.Wait()
	fmt.Println("最终rawNum:", rawNum)
	fmt.Println("最终m:", m)
	// Obtain a lock for our given mutex. After this is successful, no one else
	// can obtain the same lock (the same mutex name) until we unlock it.

	// Do your work that requires the lock.

	// Release the lock so other processes or threads can obtain a lock.
}
