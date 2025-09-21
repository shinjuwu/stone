package random_test

import (
	crand "crypto/rand"
	"math/big"
	"os"
	"time"
)

var (
	_point     = 0
	_rand_hash = new(big.Int).SetUint64(uint64(time.Now().UTC().UnixNano()^int64(os.Getppid())) << 37 & 0xBF58476D1CE4E5B9)
)

func GetSeed() int64 {
	sd, _ := crand.Int(crand.Reader, _rand_hash)
	return sd.Int64()
}

// func Test_Test(t *testing.T) {
// 	// t.Logf("%+v\n", (1 << 31))
// 	// num := uint(256 + (256 << 8) + (256 << 16) + (256 << 24))
// 	// t.Logf("%+v\n", num)
// 	// Random pick with weight

// 	weights := []int64{2, 8, 90}
// 	stat := make([]float64, len(weights))
// 	round := 10000000
// 	for i := 0; i < round; i++ {
// 		pick := random.Int64w(weights)
// 		stat[pick]++
// 	}

// 	for _, v := range stat {
// 		fmt.Printf("%.4f \t\n", v/float64(round))
// 	}

// 	fmt.Printf("wc == %.2f \t\n", stat)
// }
