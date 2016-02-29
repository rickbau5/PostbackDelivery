# Ingestion Agent

__Specification__

1. Accept incoming http request
2. Push a "postback" object to Redis for each "data" object contained in accepted request.

__Sample Request__

(POST) http://{server_ip}/ingest.php
(RAW POST DATA) {  
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

__Created By__

`curl -X POST -H "Content-Type: application/json" -d '{ "endpoint": { "method":"GET", "url":"http://sample_domain_endpoint.com/data?key={key}&value={value}&foo={bar}" }, "data":[ { "key":"Azureus", "value":"Dendrobates" }, { "key":"Phyllobates", "value":"Terribilis" } ] }' http://{server-ip}/ingest.php`
