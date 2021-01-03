package cache

import (
	"github.com/go-redis/redis/v7"
)

var (
	fClient *redis.Client
	nClient *redis.Client
)

// NewFingerprintRedisClient creates a redis client for fingerprints.
func NewFingerprintRedisClient(address, password string, db int) (err error) {
	fClient, err = newRedisClient(address, password, db)
	return
}

// NewNSFWRedisClient creates a redis client for NSFW data.
func NewNSFWRedisClient(address, password string, db int) (err error) {
	nClient, err = newRedisClient(address, password, db)
	return
}

func newRedisClient(address, password string, db int) (*redis.Client, error) {

	rClient := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})

	_, err := rClient.Ping().Result()
	return rClient, err

}
