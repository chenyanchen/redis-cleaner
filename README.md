# redis-cleaner

Redis cleaner.

## How to use

1. Update your condition in `cmd/cleaner.go`:

   ```go
   // my condition
   for _, key := range keys {
       dur, _ := scanner.TTL(ctx, key).Result()
       // condition: ttl <= 0
       if dur > 0 {
           continue
       }
       cleaner.Del(ctx, key).Err()
   }
   ```

2. Build

   ```bash
   make
   ```

3. Update `config/redis-cleaner.yaml`:

4. Run

   ```bash
   ./redis-cleaner --config config/redis-cleaner.yaml
   ```

## Known issues

1. Ping Redis failed and return `EOF` error, use `github.com/go-redis/redis/v8` instead
   of `github.com/redis/go-redis/v9`.
