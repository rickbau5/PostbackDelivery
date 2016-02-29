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
  
At this point I had the environment set up with what I could immediately tell that was needed from my interpretation of the specication. From here, I moved on to Phase 1

## Phase I
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

I know had the input using `file_get_contents('php://input')` but doing a `var_dump` on it, I saw that it was of type `string`. Which meant, I couldn't get at the contents in a nice way (like accessing an array/map, for instance). I needed to decode this string into a JSON object so I could easily access the data within it. I quickly found that this was done using in the following manner: `json_decode($postdata,true)`, where `true` makes it return a so-called `associative array`.

At this point I had pretty much all I needed to get to the real meat of solving **Phase I**. A not-so-safe-or-pretty version of the ingestion agent at this point can be found [here](https://github.com/rickbau5/PostbackDelivery/blob/64f262b60e647f52a6274723342e15dfa031afa8/src/ingest.php). It also includes a sample of printing out the objects in `data`, which is an important part of this project and phase.
