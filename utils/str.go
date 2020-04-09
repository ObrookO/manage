package utils

import (
	"bytes"
	"crypto/rand"
	"math/big"
)

// ObjInIntSlice 判断obj是否在slice中
func ObjInIntSlice(obj int, s []int) bool {
	for _, item := range s {
		if obj == item {
			return true
		}
	}

	return false
}

// RandomStr 生成随机字符串
func RandomStr(len int) string {
	var container string
	var str = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	b := bytes.NewBufferString(str)
	length := b.Len()
	bigInt := big.NewInt(int64(length))
	for i := 0; i < len; i++ {
		randomInt, _ := rand.Int(rand.Reader, bigInt)
		container += string(str[randomInt.Int64()])
	}
	return container
}
