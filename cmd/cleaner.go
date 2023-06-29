package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// New returns a new cobra command.
func New() *cobra.Command {
	h := &handler{}
	c := &cobra.Command{
		Use:     "cleaner",
		Short:   "Redis cleaner",
		Long:    `Clean redis keys by pattern. More config see config/redis-cleaner.yaml`,
		Example: "redis-cleaner --config config/redis-cleaner.yaml",
		Args:    cobra.NoArgs,
		Version: "0.2.0",
		Run:     h.Run,
	}
	c.Flags().StringVar(&h.cfgFile, "config", "config/redis-cleaner.yaml", "config file")
	return c
}

type handler struct {
	cfgFile string
}

func (h *handler) Run(_ *cobra.Command, _ []string) {
	if err := h.initConfig(); err != nil {
		log.Fatal().Err(err).Str("config", h.cfgFile).Msg("init config failed")
	}

	var cfg *Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatal().Err(err).Msg("unmarshal config failed")
	}

	if err := validator.New().Struct(cfg); err != nil {
		log.Fatal().Err(err).Msg("validate config failed")
	}

	h.run(cfg)
}

// initConfig reads in config file and ENV variables if set.
func (h *handler) initConfig() error {
	// Use config file from the flag.
	viper.SetConfigFile(h.cfgFile)

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("read config file failed: %w", err)
	}
	return nil
}

func (h *handler) run(cfg *Config) {
	for i, cleanerCfg := range cfg.Cleaner {
		if err := h.cleanOne(context.Background(), cleanerCfg); err != nil {
			log.Error().Err(err).Int("index", i).Any("config", cleanerCfg).Msg("clean failed")
		}
	}
}

func (h *handler) cleanOne(ctx context.Context, cfg *CleanerConfig) error {
	// New redis scanner.
	scanner := redis.NewClient(&redis.Options{
		Addr:     cfg.Scanner.Addr,
		Username: cfg.Scanner.Username,
		Password: cfg.Scanner.Password,
	})

	// New redis cleaner.
	cleaner := scanner
	if cfg.Cleaner != nil {
		cleaner = redis.NewClient(&redis.Options{
			Addr:     cfg.Cleaner.Addr,
			Username: cfg.Cleaner.Username,
			Password: cfg.Cleaner.Password,
		})
	}

	var cursor uint64
	for {
		// scan keys
		keys, nextCursor, err := scanner.Scan(ctx, cursor, cfg.Match, cfg.Count).Result()
		if err != nil {
			return fmt.Errorf("scan redis failed: %w", err)
		}

		// delete keys
		for _, key := range keys {
			dur, err := scanner.TTL(ctx, key).Result()
			if err != nil {
				log.Err(err).Str("key", key).Msg("ttl redis key failed")
				continue
			}

			// condition: ttl <= 0
			if dur > 0 {
				continue
			}
			if err := cleaner.Del(ctx, key).Err(); err != nil {
				log.Err(err).Str("key", key).Msg("del redis key failed")
				continue
			}
			log.Info().Str("key", key).Msg("del redis key success")
		}

		if nextCursor == 0 {
			break
		}
		cursor = nextCursor

		// sleep interval
		time.Sleep(cfg.Interval)
	}
	return nil
}
