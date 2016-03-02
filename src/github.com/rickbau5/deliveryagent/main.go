package main

import (
    "encoding/json"
    "fmt"
    "gopkg.in/redis.v3"
    "io/ioutil"
    "net/http"
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

func jsonStringToMap(s string) map[string]interface{} {
    var requestData map[string]interface{}
    converted := []byte(s)
    jsonErr := json.Unmarshal(converted, &requestData)
    if jsonErr != nil {
        panic(jsonErr)
    }
    return requestData
}

func constructResponse(unformattedUrl string, dataMap map[string]interface{}) string {
    formatted := strings.Replace(unformattedUrl, "{key}", dataMap["key"].(string), 1)
    formatted = strings.Replace(formatted, "{value}", dataMap["value"].(string), 1)
    return formatted
} 

func sendResponse(response string, method string) {
    if method == "GET" {
        resp, err := http.Get(response)
        defer resp.Body.Close()
        fmt.Println(method, resp, err)
        body, err := ioutil.ReadAll(resp.Body)
        fmt.Printf("%s\n", body)
    } else if method == "POST" {
        fmt.Println("Gimme a minute or two.")
    } else {
        fmt.Println("Unknown method type", method)
    }    
}

func main() {
    client := newClient()
    
    ping := client.Ping()
    if _, errPing := ping.Result(); errPing == nil {
        popCmd := client.BRPop(0, "requests")
        if popCmd.Err() == nil {
            formatted := formatString(popCmd.String())
            mapped := jsonStringToMap(formatted)
            if end, ok := mapped["endpoint"]; ok {
                endpoint := end.(map[string]interface{})
                fmt.Println(endpoint)
                if d, ok := mapped["data"]; ok {
                    di := d.([]interface{})
                    for _, data := range di {
                        dataMap := data.(map[string]interface{})
                        response := constructResponse(endpoint["url"].(string), dataMap)
                        sendResponse(response, endpoint["method"].(string))
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
    } else {
        //Log error
        fmt.Println(errPing)
    }
}

