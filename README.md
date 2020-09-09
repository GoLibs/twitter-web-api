# GoLibs twitter-web-api - Twitter Web API in Golang

## This API uses [Twitter Mobile](https://mobile.twitter.com)

### Current Features:
 Get User & Latest Twits.
```go
client := twitter_web_api.NewTwitterClient()
user, err := client.GetUserByUsername("jack")
fmt.Println(user, err)
```
