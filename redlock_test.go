package redlock

import (
	"context"
	"github.com/go-redis/redis"
	"testing"
	"time"
)

func ExampleRedLock() {
	var db *redis.ClusterClient

	myLock := &RedLock{}

	// capture the lock, this might take a while
	myLock.Lock(db, "some-key", 5*time.Second)

	// do protected things

	// dispose of the lock
	myLock.Unlock(db)
}

func ExampleRedLock_LockWithContext() {
	var db *redis.ClusterClient

	myLock := &RedLock{}

	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)

	// capture the lock with timeout
	err := myLock.LockWithContext(ctx, db, "some-key", 5*time.Second)
	if err != nil {
		// handle timeout
		return
	}

	// do protected things

	// dispose of the lock
	myLock.Unlock(db)
}

func TestRedLock_Lock(t *testing.T) {

	db := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{"localhost:7000", "localhost:7001", "localhost:7002"},
	})

	done := make(chan bool)

	go func() {
		lock := &RedLock{}
		lock.Lock(db, "mykey", time.Second)
		time.Sleep(500 * time.Millisecond)
		lock.Unlock(db)

		done <- true
	}()

	go func() {
		lock := &RedLock{}
		lock.Lock(db, "mykey", time.Second)
		time.Sleep(500 * time.Millisecond)
		lock.Unlock(db)

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
	lock.Unlock(db)

	lock.Lock(db, "mykey", time.Second)
	lock.Unlock(db)

	if time.Now().Sub(start).Seconds() > 0 {
		t.Fatal("should be less than a second else the lock is not being unlocked")
	}
}
