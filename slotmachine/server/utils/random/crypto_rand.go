package random

import (
	"crypto/hmac"
	crand "crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"math/big"
	"os"
	"time"
)

var (
	_rand_hash = new(big.Int).SetUint64(uint64(time.Now().UTC().UnixNano()^int64(os.Getppid())) << 37 & 0xBF58476D1CE4E5B9)
)

// CryptoRandInt - crypto/rand
//
//	接近真實隨機數，但是速度較慢，用於做種子生成
func crypto_random() *big.Int {
	sd, _ := crand.Int(crand.Reader, _rand_hash)
	return sd
}

// CryptoRand - 參考 Stake.com 的 Fairness 實作，使用 HNMAC-SHA256 產生隨機十六進位字串
//
//	@param  - the server seed.
//	@param  - the client seed.
//	@@return - the random number.
func CryptoRand(server_seed, client_seed string) []byte {
	h := hmac.New(sha256.New, []byte(server_seed))
	h.Write([]byte(client_seed))

	return h.Sum(nil)
}

func CryptoRandString(server_seed, client_seed string) string {
	h := hmac.New(sha256.New, []byte(server_seed))
	h.Write([]byte(client_seed))
	sha := hex.EncodeToString(h.Sum(nil))
	return sha
}
