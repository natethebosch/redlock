package redlock

import (
	"context"
	"errors"
	"github.com/go-redis/redis"
	"math/rand"
	"strconv"
	"time"
)

const retryDelay = time.Millisecond

type RedLock struct {
	random int
	keys   []string
}

// Captures a distributed lock at the key provided.
// Duration specifies how long the lock will be held, the actual lock duration will be less than the value provided
// because of the delay in locking each node in the cluster.
// This method will block until the lock is captured.
func (r *RedLock) Lock(db *redis.ClusterClient, key string, duration time.Duration) {

	for {
		err := r.lock(db, key, duration)
		if err != nil {
			<-time.After(retryDelay)
			continue
		}

		return
	}
}

// Captures a distributed lock at the key provided using a context.
// Duration specifies how long the lock will be held, the actual lock duration will be less than the value provided
// because of the delay in locking each node in the cluster.
// Errors will only originate from the context.
func (r *RedLock) LockWithContext(ctx context.Context, db *redis.ClusterClient, key string, duration time.Duration) error {

	for {
		err := r.lock(db, key, duration)

		if err != nil {

			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(retryDelay):
				continue
			}

		}

		return nil
	}
}

func (r *RedLock) lock(db *redis.ClusterClient, key string, duration time.Duration) error {

	r.random = rand.Int()

	slots, err := db.ClusterSlots().Result()
	if err != nil {
		return err
	}

	r.keys = r.keys[0:0]

	successCount := 0

	for _, slot := range slots {

		slotKey, err := keyWithinSlotRange(key, slot.Start, slot.End)
		if err != nil {
			r.Unlock(db)
			return err
		}

		ok, err := db.SetNX(slotKey, r.random, duration).Result()
		if ok && err == nil {
			successCount++
			r.keys = append(r.keys, slotKey)
		}
	}

	if successCount <= (len(slots) / 2) {
		r.Unlock(db)
		return errors.New("failed to capture lock")
	}

	return nil
}

// Releases
func (r *RedLock) Unlock(db *redis.ClusterClient) {

	for _, key := range r.keys {

		// ignore errors
		db.Eval(`if redis.call("get",KEYS[1]) == ARGV[1] then
			    return redis.call("del",KEYS[1])
			else
			    return 0
			end`, []string{key}, r.random)

	}

	r.keys = r.keys[0:0]
}

var unableToComputeKey = errors.New("unable to compute key within slot range in a reasonable number of iterations")

func keyWithinSlotRange(baseKey string, min, max int) (string, error) {

	var key string
	var slot int

	for i := 0; i < 10000; i++ {
		key = "{" + strconv.Itoa(i) + "}" + baseKey
		slot = ComputeHashSlot(key)

		if slot >= min && slot <= max {
			return key, nil
		}
	}

	return "", unableToComputeKey
}
