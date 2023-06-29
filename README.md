# redis-cleaner

Redis cleaner.

## Known issues

1. Ping Redis failed and return `EOF` error, use `github.com/go-redis/redis/v8` instead
   of `github.com/redis/go-redis/v9`.
