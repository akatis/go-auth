package authTest

import (
	"github.com/akatis/go-auth"
	"time"
)

var a = auth.New(&auth.Config{
	Redis: struct {
		RedisAddr    string
		RedisPass    string
		RedisDb      int
		PoolSize     int
		MaxIdleConns int
		MinIdleConns int
		DialTimeout  time.Duration
	}{RedisAddr: "3.71.180.188:6379", RedisPass: "!*iEcZa20?23*!", RedisDb: 0, PoolSize: 1000, MaxIdleConns: 100, MinIdleConns: 10},
	JwtSecretKey:        "secret_key",
	EndpointPermissions: EndPointPermissions,
})

func GetAuth() *auth.Auth {
	return a
}
