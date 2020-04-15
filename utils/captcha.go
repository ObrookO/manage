package utils

import (
	"github.com/mojocn/base64Captcha"
)

//type redisStore struct {
//	conn redis.Conn
//}
//
//func newRedisStore() *redisStore {
//	c, _ := redis.Dial("tcp", beego.AppConfig.String("redis_host"))
//	return &redisStore{conn: c}
//}
//
//func (c *redisStore) Set(id string, value string) {
//}
//
//func (c *redisStore) Get(id string, clear bool) string {
//	return ""
//}
//
//func (c *redisStore) Verify(id, answer string, clear bool) bool {
//	return false
//}

// 获取验证码
func GetCaptcha() (string, string, error) {
	//store := newRedisStore()
	captcha := base64Captcha.NewCaptcha(base64Captcha.DefaultDriverDigit, base64Captcha.DefaultMemStore)
	return captcha.Generate()
}
