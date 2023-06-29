package cmd

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// New returns a new cobra command.
func New() *cobra.Command {
	h := &handler{}
	c := &cobra.Command{
		Use:     "cleaner",
		Short:   "Redis cleaner",
		Long:    `Clean redis keys by pattern.`,
		Example: "redis-cleaner --addr=127.0.0.1:6379 --match=log:*:*:* --count=1024 --interval=1s",
		Args:    cobra.NoArgs,
		Version: "0.1.0",
		Run:     h.Run,
	}
	c.Flags().StringVar(&h.addr, "addr", "", "Redis addr")
	// log:52b7f7a2-5f92-49a3-aae1-b96fdd2207c5:Article:221487313
	c.Flags().StringVar(&h.match, "match", "log:*:*:*", "Redis key pattern")
	c.Flags().Int64Var(&h.count, "count", 1024, "Redis scan count")
	c.Flags().DurationVar(&h.interval, "interval", 0, "Redis scan interval")
	return c
}

type handler struct {
	// Redis addr.
	addr string

	// Redis key pattern.
	match string

	// Redis scan count per time.
	count int64

	// interval time of per scan.
	interval time.Duration
}

func (h *handler) Run(_ *cobra.Command, _ []string) {
	// new context.
	ctx := context.Background()

	// new redis client.
	cli := redis.NewClient(&redis.Options{
		Addr: h.addr,
	})

	// ping redis.
	if err := cli.Ping(ctx).Err(); err != nil {
		log.Fatal().Err(err).Msg("ping redis failed")
	}

	var cursor uint64
	for {
		// scan keys
		keys, nextCursor, err := cli.Scan(ctx, cursor, h.match, h.count).Result()
		if err != nil {
			log.Err(err).Str("match", h.match).Uint64("cursor", cursor).Msg("scan redis failed")
		}

		// delete keys
		for _, key := range keys {
			dur, err := cli.TTL(ctx, key).Result()
			if err != nil {
				log.Err(err).Str("key", key).Msg("ttl redis key failed")
				continue
			}

			// condition: ttl <= 0
			if dur > 0 {
				continue
			}
			if err := cli.Del(ctx, key).Err(); err != nil {
				log.Err(err).Str("key", key).Msg("del redis key failed")
			}
		}

		if nextCursor == 0 {
			break
		}
		cursor = nextCursor

		// sleep interval
		time.Sleep(h.interval)
	}
}
