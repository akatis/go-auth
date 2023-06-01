package auth

import "time"

type RedisConfig struct {
	RedisAddr    string
	RedisPass    string
	RedisDb      int
	PoolSize     int
	MaxIdleConns int
	MinIdleConns int
	DialTimeout  time.Duration
}

/*type RedisConn struct {
	conn *redis.Conn
}

func NewRedisConn(rc *RedisConfig) *RedisConn {
	var client = redis.NewClient(&redis.Options{
		Addr:         rc.RedisAddr,
		Password:     rc.RedisPass,
		DB:           rc.RedisDb,
		PoolSize:     rc.PoolSize,
		MaxIdleConns: rc.MaxIdleConns,
		MinIdleConns: rc.MinIdleConns,
		DialTimeout:  rc.DialTimeout,
	})

	conn := client.Conn()
	if err := conn.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}

	return &RedisConn{
		conn: conn,
	}

}*/
