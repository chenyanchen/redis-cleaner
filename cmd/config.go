package cmd

import "time"

type Config struct {
	Cleaner []*CleanerConfig `validate:"required,dive"`
}

type CleanerConfig struct {
	// Scanner is used to scan redis keys.
	// Scanner is required.
	Scanner *RedisConfig `validate:"required"`

	// Cleaner is used to delete redis keys.
	// Cleaner is optional. If not set, Scanner will be used.
	Cleaner *RedisConfig

	// Match is the pattern of keys to scan.
	// Match is required.
	Match string `validate:"required"`

	// Count is the number of keys to scan per time.
	// Count is optional. Default is 10.
	// See: https://redis.io/commands/scan
	Count int64

	// Interval time of per scan.
	// Interval is optional. If not set, no sleep between scans.
	Interval time.Duration
}

type RedisConfig struct {
	Addr     string
	Username string
	Password string
}
