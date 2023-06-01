# GO AUTH
## What the purpose of this package?
The purpose of this package is to _generate jwt token_, _provide session control from Redis_ and _provide an authorization middleware_.

## How can use this package?
#### 1- First of all load packages and import.
```cmd 
go get "github.com/akatis/go-auth"
```
```go
import "github.com/akatis/go-auth"
```
#### 2- Secondly create Auth veriable.
Create an "auth" directory and define a global variable there to get the Auth struct with the New method.
<br>Create GetAuth() method to use the auth variable in different directories.
```go 
var a = auth.New(&auth.Config{
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

func GetAuth() *auth.Auth {
return a
}

```
#### 3- Use gofiber/fiber/v2 for middleware.
```go 
app := fiber.New()
api := app.Group("/api")

a := auth.GetAuth()

api.Use(a.Middleware)
```
You can control auth.go file for requirement.
