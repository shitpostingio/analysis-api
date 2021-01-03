package structs

// RedisConfiguration represents the Redis configuration.
type RedisConfiguration struct {
	Address             string
	NSFWDatabase        int
	FingerprintDatabase int
}
