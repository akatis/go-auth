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
#### 2- Secondly create Auth variable.
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

### Creating Access Token
You can generate tokens using the CreateAccessToken function, which is the method of the Auth struct.
<br>The generated access token is added to the AccessToken field in the Auth struct.
```go
token := a.CreateAccessToken("uuid3", nil)
```
The CreateAccess Token function takes two parameters. The first is the uuid of the token holder and the second is its role.
<br>However, if your system does not use role information, you can specify it as "nil".

### Adding to Redis
Redis NoSQL database is used to control the session of the users.<br>
Using Redis's Hash variable, users are registered according to their uuid. In the registration information, the "field" field contains the payload information of the user's token information. The "value" field contains the "user-agent" information.
<br>By using the AddToRedis function, the user's session is saved to Redis and all sessions opened by the user are checked.
```go
_ = a.AddToRedis("uuid3", "user agent")
```

### Deleting from Redis
If the user's refresh token has expired, the user is logged out. Session information is also deleted from Redis.
```go
_ = a.DeleteFromRedis(a.Payload)
```
