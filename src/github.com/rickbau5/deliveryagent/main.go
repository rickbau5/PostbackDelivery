package main

import (
    "encoding/json"
    "fmt"
    "gopkg.in/redis.v3"
    "io/ioutil"
    "net/http"
    "net/url"
    "strings"
)

func redisClient() *redis.Client {
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

func jsonStringToMap(s string) map[string]interface{} {
    var requestData map[string]interface{}
    converted := []byte(s)
    jsonErr := json.Unmarshal(converted, &requestData)
    if jsonErr != nil {
        panic(jsonErr)
    }
    return requestData
}

func constructGet(unformattedUrl string, dataMap map[string]interface{}) string {
    formatted := strings.Replace(unformattedUrl, "{key}", dataMap["key"].(string), 1)
    formatted = strings.Replace(formatted, "{value}", dataMap["value"].(string), 1)
    return formatted
} 

func sendResponse(uResponse string, dataMap map[string]interface{}, method string) {
    if method == "GET" {
        response := constructGet(uResponse, dataMap)
        if resp, respErr := http.Get(response); respErr == nil {
            defer resp.Body.Close()
            if body, err := ioutil.ReadAll(resp.Body); err == nil {
                fmt.Printf("%s\n", body)
            } else {
                fmt.Println("Error reading body", err)
            }
        } else {
            fmt.Println("Error with response", respErr)
        }
    } else if method == "POST" {
        response := uResponse
        if idx := strings.Index(response, "?"); idx != -1 {
            response = response[:idx]
        }
        values := url.Values{"key":{dataMap["key"].(string)}, "value":{dataMap["value"].(string)}}
        if resp, err := http.PostForm(response, values); err == nil{
            defer resp.Body.Close()
            body, _ := ioutil.ReadAll(resp.Body)
            fmt.Printf("%s\n", body)
        } else {
            fmt.Println(err)
        }
    } else {
        fmt.Println("Unknown method type", method)
    }    
}

func main() {
    client := redisClient()
    
    ping := client.Ping()
    if _, errPing := ping.Result(); errPing == nil {
        for {
            popCmd := client.BRPop(0, "requests")
            if popCmd.Err() == nil {
                formatted := formatString(popCmd.String())
                mapped := jsonStringToMap(formatted)
                if end, ok := mapped["endpoint"]; ok {
                    endpoint := end.(map[string]interface{})
                    if d, ok := mapped["data"]; ok {
                        di := d.([]interface{})
                        for _, data := range di {
                            dataMap := data.(map[string]interface{})
                            sendResponse(endpoint["url"].(string), dataMap, endpoint["method"].(string))
                        }
                    }
                } else {
                    fmt.Println("Endpoint is nil.")
                }
            } else {
                //on error
                fmt.Println("Error while popping.")
                fmt.Println(popCmd.Err())
            }
        }
    } else {
        //Log error
        fmt.Println(errPing)
    }
}

