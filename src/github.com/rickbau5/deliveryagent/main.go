package main

import (
    "fmt"
    "gopkg.in/redis.v3"
    "encoding/json"
    "strings"
)

func newClient() *redis.Client {
    client := redis.NewClient(&redis.Options{
        Addr: "127.0.0.1:6379",
        Password: "",
        DB: 0,
    })    
    return client
}

func formatString(s string) string {
    t := strings.Trim(s, "]")
    formatted := strings.Replace(t, "BRPOP requests 0: [requests ", "", 1)
    return formatted
}

func main() {
    client := newClient()
    
    ping := client.Ping()
    if _, errPing := ping.Result(); errPing == nil {
        popCmd := client.BRPop(0, "requests")
        if popCmd.Err() == nil {
            formatted := formatString(popCmd.String())
            var requestData map[string]interface{}
            converted := []byte(formatted)
            jsonErr := json.Unmarshal(converted, &requestData)
            if jsonErr == nil {
                fmt.Println(requestData)
            } else {
                fmt.Println("Couldn't parse JSON :( ")
                fmt.Println(jsonErr)
            }
        } else {
            //on error
            fmt.Println("Error while popping.")
            fmt.Println(popCmd.Err())
        }
    } else {
        //Log error
        fmt.Println(errPing)
    }
}

