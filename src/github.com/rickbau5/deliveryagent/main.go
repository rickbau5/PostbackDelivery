package main

import (
    "fmt"
    "gopkg.in/redis.v3"
)

func newClient() *redis.Client {
    client := redis.NewClient(&redis.Options{
        Addr: "127.0.0.1:6379",
        Password: "",
        DB: 0,
    })    
    return client
}

func main() {
    fmt.Println("Hello world")
    client := newClient()
    
    pong, err := client.Ping().Result()
    fmt.Println(pong, err)
}

