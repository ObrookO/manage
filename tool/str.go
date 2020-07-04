package tool

import (
	"github.com/ObrookO/go-utils"
	"github.com/astaxie/beego"
)

// GenerateRawPassword 生成初始密码
func GenerateRawPassword() string {
	return utils.RandomStr(32)
}

// GenerateEncryptedPassword 生成加密后的密码
func GenerateEncryptedPassword(password string) string {
	key := beego.AppConfig.String("aes_key")
	encryptedPassword, err := utils.AesEncrypt(password, key)
	if err != nil {
		return ""
	}

	return encryptedPassword
}
