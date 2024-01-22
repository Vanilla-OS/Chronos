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
	"fmt"

	"github.com/dgraph-io/ristretto"
	"github.com/eko/gocache/lib/v4/cache"
	ristretto_store "github.com/eko/gocache/store/ristretto/v4"
)

var (
	cacheManager *cache.Cache[[]byte]
)

func prepareCache() {
	ristrettoCache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1000,
		MaxCost:     100,
		BufferItems: 64,
	})
	if err != nil {
		fmt.Println("Error while creating cache manager")
		panic(err)
	}
	ristrettoStore := ristretto_store.NewRistretto(ristrettoCache)
	cacheManager = cache.New[[]byte](ristrettoStore)
}
