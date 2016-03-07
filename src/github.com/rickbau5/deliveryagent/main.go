package main

import (
    "encoding/json"
    "fmt"
    "gopkg.in/redis.v3"
    "io/ioutil"
    "log"
    "net/http"
    "net/url"
    "os"
    "regexp"
    "strings"
    "time"
)

func redisClient() *redis.Client {
    client := redis.NewClient(&redis.Options{
        Addr: "127.0.0.1:6379",
        Password: "",
        DB: 0,
    })    
    return client
}

func setupLogger(path string) *os.File {
    file, err := os.OpenFile(path, os.O_CREATE | os.O_RDWR | os.O_APPEND, 0666)
    if err != nil {
        panic("Couldn't open log file.")
    }
    log.SetOutput(file)
    return file
}

func main() {
    if len(os.Args) <= 1 {
        fmt.Println("Must supply log file path.\n  deliveryagent /path/to/file.log")
        os.Exit(1)
    }
    client := redisClient()
    logFile := setupLogger(os.Args[1])
    defer logFile.Close()

    ping := client.Ping()
    if _, errPing := ping.Result(); errPing == nil {
        for {
            if str, errPop := client.BRPop(0, "requests").Result(); errPop == nil {
                deliveryStart := time.Now()
                mapped := jsonStringToMap(str[1])

                if end, ok := mapped["endpoint"]; ok {
                    endpoint := end.(map[string]interface{})
                    data := mapped["data"].(map[string]interface{})

                    sendResponse(endpoint["url"].(string), data, endpoint["method"].(string), deliveryStart)
                } else {
                    log.Println("No endpoint in request.")
                    log.Println(mapped)
                }
            } else {
                //on error
                log.Println("Error while popping.")
                log.Println(errPop)
            }
        }
    } else {
        //Log error
        log.Println(errPing)
        log.Panicln("Couldn't ping the database.")
    }
}

func constructGet(unformattedUrl string, dataMap map[string]interface{}) string {
    formatted := unformattedUrl
    for key, val := range dataMap {
        formatted = strings.Replace(formatted, braced(key), val.(string), 1)
    }
    regex := regexp.MustCompile("{[[:word:]]*}")
    formatted = regex.ReplaceAllString(formatted, "")
    return formatted
}

func constructPost(unformattedUrl string, dataMap map[string]interface{}) (string, url.Values) {
    if idx := strings.Index(unformattedUrl, "?"); idx != -1 {
        unformattedUrl = unformattedUrl[:idx]
    }
    data := url.Values{}
    for key, val := range dataMap {
        data.Add(key, val.(string))
    }

    return unformattedUrl, data
}

func sendResponse(uResponse string, dataMap map[string]interface{}, method string, deliveryStart time.Time) {
    url := uResponse    // a bit of a naive way to do it. Not prod ;)
    if split := strings.LastIndex(uResponse, "/"); split > 0 {
        url = uResponse[:split]
    }
    var response *http.Response
    var err error
    var deliveryEnd time.Time
    var start time.Time
    var end time.Time

    if method == "GET" {
        resp := constructGet(uResponse, dataMap)
        deliveryEnd = time.Now()

        start = time.Now()
        response, err = http.Get(resp)
        end = time.Now()
    } else if method == "POST" {
        resp, values := constructPost(uResponse, dataMap)
        deliveryEnd = time.Now()

        start = time.Now()
        response, err = http.PostForm(resp, values)
        end = time.Now()
    } else {
        log.Println("Unknown method type %s for %s.\n", method, url)
    }

    if err == nil {
        log.Printf("Delivery to %s processed in %dms.\n", url, (deliveryEnd.Nanosecond() - deliveryStart.Nanosecond()) / 1000000)

        defer response.Body.Close()
        if body, err := ioutil.ReadAll(response.Body); err == nil {
            log.Printf("Code: %s in %dms.\n", response.Status, (end.Nanosecond() - start.Nanosecond()) / 1000000)
            log.Printf("Body: %s\n", body)
        } else {
            log.Println("Error reading body:", err)
        }
    } else {
        log.Println("Error with response:", err)
    }
}

func braced(s string) string {
    return fmt.Sprintf("{%s}", s)
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
