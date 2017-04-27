package redlock

import (
	"github.com/go-redis/redis"
	"testing"
)

func TestComputeHashSlot_TrivialCase(t *testing.T) {

	sampleKey := "mykey"

	calc := ComputeHashSlot(sampleKey)

	cli := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{"localhost:7000", "localhost:7001", "localhost:7002"},
	})

	expected, err := cli.ClusterKeySlot(sampleKey).Result()
	if err != nil {
		t.Fatal(err)
	}

	if int64(calc) != expected {
		t.Fatal("expected", expected, "got", calc)
	}
}

func TestComputeHashSlot_ValidBraces(t *testing.T) {

	sampleKey := "my{key}"

	calc := ComputeHashSlot(sampleKey)

	cli := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{"localhost:7000", "localhost:7001", "localhost:7002"},
	})

	expected, err := cli.ClusterKeySlot(sampleKey).Result()
	if err != nil {
		t.Fatal(err)
	}

	if int64(calc) != expected {
		t.Fatal("expected", expected, "got", calc)
	}
}

func TestComputeHashSlot_NotSpecialBraces(t *testing.T) {

	sampleKey := "my{}"

	calc := ComputeHashSlot(sampleKey)

	cli := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{"localhost:7000", "localhost:7001", "localhost:7002"},
	})

	expected, err := cli.ClusterKeySlot(sampleKey).Result()
	if err != nil {
		t.Fatal(err)
	}

	if int64(calc) != expected {
		t.Fatal("expected", expected, "got", calc)
	}
}

func TestComputeHashSlot_InvalidBraces(t *testing.T) {

	sampleKey := "my}key{"

	calc := ComputeHashSlot(sampleKey)

	cli := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{"localhost:7000", "localhost:7001", "localhost:7002"},
	})

	expected, err := cli.ClusterKeySlot(sampleKey).Result()
	if err != nil {
		t.Fatal(err)
	}

	if int64(calc) != expected {
		t.Fatal("expected", expected, "got", calc)
	}
}

func BenchmarkComputeHashSlot(b *testing.B) {

	sampleKey := "my}key{"

	for i := 0; i < b.N; i++ {
		_ = ComputeHashSlot(sampleKey)
	}
}
