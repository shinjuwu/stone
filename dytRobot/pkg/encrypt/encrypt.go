package encrypt

import (
	"crypto/aes"
	"crypto/md5"
	"crypto/rand"
	"dytRobot/pkg/encrypt/aescbc"
	"dytRobot/pkg/encrypt/base64url"
	"encoding/hex"
	"errors"
	"log"
	"strings"
)

// 加密
func EncryptSaltToken(encryptData string, gameKey string) (string, error) {
	/*
		key: 6b14314760bd9280695a95d38082478b
		data: 0c95c25eb25d666ebb01528556a45cc0
		先 AES256 aes-256-cbc 使用PKCS#7填充 加密,  產出 密文後, 將salt組合拼湊為  密文::iv   在使用 base64url encode.
		encrypt_data:  pCBS7aa1Xkld3UFyQVVoWC6v87eTsGPQQeajU6bm12Tkl2h-hvpb07BQ5WDn_G9HOjoJ6UimMlf45OimUrO84gAT
	*/

	iv := CreateIV()

	encrypted := aescbc.AesIvEncrypt(encryptData, gameKey, iv)

	base64Encrypted := base64url.Encode([]byte(encrypted + "::" + iv))

	return base64Encrypted, nil
}

// 解密
func DecryptSaltToken(encryptData string, aeskey string) (string, error) {
	base64DecodeData, err := base64url.Decode(encryptData)
	if err != nil {
		log.Printf("Decode result error: %v", err)
		return "", err
	}

	trimBase64DecodeData := strings.Split(string(base64DecodeData), "::")

	if len(trimBase64DecodeData) == 2 {
		aesEncodeData := []byte(trimBase64DecodeData[0])
		iv := []byte(trimBase64DecodeData[1])
		key := []byte(aeskey)

		origData := aescbc.AesIvDecrypt(key, iv, aesEncodeData)

		return string(origData), nil
	} else {
		return "", errors.New("wrong number of parameters")
	}
}

// 產生 16 碼隨機碼 (aes.BlockSize = 16)
func CreateIV() string {
	bIV := make([]byte, aes.BlockSize)
	_, _ = rand.Read(bIV)
	return string(bIV)
}

// MD5Hash can't not decrypt
func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
