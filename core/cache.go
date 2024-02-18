package core

/*	License: GPLv3
	Authors:
		Mirko Brombin <send@mirko.pm>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2024
	Description:
		Chronos is a simple, fast and lightweight documentation server written in Go.
*/

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/dgraph-io/ristretto"
	"github.com/eko/gocache/lib/v4/cache"
	big_cache_backend "github.com/eko/gocache/store/bigcache/v4"
	go_cache_backend "github.com/eko/gocache/store/go_cache/v4"
	ristretto_backend "github.com/eko/gocache/store/ristretto/v4"
	go_cache "github.com/patrickmn/go-cache"
	redis "github.com/redis/go-redis/v9"
	"github.com/vanilla-os/Chronos/settings"
	redis_backend "github.com/vanilla-os/Chronos/utils"
)

var (
	cacheManager *cache.Cache[[]byte]
)

func NewRistrettoCache() (*cache.Cache[[]byte], error) {
	ristrettoCache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1000,
		MaxCost:     100,
		BufferItems: 64,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to create ristretto cache: %w", err)
	}

	ristrettoStore := ristretto_backend.NewRistretto(ristrettoCache)
	cacheManager = cache.New[[]byte](ristrettoStore)

	return cacheManager, nil
}

func NewBigCache() (*cache.Cache[[]byte], error) {
	bigcacheClient, err := bigcache.New(context.Background(), bigcache.DefaultConfig(5*time.Minute))
	if err != nil {
		return nil, fmt.Errorf("unable to create bigcache client: %w", err)
	}

	bigcacheStore := big_cache_backend.NewBigcache(bigcacheClient)

	cacheManager := cache.New[[]byte](bigcacheStore)
	return cacheManager, nil
}

func NewRedisCache() (*cache.Cache[[]byte], error) {
	redisServer := settings.Cnf.RedisCacheServer
	if redisServer == "" {
		log.Println("(cache): No Redis server specified, using default: localhost")
		redisServer = "localhost"
	}

	redisPort := settings.Cnf.RedisCachePort
	if redisPort == "" {
		log.Println("(cache): No Redis port specified, using default: 6379")
		redisPort = "6379"
	}

	redisUsername := settings.Cnf.RedisCacheUsername
	if redisUsername == "" {
		log.Println("(cache): No Redis username specified, using no username")
	}

	redisPassword := settings.Cnf.RedisCachePassword
	if redisPassword == "" {
		log.Println("(cache): No Redis password specified, using no password")
	}

	redisDB := settings.Cnf.RedisCacheDB
	if redisDB == 0 {
		log.Println("(cache): No Redis DB specified, using default: 0")
		redisDB = 0
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisServer, redisPort),
		Username: redisUsername,
		Password: redisPassword,
		DB:       redisDB,
	})

	redisStore := redis_backend.NewRedis(redisClient)
	cacheManager = cache.New[[]byte](redisStore)

	status := redisClient.Ping(context.Background())
	if status.Err() != nil {
		return nil, fmt.Errorf("unable to connect to Redis server: %w", status.Err())
	}
	if status.Val() != "PONG" {
		return nil, fmt.Errorf("unable to connect to Redis server: unexpected response: %s", status.Val())
	}

	return cacheManager, nil
}

func NewGoCache() (*cache.Cache[[]byte], error) {
	gocacheClient := go_cache.New(5*time.Minute, 10*time.Minute)
	gocacheStore := go_cache_backend.NewGoCache(gocacheClient)

	cacheManager := cache.New[[]byte](gocacheStore)
	return cacheManager, nil
}

func InitCacheManager() error {
	var err error
	backend := settings.Cnf.CacheBackend
	if backend == "" {
		backend = "ristretto"
		log.Printf("(cache): No cache backend specified, using default: %s", backend)
	}

	switch backend {
	case "ristretto":
		cacheManager, err = NewRistrettoCache()
	case "bigcache":
		cacheManager, err = NewBigCache()
	case "redis":
		cacheManager, err = NewRedisCache()
	case "gocache":
		cacheManager, err = NewGoCache()
	default:
		err = fmt.Errorf("unknown cache backend: %s", backend)
	}

	if err != nil {
		return fmt.Errorf("unable to initialize cache manager (backend: %s): %w", backend, err)
	}

	log.Printf("(cache): Cache manager initialized with backend: %s", backend)
	return nil
}
