package md5hash

import (
	"crypto/md5"
	"encoding/hex"
)

/*
MD5 is a hash function, not an encryption function. The point of a hash function
is that it is impossible to convert the output back into the input.
*/

func Hash32bit(value string) string {
	h := md5.New()
	h.Write([]byte(value))
	hex.EncodeToString(h.Sum(nil))

	return hex.EncodeToString(h.Sum(nil))
}

func Hash16bit(value string) string {
	return Hash32bit(value)[8:24]
}

func Hash8bit(value string) string {
	return Hash32bit(value)[12:20]
}
