package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"fmt"

	"github.com/astaxie/beego/logs"
)

// Md5Str 使用md5加密字符串
func Md5Str(str string) string {
	md5Str := fmt.Sprintf("%x", md5.Sum([]byte(str)))
	return md5Str
}

func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// AesEncrypt aes加密
func AesEncrypt(origData, key string) string {
	k := []byte(key)
	block, err := aes.NewCipher(k)
	if err != nil {
		logs.Error(err)
		return ""
	}

	b := []byte(origData)
	blockSize := block.BlockSize()
	b = pkcs5Padding(b, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, k[:blockSize])
	crypted := make([]byte, len(b))
	blockMode.CryptBlocks(crypted, b)
	return base64.StdEncoding.EncodeToString(crypted)
}

// AesDecrypt aes解密
func AesDecrypt(crypted, key string) string {
	k := []byte(key)
	block, err := aes.NewCipher(k)
	if err != nil {
		logs.Error(err)
		return ""
	}

	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, k[:blockSize])
	b, _ := base64.StdEncoding.DecodeString(crypted)
	origData := make([]byte, len(b))
	blockMode.CryptBlocks(origData, b)
	origData = pkcs5UnPadding(origData)
	return string(origData)
}
