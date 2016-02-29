# Ingestion Agent (PHP)
Handles incomng http requests from the boundless web and pushes data to Redis. Implemented in PHP.

__Specification__

1. Accept incoming http request
2. Push a "postback" object to Redis for each "data" object contained in accepted request.

__Sample Request__

*Post*

`http://{server_ip}/ingest.php`

*Raw Post Data*

```json 
{ 
  "endpoint":	{
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

__Created By__

```bash
curl -X POST -H "Content-Type: application/json" -d '{ "endpoint": { "method":"GET", "url":"http://sample_domain_endpoint.com/data?key={key}&value={value}&foo={bar}" }, "data":[ { "key":"Azureus", "value":"Dendrobates" }, { "key":"Phyllobates", "value":"Terribilis" } ] }' http://{server-ip}/ingest.php
```
