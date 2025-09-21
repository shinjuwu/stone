package random

import (
	"sync"

	"github.com/bofry/random"
	"github.com/bofry/random/mt19937"
)

var (
	rng_mt19937 = random.New(mt19937.New())
	randlk      = sync.Mutex{}
)

func init() {
	// rng_mt19937.Seed(time.Now().UTC().UnixNano())
	randlk.Lock()
	init := []uint64{
		crypto_random().Uint64(),
		crypto_random().Uint64(),
		crypto_random().Uint64(),
		crypto_random().Uint64(),
	}
	mt := mt19937.New()
	mt.SeedbyArray(init)

	rng_mt19937 = random.New(mt)
	randlk.Unlock()
}

// Uint63n - 產生 0 ~ n-1 的隨機數
//
//	@param n		隨機數範圍
//	@return uint64	隨機數
func Uint63n(n uint64) uint64 {
	randlk.Lock()
	num := rng_mt19937.Uint64n(n)
	randlk.Unlock()
	return num
}

// Uint63sn - 產生 0 ~ n-1 的隨機數
//
//	@param ns		隨機數範圍
func Uint63sn(n []uint64) []uint64 {
	if n == nil {
		return nil
	}

	randlk.Lock()
	nums := make([]uint64, len(n))
	for i, v := range n {
		nums[i] = rng_mt19937.Uint64n(v)
	}
	randlk.Unlock()

	return nums
}

// Int63n - 產生 0 ~ n-1 的隨機數
//
//	@param n		隨機數範圍
//	@return int64	隨機數
func Int63n(n int64) int64 {
	randlk.Lock()
	num := rng_mt19937.Int63n(n)
	randlk.Unlock()
	return num
}

// Int63sn
//
//	@param ns		隨機數範圍
func Int63sn(n []int64) []int64 {
	if n == nil {
		return nil
	}

	randlk.Lock()
	nums := make([]int64, len(n))
	for i, v := range n {
		nums[i] = rng_mt19937.Int63n(v)
	}
	randlk.Unlock()

	return nums
}

// Intn - 產生 0 ~ n-1 的隨機數
func Intn(n int) int {
	randlk.Lock()
	num := rng_mt19937.Intn(n)
	randlk.Unlock()
	return num
}

// Intsn - 產生 0 ~ n-1 的隨機數
// @param ns		隨機數範圍
func Intsn(n []int) []int {
	if n == nil {
		return nil
	}

	randlk.Lock()
	nums := make([]int, len(n))
	for i, v := range n {
		nums[i] = rng_mt19937.Intn(v)
	}
	randlk.Unlock()

	return nums
}
