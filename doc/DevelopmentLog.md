# Development Log
This doc contains the development log for this project.

## Setting it Up
First things first, rest in peace Un...wait no. 

First things first, had to set up the environment. 
- Installed [Redis](http://redis.io/) from source, using `make` (installing all deps as I went). 
- Installed [Go](https://golang.org/) using `apt-get`, adding `go/bin` to `$PATH` and specified `$GOPATH`.
- Installed [PHP](http://php.net/) using `apt-get`. 
- Installed [Git](https://git-scm.com/) using `apt-get`
  - Added identity file for my Github account
  
At this point I had the environment set up with what I could immediately tell that was needed from my interpretation of the specication. From here, I moved on to Phase I

## Phase I - Ingesting
The main hurdle in Phase I was learning the basics of **PHP**. Also I immediately realized I needed a way to generate the sample request (as detailed in the [Ingestion Agent doc](doc/IngestionAgent.md)), so I went to work on this first.

### Generating the Sample Request
After playing around with cURL and seeing various examples, I realized it was formatted as a JSON object.

*Sample Request*
``` JSON
{ 
  "endpoint":   {
    "method":"GET",
    "url":"http://sample_domain_endpoint.com/data?key={key}&value={value}&foo={bar}"
  },
  "data":[
    {
      "key":"Azureus",
      "value":"Dendrobates"
    },
    {
      "key":"Phyllobates",
      "value":"Terribilis"
    }
  ]
}
```

I came up with this command (taken from the [Ingestion Agent doc](doc/IngestionAgent.md)) using cURL to post the request.

``` bash
curl -X POST -H "Content-Type: application/json" -d '{ "endpoint": { "method":"GET", "url":"http://sample_domain_endpoint.com/data?key={key}&value={value}&foo={bar}" }, "data":[ { "key":"Azureus", "value":"Dendrobates" }, { "key":"Phyllobates", "value":"Terribilis" } ] }' http://{server-ip}/ingest.php
```

### Ingesting the Request
Now that I was able to post the request to a server, I needed to get at the post data through PHP. I ran into the issue here that $_POST didn't contain anything. I had already done some testing and had found that it could handle other POST requests so I knew that there wasn't a configuration issue or something to that effect. Eventually I found that this was because PHP doesn't decode JSON content type automatically, like it does for several other types of POST requests. I immediately found, though, that the data was accessible from `php://input`.

I now had the input using `file_get_contents('php://input')` but doing a `var_dump` on it, I saw that it was of type `string`. Which meant, I couldn't get at the contents in a nice way (like accessing an array/map, for instance). I needed to decode this string into a JSON object so I could easily access the data within it. I quickly found that this was done using in the following manner: `json_decode($postdata,true)`, where `true` makes it return a so-called `associative array`.

At this point I had pretty much all I needed to get to the real meat of solving **Phase I**. A not-so-safe-or-pretty version of the ingestion agent at this point can be found below. It also includes a sample of printing out the objects in `data`, which is an important part of this project and phase.

[State of ingest.php](https://github.com/rickbau5/PostbackDelivery/blob/64f262b60e647f52a6274723342e15dfa031afa8/src/ingest.php) at the conclusion of this section.

### Processing the Request
Now to get to processing the actual request. My first pass was completed quickly, but involved little checking and just served as a rudimentary proof of concept. Really nothing fancy going on besides playing with variables, utilizing `str_replace`, and...well that's it. 

Here, I was going for simplicity, so I could get move on to the next part of this phase, which will be pushing the formatted request into Redis. I planned on coming back to this part after I started working with Redis and Go if I need to restructure the data in any way, or in the *unlikely* case this implementation is perfect, come back at the end to add in all the error handling.

[State of ingest.php](https://github.com/rickbau5/PostbackDelivery/blob/4f8bd6cb687a4c188b76a9d6a7f9cd171a97b286/src/ingest.php#L17-L28) at the conclusion of this section with lines of interest highlighted. From here I could move on to the extremely brief Phase II.

## Phase II - Queueing
This phase concerned getting the data into Redis from PHP and then back out to Go. This proved very simple.

### Pushing the processed request into Redis
Now that I had the formatted request ready to be put into Redis, just needed to get that connection set up. To achieve this, I installed [phpredis](https://github.com/phpredis/phpredis), which provides a way to interact with Redis from PHP. As part of this, I had to toy with the PHP config to enable this extension. Also during this stage I created an Upstart service to maintain `redis-server` and ensure that it would start and be up constantly (I think?).

Once the extension was enabled, I really only added two new lines of code to create this interaction:
- `$redis->connect('127.0.0.1')`
- `$redis->lPush('reqests', $populated)` where `$populated` is the formatted URL. (I'm pretty sure I'm going to need to come back to this point and restructure the data a bit, but at this point I'm more interested in getting the whole flow set up.)

[State of ingest.php](https://github.com/rickbau5/PostbackDelivery/blob/7e4991aaadad63114abc796fef5fd886d429ac40/src/ingest.php#L28) at the conclusion of this section. From here, I could move onto the Go portion of this project.

### Pulling the data out of Redis into Go
Here I fought with Go for a bit to figure out exactly how to import the go-redis package to have access to the Redis connector. Once that was resolved, connecting to Redis from Go and verifying that it was working and communicating correctly was simple. [State of delivery agent](https://github.com/rickbau5/PostbackDelivery/blob/abf9f62f9d20266d141035146e48f07c7311049d/src/github.com/rickbau5/deliveryagent/main.go) at this point was essentially the "hello world" of this situation.
