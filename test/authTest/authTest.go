package authTest

import (
	"github.com/akatis/go-auth"
	"time"
)

var a = auth.NewConfig(&auth.Config{
	Redis: struct {
		RedisAddr    string
		RedisPass    string
		RedisDb      int
		PoolSize     int
		MaxIdleConns int
		MinIdleConns int
		DialTimeout  time.Duration
	}{RedisAddr: "localhost:6379", RedisPass: "password", RedisDb: 0, PoolSize: 1000, MaxIdleConns: 100, MinIdleConns: 10},
	JwtSecretKey: "secret_key",
})

func GetConf() *auth.Auth {
	return a
}
