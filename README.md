# Redlock

Redlock is a Golang package that implements the redlock distributed lock algorithm as described in the [redis docs](https://redis.io/topics/distlock). 
This package also implements the [redis cluster hash slot calculation](https://redis.io/topics/cluster-spec).

### Installation
```
go get github.com/natethebosch/redlock

# depenancy
go get github.com/snksoft/crc
```

### API

[gopkg.in docs](https://gopkg.in/natethebosch/redlock.v1)