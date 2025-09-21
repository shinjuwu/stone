package aescbc

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

var errByte = []byte("")

func pkcs7Padding(ciphertext []byte) []byte {
	padding := aes.BlockSize - len(ciphertext)%aes.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
func pkcs7UnPadding(origData []byte) []byte {
	length := len(origData)
	if length == 0 {
		return origData
	}
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

/*
AES256 aes-256-cbc 使用PKCS#7填充加密,產出密文
*/
func AesIvEncrypt(plaintext string, key string, iv string) string {

	bKey := []byte(key)
	bIV := []byte(iv)
	bPlaintext := pkcs7Padding([]byte(plaintext))
	ciphertext := make([]byte, len(bPlaintext))
	block, err := aes.NewCipher(bKey)
	if err != nil {
		return ""
	}

	mode := cipher.NewCBCEncrypter(block, bIV)
	mode.CryptBlocks(ciphertext, bPlaintext)
	return string(ciphertext)
}

/*
AES256 aes-256-cbc 使用PKCS#7解密
*/
func AesIvDecrypt(key, iv []byte, aesEncodeData []byte) []byte {

	block, err := aes.NewCipher(key)
	if err != nil {
		return errByte
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(aesEncodeData, aesEncodeData)
	origData := pkcs7UnPadding(aesEncodeData)

	return origData
}

//AesEncrypt 加密
func AesEncrypt(data []byte, key []byte) ([]byte, error) {
	//创建加密实例
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	//判断加密快的大小
	blockSize := block.BlockSize()
	//填充
	encryptBytes := pkcs7Padding(data)
	//初始化加密数据接收切片
	crypted := make([]byte, len(encryptBytes))
	//使用cbc加密模式
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	//执行加密
	blockMode.CryptBlocks(crypted, encryptBytes)
	return crypted, nil
}

//AesDecrypt 解密
func AesDecrypt(data []byte, key []byte) ([]byte, error) {
	//创建实例
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	//获取块的大小
	blockSize := block.BlockSize()
	//使用cbc
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	//初始化解密数据接收切片
	crypted := make([]byte, len(data))
	//执行解密
	blockMode.CryptBlocks(crypted, data)
	//去除填充
	crypted = pkcs7UnPadding(crypted)
	if err != nil {
		return nil, err
	}
	return crypted, nil
}
