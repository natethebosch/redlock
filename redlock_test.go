package redlock

import (
	"github.com/go-redis/redis"
	"testing"
	"time"
)

func TestRedLock_Lock(t *testing.T) {

	db := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{"localhost:7000", "localhost:7001", "localhost:7002"},
	})

	done := make(chan bool)

	go func() {
		lock := &RedLock{}
		lock.Lock(db, "mykey", time.Second)
		time.Sleep(500 * time.Millisecond)
		lock.Unlock(db, "mykey")

		done <- true
	}()

	go func() {
		lock := &RedLock{}
		lock.Lock(db, "mykey", time.Second)
		time.Sleep(500 * time.Millisecond)
		lock.Unlock(db, "mykey")

		done <- true
	}()

	start := time.Now()

	<-done
	<-done

	duration := time.Now().Sub(start)

	if duration < 1*time.Second {
		t.Fatal("locks aren't working, should take at least 1 second")
	}

	if duration >= 2*time.Second {
		t.Fatal("locks aren't unlocking, shouldn't take 2 seconds")
	}

}

func TestRedLock_LockUnLock(t *testing.T) {

	db := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{"localhost:7000", "localhost:7001", "localhost:7002"},
	})

	start := time.Now()

	lock := &RedLock{}
	lock.Lock(db, "mykey", time.Second)
	lock.Unlock(db, "mykey")

	lock.Lock(db, "mykey", time.Second)
	lock.Unlock(db, "mykey")

	if time.Now().Sub(start).Seconds() > 0 {
		t.Fatal("should be less than a second else the lock is not being unlocked")
	}
}
